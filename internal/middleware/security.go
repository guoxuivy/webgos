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
	{regexp.MustCompile(`(?i)\.env\b`), "env-file"},
	{regexp.MustCompile(`(?i)\.bak\b`), "backup-file"},
	{regexp.MustCompile(`(?i)config\.(json|yaml|yml|toml|ini)`), "config-file"},
	{regexp.MustCompile(`(?i)\.git/config`), "git-config"},
	{regexp.MustCompile(`(?i)\.gitignore`), "gitignore"},
	{regexp.MustCompile(`(?i)\.svn`), "svn"},
	{regexp.MustCompile(`(?i)adminer`), "adminer"},
	{regexp.MustCompile(`(?i)phpmyadmin`), "phpmyadmin"},
	{regexp.MustCompile(`(?i)wp-admin`), "wp-admin"},
	{regexp.MustCompile(`(?i)wp-content`), "wp-content"},
	{regexp.MustCompile(`(?i)\.sql\b`), "sql-dump"},
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

// ============================================================
// 公开 API
// ============================================================

// CheckSensitivePath 检测并记录敏感路径访问，在 404 handler 中调用
func CheckSensitivePath(c *gin.Context) {
	path := c.Request.URL.Path
	if len(path) > maxPathLen {
		path = path[:maxPathLen]
	}

	isHit := false
	for _, sp := range sensitivePatterns {
		if sp.pattern.MatchString(path) {
			isHit = true
			break
		}
	}
	if !isHit {
		lowerPath := strings.ToLower(path)
		for _, kw := range []string{"/shell", "/cmd", "/exec", "/eval", "/passwd", "/shadow", "/htaccess"} {
			if strings.Contains(lowerPath, kw) {
				isHit = true
				break
			}
		}
	}
	if !isHit {
		return
	}

	ip := c.ClientIP()
	pattern := getSensitivePatternName(path)

	xlog.Warn("[SECURITY] 敏感路径访问 IP=%s Path=%s Pattern=%s UserAgent=%s",
		ip, path, pattern, c.Request.UserAgent())

	trackerIns.record(ip, path, pattern)

	if trackerIns.shouldBan(ip) {
		globalIPBlacklist.add(ip)
		xlog.Warn("[SECURITY] 恶意IP自动封禁 IP=%s", ip)
	}
}

// IPBlacklistMiddleware IP黑名单中间件，在全局路由注册拦截黑名单IP的所有请求
func IPBlacklistMiddleware(dir string) gin.HandlerFunc {
	blacklistOnce.Do(func() {
		globalIPBlacklist.init(dir)
	})

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if globalIPBlacklist.contains(ip) {
			xlog.Warn("[SECURITY] 黑名单IP请求被拒绝 IP=%s Path=%s", ip, c.Request.URL.Path)
			c.AbortWithStatusJSON(403, gin.H{
				"code":    403,
				"message": "请求被拒绝",
			})
			return
		}
		c.Next()
	}
}

// ============================================================
// 敏感路径检测
// ============================================================

// getSensitivePatternName 根据路径匹配敏感模式名称，用于日志记录
func getSensitivePatternName(path string) string {
	if len(path) > maxPathLen {
		path = path[:maxPathLen]
	}
	for _, sp := range sensitivePatterns {
		if sp.pattern.MatchString(path) {
			return sp.name
		}
	}
	return "unknown"
}

// ============================================================
// 恶意 IP 追踪器
// ============================================================

// ipRecord 恶意IP单条命中记录
type ipRecord struct {
	IP        string    // 请求IP
	Path      string    // 最后一次命中的路径
	Pattern   string    // 命中的敏感模式名称
	FirstSeen time.Time // 首次命中时间
	LastSeen  time.Time // 最近命中时间
	Count     int       // 统计窗口内命中次数
}

// iplogTracker 恶意IP访问追踪器，基于时间窗口统计敏感路径命中次数，超过阈值自动封禁
type iplogTracker struct {
	mu       sync.RWMutex
	records  map[string]*ipRecord
	maxTrack int // 最大追踪IP数量
}

// trackerIns 恶意IP追踪器全局实例
var trackerIns = &iplogTracker{
	records:  make(map[string]*ipRecord),
	maxTrack: 10000,
}

// record 记录一次敏感路径访问
func (t *iplogTracker) record(ip, path, pattern string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	rec, exists := t.records[ip]

	if exists {
		// 超过统计窗口则重置计数，避免误触用户被永久拉黑
		if now.Sub(rec.FirstSeen) > scanWindow {
			rec.FirstSeen = now
			rec.Count = 1
			rec.LastSeen = now
			rec.Path = path
			return
		}
		rec.LastSeen = now
		rec.Count++
		rec.Path = path
		return
	}

	// 追踪池满时淘汰最旧记录
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

	t.records[ip] = &ipRecord{
		IP:        ip,
		Path:      path,
		Pattern:   pattern,
		FirstSeen: now,
		LastSeen:  now,
		Count:     1,
	}
}

