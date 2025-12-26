package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/patrickmn/go-cache"

	"webgos/internal/config"
	"webgos/internal/models"
)

// 全局缓存实例（过期时间24小时，每10分钟清理一次过期数据）
var tokenCache = cache.New(24*time.Hour, 10*time.Minute)

// AuthService 认证服务接口
type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
	Logout(tokenString string)
}

// authService 实现 AuthService 接口
type authService struct{}

// NewAuthService 创建认证服务实例
func NewAuthService() AuthService {
	return &authService{}
}

// Login 用户登录并生成JWT令牌
func (s *authService) Login(username, password string) (string, error) {
	var userModel models.User

	jwtConfig := config.GlobalConfig.JWT

	// 根据用户名查找用户（使用 BaseModel）
	user, err := userModel.Where("username = ?", username).One()
	if err != nil {
		return "", errors.New("用户不存在")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return "", errors.New("密码错误")
	}

	// 生成JWT令牌
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(jwtConfig.Expiry)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	// 将令牌存入缓存，设置与令牌相同的过期时间
	tokenCache.Set(tokenString, true, time.Duration(jwtConfig.Expiry)*time.Hour)

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func (s *authService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	// 检查令牌是否已被登出（从缓存中删除）
	if _, found := tokenCache.Get(tokenString); !found {
		return nil, errors.New("令牌已失效")
	}

	jwtConfig := config.GlobalConfig.JWT
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的令牌")
}

// 登出：从缓存中删除令牌（使其立即失效）
func (s *authService) Logout(tokenString string) {
	tokenCache.Delete(tokenString)
}
