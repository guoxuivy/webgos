package middleware

import (
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"webgos/common/json"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
)

const (
	maxPathLen   = 1024      // 路径最大检测长度，防止 ReDoS
	scanWindow   = time.Hour // 敏感路径命中统计窗口
	banThreshold = 5         // 窗口内命中次数阈值
)

// 敏感路径模式列表
var sensitivePatterns = []struct {
	pattern *regexp.Regexp
	name    string
}{
	{regexp.MustCompile(`(?i)\.env`), "env-file"},
	{regexp.MustCompile(`(?i)\.bak`), "backup-file"},
	{regexp.MustCompile(`(?i)config\.(json|yaml|yml|toml|ini)`), "config-file"},
	{regexp.MustCompile(`(?i)\.git/config`), "git-config"},
	{regexp.MustCompile(`(?i)\.gitignore`), "gitignore"},
	{regexp.MustCompile(`(?i)\.svn`), "svn"},
	{regexp.MustCompile(`(?i)adminer`), "adminer"},
	{regexp.MustCompile(`(?i)phpmyadmin`), "phpmyadmin"},
	{regexp.MustCompile(`(?i)wp-admin`), "wp-admin"},
	{regexp.MustCompile(`(?i)wp-content`), "wp-content"},
	{regexp.MustCompile(`(?i)\.sql`), "sql-dump"},
	{regexp.MustCompile(`(?i)\.(log|txt|md)$`), "log-txt-file"},
	{regexp.MustCompile(`(?i)\.(tar|gz|zip|rar|7z)`), "archive-file"},
	{regexp.MustCompile(`(?i)robots\.txt`), "robots-txt"},
	{regexp.MustCompile(`(?i)crossdomain\.xml`), "crossdomain"},
	{regexp.MustCompile(`(?i)\.aws/credentials`), "aws-credentials"},
	{regexp.MustCompile(`(?i)\.ssh`), "ssh-dir"},
	{regexp.MustCompile(`(?i)composer\.(json|lock)`), "composer"},
	{regexp.MustCompile(`(?i)npmrc`), "npmrc"},
	{regexp.MustCompile(`(?i)docker-compose`), "docker-compose"},
	{regexp.MustCompile(`(?i)Dockerfile`), "dockerfile"},
}

// ===== 恶意 IP 追踪器 =====
type MaliciousIPRecord struct {
	IP        string
	Path      string
	Pattern   string
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

type MaliciousIPTracker struct {
	mu       sync.RWMutex
	records  map[string]*MaliciousIPRecord
	maxTrack int
}

var maliciousIPTracker = &MaliciousIPTracker{
	records:  make(map[string]*MaliciousIPRecord),
	maxTrack: 10000,
}

func IsSensitivePath(path string) bool {
	if len(path) > maxPathLen {
		path = path[:maxPathLen]
	}
	for _, sp := range sensitivePatterns {
		if sp.pattern.MatchString(path) {
			return true
		}
	}
	lowerPath := strings.ToLower(path)
	sensitiveKeywords := []string{"shell", "cmd", "exec", "eval", "passwd", "shadow", "htaccess"}
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lowerPath, kw) {
			return true
		}
	}
	return false
}

func getSensitivePatternName(path string) string {
	for _, sp := range sensitivePatterns {
		if sp.pattern.MatchString(path) {
			return sp.name
		}
	}
	return "unknown"
}

func LogSensitiveAccess(c *gin.Context) {
	ip := c.ClientIP()
	path := c.Request.URL.Path
	pattern := getSensitivePatternName(path)

	xlog.Warn("[SECURITY] 敏感路径访问 IP=%s Path=%s Pattern=%s UserAgent=%s",
		ip, path, pattern, c.Request.UserAgent())

	maliciousIPTracker.record(ip, path, pattern)

	if maliciousIPTracker.shouldBan(ip) {
		AddToBlacklist(ip)
		xlog.Warn("[SECURITY] 恶意IP自动封禁 IP=%s", ip)
	}
}

func (t *MaliciousIPTracker) record(ip, path, pattern string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	record, exists := t.records[ip]

	if exists {
		// 超过统计窗口则重置计数，避免误触用户被永久拉黑
		if now.Sub(record.FirstSeen) > scanWindow {
			record.FirstSeen = now
			record.Count = 1
			record.LastSeen = now
			record.Path = path
			return
		}
		record.LastSeen = now
		record.Count++
		record.Path = path
		return
	}

	if len(t.records) >= t.maxTrack {
		var oldestIP string
		var oldestTime time.Time
		for k, v := range t.records {
			if oldestTime.IsZero() || v.FirstSeen.Before(oldestTime) {
				oldestIP = k
				oldestTime = v.FirstSeen
			}
		}
		delete(t.records, oldestIP)
	}

	t.records[ip] = &MaliciousIPRecord{
		IP:        ip,
		Path:      path,
		Pattern:   pattern,
		FirstSeen: now,
		LastSeen:  now,
		Count:     1,
	}
}

func (t *MaliciousIPTracker) shouldBan(ip string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	record, exists := t.records[ip]
	return exists && record.Count >= banThreshold
}

func GetMaliciousIPs() []MaliciousIPRecord {
	maliciousIPTracker.mu.RLock()
	defer maliciousIPTracker.mu.RUnlock()

	result := make([]MaliciousIPRecord, 0, len(maliciousIPTracker.records))
	for _, record := range maliciousIPTracker.records {
		result = append(result, *record)
	}
	return result
}

