package services

import (
	"errors"

	"webgos/internal/xdb"
	"webgos/internal/dto"
	"webgos/internal/models"

	"gorm.io/gorm"
)

type RBACService interface {
	AddRole(dtoModel dto.AddRoleDTO) (*models.RBACRole, error)
	EditRole(dtoModel dto.EditRoleDTO) error
	AssignRolesToUser(userID int, roleIDs []int) error
	AssignPermissionsToRole(roleID int, permissionIDs []int) error
	GetRoleByID(id int) (*models.RBACRole, error)
	GetUserRoles(userID int) ([]models.RBACRole, error)
	GetRoles() ([]models.RBACRole, error)
	GetPermissions() ([]models.RBACPermission, error)
	GetRolePermissions(roleID int) ([]models.RBACPermission, error)
	DeletePermission(id int) error
}

type rbacService struct{}

func NewRBACService() RBACService {
	return &rbacService{}
}

func (s *rbacService) AddRole(dtoModel dto.AddRoleDTO) (*models.RBACRole, error) {
	role := &models.RBACRole{
		Name:   dtoModel.Name,
		Remark: dtoModel.Remark,
		Status: dtoModel.Status,
	}
	role.SetMenuIDs(dtoModel.Menus)

	if err := xdb.GetDB().Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (s *rbacService) EditRole(dtoModel dto.EditRoleDTO) error {
	var role models.RBACRole
	if err := xdb.GetDB().First(&role, dtoModel.ID).Error; err != nil {
		return err
	}

	if dtoModel.Menus != nil {
		role.SetMenuIDs(dtoModel.Menus)
	}
	if dtoModel.Name != nil {
		role.Name = *dtoModel.Name
	}
	if dtoModel.Remark != nil {
		role.Remark = *dtoModel.Remark
	}
	if dtoModel.Status != nil {
		role.Status = *dtoModel.Status
	}

	return xdb.GetDB().Select("*").Updates(&role).Error
}

func (s *rbacService) AssignRolesToUser(userID int, roleIDs []int) error {
	var user models.User
	if err := xdb.GetDB().First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var roles []models.RBACRole
	if err := xdb.GetDB().Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return errors.New("查询角色时出错")
	}

	if len(roles) != len(roleIDs) {
		return errors.New("部分角色不存在")
	}

	return xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).Association("Roles").Replace(roles)
	})
}

func (s *rbacService) AssignPermissionsToRole(roleID int, permissionIDs []int) error {
	var role models.RBACRole
	if err := xdb.GetDB().First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	var permissions []models.RBACPermission
	if err := xdb.GetDB().Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return errors.New("查询权限时出错")
	}

	if len(permissions) != len(permissionIDs) {
		return errors.New("部分权限不存在")
	}

	return xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&role).Association("Permissions").Replace(permissions)
	})
}

func (s *rbacService) GetRoleByID(id int) (*models.RBACRole, error) {
	var role models.RBACRole
	if err := xdb.GetDB().First(&role, id).Error; err != nil {
		return nil, err
	}

	role.Menus = role.GetMenuIDs()
	return &role, nil
}

func (s *rbacService) GetUserRoles(userID int) ([]models.RBACRole, error) {
	var user models.User
	if err := xdb.GetDB().Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	for i := range user.Roles {
		user.Roles[i].Menus = user.Roles[i].GetMenuIDs()
	}

	return user.Roles, nil
}

func (s *rbacService) GetRoles() ([]models.RBACRole, error) {
	var roles []models.RBACRole
	if err := xdb.GetDB().Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	for i := range roles {
		roles[i].Menus = roles[i].GetMenuIDs()

		permissionIDs := make([]int, len(roles[i].Permissions))
		for j, perm := range roles[i].Permissions {
			permissionIDs[j] = perm.ID
		}
		roles[i].PermissionIDs = permissionIDs
		roles[i].Permissions = nil
	}

	return roles, nil
}

func (s *rbacService) GetPermissions() ([]models.RBACPermission, error) {
	var permissions []models.RBACPermission
	err := xdb.GetDB().Find(&permissions).Error
	return permissions, err
}

func (s *rbacService) GetRolePermissions(roleID int) ([]models.RBACPermission, error) {
	var role models.RBACRole
	if err := xdb.GetDB().Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (s *rbacService) DeletePermission(id int) error {
	var permission models.RBACPermission
	if err := xdb.GetDB().First(&permission, id).Error; err != nil {
		return errors.New("权限不存在")
	}

	return xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rbac_permission_id = ?", id).Delete(&models.RBACRolePermission{}).Error; err != nil {
			return err
		}

		return tx.Delete(&permission, id).Error
	})
}
