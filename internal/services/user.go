package services

import (
	"errors"

	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xdb"
)

type UserService interface {
	CreateOrUpdateUser(user *models.User) error
	ResetPassword(username, password string) error
	UsersPage(query dto.UserQuery) ([]models.User, int64)
	GetUserInfo(userID int) (*models.User, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) CreateOrUpdateUser(user *models.User) error {
	if user.Password != "" {
		if err := user.SetPassword(user.Password); err != nil {
			return err
		}
	}

	if user.ID > 0 {
		return xdb.GetDB().Updates(user).Error
	}

	return xdb.GetDB().Create(user).Error
}

func (s *userService) ResetPassword(username, password string) error {
	var user models.User

	if err := xdb.GetDB().Where("username = ?", username).Take(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return xdb.GetDB().Model(&user).Update("Password", user.Password).Error
}

func (s *userService) UsersPage(query dto.UserQuery) ([]models.User, int64) {
	var users []models.User
	var total int64

	db := xdb.GetDB().Model(&models.User{})

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

func (s *userService) GetUserInfo(userID int) (*models.User, error) {
	var user models.User
	err := xdb.GetDB().Preload("Roles").First(&user, userID).Error
	return &user, err
}