func ClearMaliciousIPs() {
	maliciousIPTracker.mu.Lock()
	defer maliciousIPTracker.mu.Unlock()
	maliciousIPTracker.records = make(map[string]*MaliciousIPRecord)
}

// ===== IP 黑名单管理器 =====

type BlacklistData struct {
	IPs     []string `json:"ips"`
	CIDRs   []string `json:"cidrs"`
	Version string   `json:"version"`
}

type IPBlacklist struct {
	mu        sync.RWMutex
	ips       map[string]bool
	cidrs     []*net.IPNet
	enable    bool
	savePath  string
	saveTimer *time.Timer
}

var GlobalIPBlacklist = &IPBlacklist{
	ips:    make(map[string]bool),
	enable: true,
}

func InitIPBlacklist(dir string) {
	GlobalIPBlacklist.savePath = filepath.Join(dir, "blacklist.json")
	GlobalIPBlacklist.loadFromFile()
	GlobalIPBlacklist.startAutoSave()
}

func (b *IPBlacklist) loadFromFile() {
	b.mu.Lock()
	defer b.mu.Unlock()

	data, err := os.ReadFile(b.savePath)
	if err != nil {
		if !os.IsNotExist(err) {
			xlog.Warn("[SECURITY] 加载黑名单文件失败: %v", err)
		}
		return
	}

	var blacklistData BlacklistData
	if err := json.Unmarshal(data, &blacklistData); err != nil {
		xlog.Warn("[SECURITY] 解析黑名单文件失败: %v", err)
		return
	}

	b.ips = make(map[string]bool)
	b.cidrs = b.cidrs[:0]

	for _, ip := range blacklistData.IPs {
		parsedIP := net.ParseIP(ip)
		if parsedIP != nil {
			b.ips[ip] = true
		}
	}

	for _, cidrStr := range blacklistData.CIDRs {
		_, cidr, err := net.ParseCIDR(cidrStr)
		if err == nil {
			b.cidrs = append(b.cidrs, cidr)
		}
	}

	xlog.Info("[SECURITY] 黑名单加载完成 IP=%d CIDR=%d", len(b.ips), len(b.cidrs))
}

func (b *IPBlacklist) saveToFile() {
	b.mu.RLock()
	defer b.mu.RUnlock()

	blacklistData := BlacklistData{
		IPs:     make([]string, 0, len(b.ips)),
		CIDRs:   make([]string, 0, len(b.cidrs)),
		Version: "1.0",
	}

	for ip := range b.ips {
		blacklistData.IPs = append(blacklistData.IPs, ip)
	}

	for _, cidr := range b.cidrs {
		blacklistData.CIDRs = append(blacklistData.CIDRs, cidr.String())
	}

	data, err := json.MarshalIndent(blacklistData, "", "  ")
	if err != nil {
		xlog.Warn("[SECURITY] 序列化黑名单失败: %v", err)
		return
	}

	err = os.MkdirAll(filepath.Dir(b.savePath), 0755)
	if err != nil {
		xlog.Warn("[SECURITY] 创建目录失败: %v", err)
		return
	}

	if err := os.WriteFile(b.savePath, data, 0644); err != nil {
		xlog.Warn("[SECURITY] 保存黑名单文件失败: %v", err)
	}
}

func (b *IPBlacklist) startAutoSave() {
	b.saveTimer = time.AfterFunc(5*time.Minute, func() {
		b.saveToFile()
		b.startAutoSave()
	})
}

func (b *IPBlacklist) Stop() {
	if b.saveTimer != nil {
		b.saveTimer.Stop()
		b.saveToFile()
	}
}

func AddToBlacklist(ip string) {
	GlobalIPBlacklist.mu.Lock()
	defer GlobalIPBlacklist.mu.Unlock()

	parsed := net.ParseIP(ip)
	if parsed == nil {
		return
	}
	GlobalIPBlacklist.ips[ip] = true
	xlog.Warn("[SECURITY] 添加黑名单 IP=%s", ip)
}

func RemoveFromBlacklist(ip string) {
	GlobalIPBlacklist.mu.Lock()
	defer GlobalIPBlacklist.mu.Unlock()
	delete(GlobalIPBlacklist.ips, ip)
	xlog.Info("[SECURITY] 移除黑名单 IP=%s", ip)
}

func IsBlacklisted(ip string) bool {
	GlobalIPBlacklist.mu.RLock()
	defer GlobalIPBlacklist.mu.RUnlock()

	if !GlobalIPBlacklist.enable {
		return false
	}

	if GlobalIPBlacklist.ips[ip] {
		return true
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range GlobalIPBlacklist.cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}

func GetBlacklistCount() int {
	GlobalIPBlacklist.mu.RLock()
	defer GlobalIPBlacklist.mu.RUnlock()
	return len(GlobalIPBlacklist.ips) + len(GlobalIPBlacklist.cidrs)
}

func IPBlacklistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if IsBlacklisted(ip) {
			xlog.Warn("[SECURITY] 黑名单IP请求被拒绝 IP=%s Path=%s", ip, c.Request.URL.Path)
			c.JSON(403, gin.H{
				"code":    403,
				"message": "请求被拒绝",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
