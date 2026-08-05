package services

import (
	"errors"
	"strconv"

	"webgos/internal/cache"
	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xdb"

	"gorm.io/gorm"
)

type RBACService interface {
	AddRole(dtoModel dto.AddRoleDTO) (*models.RBACRole, error)
	EditRole(dtoModel dto.EditRoleDTO) error
	AssignRolesToUser(userID int, roleIDs []int) error
	AssignMenusToRole(roleID int, menuIDs []int) error
	GetRoleByID(id int) (*models.RBACRole, error)
	GetUserRoles(userID int) ([]models.RBACRole, error)
	GetRoles() ([]models.RBACRole, error)
	GetPermissions() ([]models.RBACPermission, error)
	GetRolePermissions(roleID int) ([]models.RBACPermission, error)
	GetMenuPermissions(menuID int) ([]models.RBACPermission, error)
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

	if err := xdb.GetDB().Create(role).Error; err != nil {
		return nil, err
	}

	// 绑定菜单（多对多）
	if len(dtoModel.MenuIDs) > 0 {
		if err := s.AssignMenusToRole(role.ID, dtoModel.MenuIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (s *rbacService) EditRole(dtoModel dto.EditRoleDTO) error {
	var role models.RBACRole
	if err := xdb.GetDB().First(&role, dtoModel.ID).Error; err != nil {
		return err
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
	if err := xdb.GetDB().Select("*").Updates(&role).Error; err != nil {
		return err
	}

	if dtoModel.MenuIDs != nil {
		if err := s.AssignMenusToRole(role.ID, dtoModel.MenuIDs); err != nil {
			return err
		}
	}
	return nil
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

	if err := xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).Association("Roles").Replace(roles)
	}); err != nil {
		return err
	}

	// 角色变更后失效该用户的权限缓存
	InvalidateUserPermissionCache(userID)
	return nil
}

func (s *rbacService) AssignMenusToRole(roleID int, menuIDs []int) error {
	var role models.RBACRole
	if err := xdb.GetDB().First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	var menus []models.Menu
	if len(menuIDs) > 0 {
		if err := xdb.GetDB().Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return errors.New("查询菜单时出错")
		}
		if len(menus) != len(menuIDs) {
			return errors.New("部分菜单不存在")
		}
	}

	if err := xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&role).Association("Menus").Replace(menus)
	}); err != nil {
		return err
	}

	// 菜单变更后失效拥有该角色的所有用户的权限缓存
	InvalidateRolePermissionCache(roleID)
	return nil
}

func (s *rbacService) GetRoleByID(id int) (*models.RBACRole, error) {
	var role models.RBACRole
	if err := xdb.GetDB().Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}

	role.MenuIDs = menuIDsOf(role.Menus)
	return &role, nil
}

func (s *rbacService) GetUserRoles(userID int) ([]models.RBACRole, error) {
	var user models.User
	if err := xdb.GetDB().Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	for i := range user.Roles {
		user.Roles[i].MenuIDs = menuIDsOf(user.Roles[i].Menus)
	}

	return user.Roles, nil
}

func (s *rbacService) GetRoles() ([]models.RBACRole, error) {
	var roles []models.RBACRole
	if err := xdb.GetDB().Preload("Menus").Find(&roles).Error; err != nil {
		return nil, err
	}

	for i := range roles {
		roles[i].MenuIDs = menuIDsOf(roles[i].Menus)
	}

	return roles, nil
}

// menuIDsOf 从菜单切片中提取 id 列表
func menuIDsOf(menus []models.Menu) []int {
	ids := make([]int, 0, len(menus))
	for _, m := range menus {
		ids = append(ids, m.ID)
	}
	return ids
}

func (s *rbacService) GetPermissions() ([]models.RBACPermission, error) {
	var permissions []models.RBACPermission
	err := xdb.GetDB().Find(&permissions).Error
	return permissions, err
}

func (s *rbacService) GetRolePermissions(roleID int) ([]models.RBACPermission, error) {
	var role models.RBACRole
	if err := xdb.GetDB().Preload("Menus.Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	// 角色权限 = 所绑菜单下所有权限点的去重集合
	permMap := make(map[int]models.RBACPermission)
	for _, menu := range role.Menus {
		for _, perm := range menu.Permissions {
			permMap[perm.ID] = perm
		}
	}
	permissions := make([]models.RBACPermission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

// GetMenuPermissions 获取菜单绑定的权限点列表（多对多）
func (s *rbacService) GetMenuPermissions(menuID int) ([]models.RBACPermission, error) {
	var menu models.Menu
	if err := xdb.GetDB().Preload("Permissions").First(&menu, menuID).Error; err != nil {
		return nil, err
	}
	return menu.Permissions, nil
}

func (s *rbacService) DeletePermission(id int) error {
	var permission models.RBACPermission
	if err := xdb.GetDB().First(&permission, id).Error; err != nil {
		return errors.New("权限不存在")
	}

	if err := xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		// 清除权限点-菜单关联
		if err := tx.Model(&permission).Association("Menus").Clear(); err != nil {
			return err
		}
		return tx.Delete(&permission, id).Error
	}); err != nil {
		return err
	}
	return nil
}

// InvalidateUserPermissionCache 失效指定用户的权限缓存
func InvalidateUserPermissionCache(userID int) {
	cache.GetCache().Delete(cache.PermissionPrefix + ":" + strconv.Itoa(userID))
}

// InvalidateRolePermissionCache 失效拥有指定角色的所有用户的权限缓存
func InvalidateRolePermissionCache(roleID int) {
	var userIDs []int
	xdb.GetDB().Table("rbac_user_roles").Where("rbac_role_id = ?", roleID).Pluck("user_id", &userIDs)
	for _, uid := range userIDs {
		InvalidateUserPermissionCache(uid)
	}
}

// InvalidateMenuPermissionCache 失效绑定了指定菜单的角色下所有用户的权限缓存
func InvalidateMenuPermissionCache(menuID int) {
	var roleIDs []int
	xdb.GetDB().Table("rbac_role_menus").Where("menu_id = ?", menuID).Pluck("rbac_role_id", &roleIDs)
	for _, rid := range roleIDs {
		InvalidateRolePermissionCache(rid)
	}
}
