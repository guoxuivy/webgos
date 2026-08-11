package services

import (
	"context"
	"errors"

	"webgos/internal/cache"
	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xlog"
)

type UserService interface {
	CreateOrUpdateUser(ctx context.Context, user *models.User) error
	ResetPassword(ctx context.Context, username, password string) error
	UsersPage(ctx context.Context, query dto.UserQuery) ([]models.User, int64)
	GetUserInfo(ctx context.Context, userID int) (*models.User, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) CreateOrUpdateUser(ctx context.Context, user *models.User) error {
	if user.Password != "" {
		if err := user.SetPassword(user.Password); err != nil {
			return err
		}
	}

	db := ctxDB(ctx)

	if user.ID > 0 {
		return db.Updates(user).Error
	}

	return db.Create(user).Error
}

func (s *userService) ResetPassword(ctx context.Context, username, password string) error {
	var user models.User

	if err := ctxDB(ctx).Where("username = ?", username).Take(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return ctxDB(ctx).Model(&user).Update("Password", user.Password).Error
}

func (s *userService) UsersPage(ctx context.Context, query dto.UserQuery) (users []models.User, total int64) {
	cacheKey := cache.GenerateKey(cache.UserPagePrefix, query)
	if cache.GetPage(cacheKey, &users, &total) {
		return users, total
	}

	db := ctxSDB(ctx).Model(&models.User{})

	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		xlog.Error("user page error:%v", err)
		return users, total
	}
	db = db.Scopes(models.Page(query.Page, query.PageSize))
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return users, total
	}
	cache.SetPage(cacheKey, users, total) // 仅成功才缓存
	return users, total
}

func (s *userService) GetUserInfo(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := ctxSDB(ctx).Preload("Roles").First(&user, userID).Error
	return &user, err
}
