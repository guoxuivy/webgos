package services

import (
	"errors"

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

// AssignRolesToUser 分配角色给用户（使用 BaseModel）
func (s *rbacService) AssignRolesToUser(userID int, roleIDs []int) error {
	var userModel models.User

	// 检查用户是否存在（使用 BaseModel）
	user, err := userModel.Read(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在（使用 BaseModel）
	var roleModel models.RBACRole
	roles, err := roleModel.Where("id IN ?", roleIDs).More()
	if err != nil {
		return errors.New("查询角色时出错")
	}

	if len(roles) != len(roleIDs) {
		return errors.New("部分角色不存在")
	}

	// 使用事务更新用户角色关联
	return userModel.Transaction(func(tx *gorm.DB) error {
		// 清空用户的现有角色
		if err := tx.Model(user).Association("Roles").Clear(); err != nil {
			return err
		}

		// 添加新的角色
		if len(roles) > 0 {
			if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// AssignPermissionsToRole 分配权限给角色（使用 BaseModel）
func (s *rbacService) AssignPermissionsToRole(roleID int, permissionIDs []int) error {
	var roleModel models.RBACRole

	// 检查角色是否存在（使用 BaseModel）
	role, err := roleModel.Read(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查权限是否存在（使用 BaseModel）
	var permissionModel models.RBACPermission
	permissions, err := permissionModel.Where("id IN ?", permissionIDs).More()
	if err != nil {
		return errors.New("查询权限时出错")
	}

	if len(permissions) != len(permissionIDs) {
		return errors.New("部分权限不存在")
	}

	// 使用事务更新角色权限关联
	return roleModel.Transaction(func(tx *gorm.DB) error {
		// 清空角色的现有权限
		if err := tx.Model(role).Association("Permissions").Clear(); err != nil {
			return err
		}

		// 添加新的权限
		if len(permissions) > 0 {
			if err := tx.Model(role).Association("Permissions").Append(permissions); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetRoleByID 根据ID获取角色（使用 BaseModel）
func (s *rbacService) GetRoleByID(id int) (*models.RBACRole, error) {
	var roleModel models.RBACRole
	role, err := roleModel.Read(id)
	if err != nil {
		return nil, err
	}

	// 获取角色的菜单ID列表
	role.Menus = role.GetMenuIDs()

	return role, nil
}

// GetUserRoles 获取用户的角色（使用 BaseModel）
func (s *rbacService) GetUserRoles(userID int) ([]models.RBACRole, error) {
	var userModel models.User
	user, err := userModel.Preload("Roles").Read(userID)
	if err != nil {
		return nil, err
	}

	// 为每个角色加载菜单ID列表
	for i := range user.Roles {
		user.Roles[i].Menus = user.Roles[i].GetMenuIDs()
	}

	return user.Roles, nil
}

// GetRoles 获取所有角色列表（使用 BaseModel）
func (s *rbacService) GetRoles() ([]models.RBACRole, error) {
	var roleModel models.RBACRole
	roles, err := roleModel.Preload("Permissions").More()
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

// GetPermissions 获取所有权限项列表（使用 BaseModel）
func (s *rbacService) GetPermissions() ([]models.RBACPermission, error) {
	var permissionModel models.RBACPermission
	return permissionModel.More()
}

// GetRolePermissions 获取角色的所有权限（使用 BaseModel）
func (s *rbacService) GetRolePermissions(roleID int) ([]models.RBACPermission, error) {
	var roleModel models.RBACRole
	role, err := roleModel.Preload("Permissions").Read(roleID)
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// DeletePermission 删除权限点（使用 BaseModel）
func (s *rbacService) DeletePermission(id int) error {
	var permissionModel models.RBACPermission

	// 检查权限是否存在
	permission, err := permissionModel.Read(id)
	if err != nil {
		return errors.New("权限不存在")
	}

	return permissionModel.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联数据
		if err := tx.Where("rbac_permission_id = ?", id).Delete(&models.RBACRolePermission{}).Error; err != nil {
			return err
		}

		// 删除权限数据（使用 BaseModel）
		return permission.WithTx(tx).Delete(id)
	})
}
