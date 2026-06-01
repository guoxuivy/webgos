package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"webgos/internal/cache"
	"webgos/internal/config"
	"webgos/internal/xdb"
	"webgos/internal/models"
)

type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
	Logout(tokenString string)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Login(username, password string) (string, error) {
	jwtConfig := config.GlobalConfig.JWT

	var user models.User
	if err := xdb.GetDB().Where("username = ?", username).Take(&user).Error; err != nil {
		return "", errors.New("用户不存在")
	}

	if !user.CheckPassword(password) {
		return "", errors.New("密码错误")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(jwtConfig.Expiry)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", err
	} // 将令牌存入缓存，设置与令牌相同的过期时间
	// todo也可以使用黑名单方式实现登出功能，不用缓存储大量令牌
	cache.GetCache().Set(tokenString, true, time.Duration(jwtConfig.Expiry)*time.Hour)
	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	if _, found := cache.GetCache().Get(tokenString); !found {
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

func (s *authService) Logout(tokenString string) {
	cache.GetCache().Delete(tokenString)
}
