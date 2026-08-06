package services

import (
	"context"
	"errors"

	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xdb"
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

	db := xdb.GetDB().WithContext(ctx)

	if user.ID > 0 {
		return db.Updates(user).Error
	}

	return db.Create(user).Error
}

func (s *userService) ResetPassword(ctx context.Context, username, password string) error {
	var user models.User

	if err := xdb.GetDB().WithContext(ctx).Where("username = ?", username).Take(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return xdb.GetDB().WithContext(ctx).Model(&user).Update("Password", user.Password).Error
}

func (s *userService) UsersPage(ctx context.Context, query dto.UserQuery) ([]models.User, int64) {
	var users []models.User
	var total int64

	db := xdb.GetDB().WithContext(ctx).Model(&models.User{})

	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return []models.User{}, 0
	}
	db = db.Scopes(models.Page(query.Page, query.PageSize))
	if err := db.Preload("Roles").Find(&users).Error; err != nil {
		return []models.User{}, 0
	}
	return users, total
}

func (s *userService) GetUserInfo(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := xdb.GetDB().WithContext(ctx).Preload("Roles").First(&user, userID).Error
	return &user, err
}
