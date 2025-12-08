package services

import (
	"errors"

	"webgos/internal/database"
	"webgos/internal/models"

	"gorm.io/gorm"
)

// RBACService RBAC服务接口
type RBACService interface {
	AssignRolesToUser(userID int, roleIDs []int) error
	AssignPermissionsToRole(roleID int, permissionIDs []int) error
	GetRoleByID(id int) (*models.RBACRole, error)
	GetUserRoles(userID int) ([]models.RBACRole, error)
	GetRoles() ([]models.RBACRole, error)
	GetPermissions() ([]models.RBACPermission, error)
	GetRolePermissions(roleID int) ([]models.RBACPermission, error)
	DeletePermission(id int) error
}

// rbacService 实现 RBACService 接口
type rbacService struct{}

// NewRBACService 创建RBAC服务实例
func NewRBACService() RBACService {
	return &rbacService{}
}

// AssignRolesToUser 分配角色给用户
func (s *rbacService) AssignRolesToUser(userID int, roleIDs []int) error {
	// 检查用户是否存在
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在
	var roles []models.RBACRole
	if err := database.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return errors.New("查询角色时出错")
	}

	if len(roles) != len(roleIDs) {
		return errors.New("部分角色不存在")
	}

	// 分配角色给用户
	if err := database.DB.Model(&user).Association("Roles").Replace(&roles); err != nil {
		return err
	}

	return nil
}

// AssignPermissionsToRole 分配权限给角色
func (s *rbacService) AssignPermissionsToRole(roleID int, permissionIDs []int) error {
	// 检查角色是否存在
	var role models.RBACRole
	if err := database.DB.First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	// 检查权限是否存在
	var permissions []models.RBACPermission
	if err := database.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return errors.New("查询权限时出错")
	}

	if len(permissions) != len(permissionIDs) {
		return errors.New("部分权限不存在")
	}

	// 分配权限给角色
	if err := database.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		return err
	}

	return nil
}

// GetRoleByID 根据ID获取角色
func (s *rbacService) GetRoleByID(id int) (*models.RBACRole, error) {
	var role models.RBACRole
	if err := database.DB.First(&role, id).Error; err != nil {
		return nil, err
	}

	// 获取角色的菜单ID列表
	role.Menus = role.GetMenuIDs()

	return &role, nil
}

// GetUserRoles 获取用户的角色
func (s *rbacService) GetUserRoles(userID int) ([]models.RBACRole, error) {
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 为每个角色加载菜单ID列表
	for i := range user.Roles {
		user.Roles[i].Menus = user.Roles[i].GetMenuIDs()
	}

	return user.Roles, nil
}

// GetRoles 获取所有角色列表
func (s *rbacService) GetRoles() ([]models.RBACRole, error) {

	model := models.RBACRole{}
	roles, err := model.Preload("Permissions").More()
	if err != nil {
		return nil, err
	}

	// 为每个角色加载菜单ID列表和权限ID列表
	for i := range roles {
		roles[i].Menus = roles[i].GetMenuIDs()

		// 获取权限ID列表
		permissionIDs := make([]int, len(roles[i].Permissions))
		for j, perm := range roles[i].Permissions {
			permissionIDs[j] = perm.ID
		}
		roles[i].PermissionIDs = permissionIDs
		roles[i].Permissions = nil // 清空权限详细信息，避免冗余数据返回
	}

	return roles, nil
}

// GetPermissions 获取所有权限项列表
func (s *rbacService) GetPermissions() ([]models.RBACPermission, error) {
	var permissions []models.RBACPermission
	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetRolePermissions 获取角色的所有权限
func (s *rbacService) GetRolePermissions(roleID int) ([]models.RBACPermission, error) {
	var role models.RBACRole
	if err := database.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// DeletePermission 删除权限点
func (s *rbacService) DeletePermission(id int) error {
	permission := &models.RBACPermission{}
	permission.ID = id

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联数据
		if err := tx.Where("rbac_permission_id = ?", id).Delete(&models.RBACRolePermission{}).Error; err != nil {
			return err
		}

		// 删除权限数据
		if err := tx.Delete(permission).Error; err != nil {
			return err
		}

		return nil
	})
}
