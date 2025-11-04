package services

import (
	"errors"

	"webgos/internal/database"
	"webgos/internal/dto"
	"webgos/internal/models"
)

// UserService 用户服务接口
type UserService interface {
	CreateOrUpdateUser(user *models.User) error
	ResetPassword(username, password string) error
	UsersPage(query dto.UserQuery) ([]models.User, int)
}

// userService 实现 UserService 接口
type userService struct{}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// CreateUser 创建更新用户
func (s *userService) CreateOrUpdateUser(user *models.User) error {
	// 密码加密
	if user.Password != "" {
		if err := user.SetPassword(user.Password); err != nil {
			return err
		}
	}

	if user.ID > 0 {
		// 更新用户
		if err := user.Update(user); err != nil {
			return err
		}
		return nil
	}

	// 创建用户
	if err := user.Create(user); err != nil {
		return err
	}
	return nil
}

// ResetPassword 重置用户密码
func (s *userService) ResetPassword(username, password string) error {
	var user models.User

	// 根据用户名查找用户
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 设置新密码
	if err := user.SetPassword(password); err != nil {
		return err
	}

	// 只更新密码字段，避免更新用户名导致唯一性约束冲突
	err := database.DB.Model(&user).Select("Password").Updates(&user).Error

	if err != nil {
		return err
	}

	return nil
}

// UsersPage 获取用户列表
func (s *userService) UsersPage(query dto.UserQuery) ([]models.User, int) {
	var model models.User

	if query.Username != "" {
		queryHandle := model.Where("username LIKE ?", "%"+query.Username+"%")
		model = *queryHandle.(*models.User)
	}

	return model.Preload("Roles").Page(query.Page, query.PageSize)
}