// shouldBan 检查该IP是否达到封禁阈值
func (t *iplogTracker) shouldBan(ip string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	rec, exists := t.records[ip]
	return exists && rec.Count >= banThreshold
}

// cleanupExpired 清理已过期的命中记录，防止内存泄漏
func (t *iplogTracker) cleanupExpired() {
	t.mu.Lock()
	defer t.mu.Unlock()

	expireDuration := scanWindow * 30
	now := time.Now()
	for ip, rec := range t.records {
		if now.Sub(rec.FirstSeen) > expireDuration {
			delete(t.records, ip)
		}
	}
}

// ============================================================
// IP 黑名单管理器
// ============================================================

// blacklistData 黑名单持久化数据结构
type blacklistData struct {
	IPs     []string `json:"ips"`     // 已封禁IP列表
	CIDRs   []string `json:"cidrs"`   // 已封禁CIDR网段列表
	Version string   `json:"version"` // 数据格式版本
}

// ipBlacklist IP黑名单管理器，支持精确IP和CIDR网段封禁，定期持久化到磁盘
type ipBlacklist struct {
	mu        sync.RWMutex
	ips       map[string]bool // 精确IP黑名单
	cidrs     []*net.IPNet    // CIDR网段黑名单
	savePath  string          // 持久化文件路径
	saveTimer *time.Timer     // 自动保存定时器
}

// globalIPBlacklist 黑名单全局实例
var globalIPBlacklist = &ipBlacklist{
	ips: make(map[string]bool),
}

// blacklistOnce 保证初始化只执行一次
var blacklistOnce sync.Once

// init 初始化黑名单管理器，从文件加载并启动定时保存
func (b *ipBlacklist) init(dir string) {
	b.savePath = filepath.Join(dir, "blacklist.json")
	b.loadFromFile()
	b.startAutoSave()
}

// add 将指定IP加入黑名单
func (b *ipBlacklist) add(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if net.ParseIP(ip) == nil {
		return
	}
	b.ips[ip] = true
	xlog.Warn("[SECURITY] 添加黑名单 IP=%s", ip)
}

// contains 检查指定IP是否在黑名单中（精确匹配和CIDR匹配）
func (b *ipBlacklist) contains(ip string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.ips[ip] {
		return true
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range b.cidrs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}

// loadFromFile 从磁盘加载黑名单数据
func (b *ipBlacklist) loadFromFile() {
	b.mu.Lock()
	defer b.mu.Unlock()

	data, err := os.ReadFile(b.savePath)
	if err != nil {
		if !os.IsNotExist(err) {
			xlog.Warn("[SECURITY] 加载黑名单文件失败: %v", err)
		}
		return
	}

	var bd blacklistData
	if err := json.Unmarshal(data, &bd); err != nil {
		xlog.Warn("[SECURITY] 解析黑名单文件失败: %v", err)
		return
	}

	b.ips = make(map[string]bool)
	b.cidrs = b.cidrs[:0]

	for _, ip := range bd.IPs {
		if net.ParseIP(ip) != nil {
			b.ips[ip] = true
		}
	}

	for _, cidrStr := range bd.CIDRs {
		_, cidr, err := net.ParseCIDR(cidrStr)
		if err == nil {
			b.cidrs = append(b.cidrs, cidr)
		}
	}

	xlog.Warn("[SECURITY] 黑名单加载完成 IP=%d CIDR=%d", len(b.ips), len(b.cidrs))
}

// saveToFile 将黑名单数据持久化到磁盘，间隔调用时顺便清理过期追踪记录
func (b *ipBlacklist) saveToFile() {
	trackerIns.cleanupExpired()

	b.mu.RLock()
	defer b.mu.RUnlock()

	bd := blacklistData{
		IPs:     make([]string, 0, len(b.ips)),
		CIDRs:   make([]string, 0, len(b.cidrs)),
		Version: "1.0",
	}

	for ip := range b.ips {
		bd.IPs = append(bd.IPs, ip)
	}

	for _, cidr := range b.cidrs {
		bd.CIDRs = append(bd.CIDRs, cidr.String())
	}

	bytes, err := json.MarshalIndent(bd, "", "  ")
	if err != nil {
		xlog.Warn("[SECURITY] 序列化黑名单失败: %v", err)
		return
	}

	if err := os.MkdirAll(filepath.Dir(b.savePath), 0755); err != nil {
		xlog.Warn("[SECURITY] 创建目录失败: %v", err)
		return
	}

	if err := os.WriteFile(b.savePath, bytes, 0644); err != nil {
		xlog.Warn("[SECURITY] 保存黑名单文件失败: %v", err)
	}
}

// startAutoSave 启动定时自动保存，每5分钟执行一次
func (b *ipBlacklist) startAutoSave() {
	b.saveTimer = time.AfterFunc(5*time.Minute, func() {
		b.saveToFile()
		b.startAutoSave()
	})
}

// stop 停止定时器并立即保存一次
func (b *ipBlacklist) stop() {
	if b.saveTimer != nil {
		b.saveTimer.Stop()
		b.saveToFile()
	}
}
