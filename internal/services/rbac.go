package services

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"webgos/internal/cache"
	"webgos/internal/dto"
	"webgos/internal/models"

	"gorm.io/gorm"
)

type RBACService interface {
	AddRole(ctx context.Context, dtoModel dto.AddRoleDTO) (*models.RBACRole, error)
	EditRole(ctx context.Context, dtoModel dto.EditRoleDTO) error
	AssignRolesToUser(ctx context.Context, userID int, roleIDs []int) error
	AssignMenusToRole(ctx context.Context, roleID int, menuIDs []int) error
	GetRoleByID(ctx context.Context, id int) (*models.RBACRole, error)
	GetUserRoles(ctx context.Context, userID int) ([]models.RBACRole, error)
	GetRoles(ctx context.Context) ([]models.RBACRole, error)
	GetPermissions(ctx context.Context) ([]PermissionPoint, error)
	GetRolePermissions(ctx context.Context, roleID int) ([]models.MenuPermission, error)
	GetMenuPermissions(ctx context.Context, menuID int) ([]models.MenuPermission, error)
	DeletePermission(ctx context.Context, menuID int, permKey string) error
}

type rbacService struct{}

func NewRBACService() RBACService {
	return &rbacService{}
}

func (s *rbacService) AddRole(ctx context.Context, dtoModel dto.AddRoleDTO) (*models.RBACRole, error) {
	role := &models.RBACRole{
		Name:   dtoModel.Name,
		Remark: dtoModel.Remark,
		Status: dtoModel.Status,
	}

	if err := ctxDB(ctx).Create(role).Error; err != nil {
		return nil, err
	}

	// 绑定菜单（多对多）
	if len(dtoModel.MenuIDs) > 0 {
		if err := s.AssignMenusToRole(ctx, role.ID, dtoModel.MenuIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (s *rbacService) EditRole(ctx context.Context, dtoModel dto.EditRoleDTO) error {
	var role models.RBACRole
	if err := ctxDB(ctx).First(&role, dtoModel.ID).Error; err != nil {
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
	if err := ctxDB(ctx).Select("*").Updates(&role).Error; err != nil {
		return err
	}

	if dtoModel.MenuIDs != nil {
		if err := s.AssignMenusToRole(ctx, role.ID, dtoModel.MenuIDs); err != nil {
			return err
		}
	}
	return nil
}

func (s *rbacService) AssignRolesToUser(ctx context.Context, userID int, roleIDs []int) error {
	var user models.User
	if err := ctxDB(ctx).First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var roles []models.RBACRole
	if err := ctxDB(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return errors.New("查询角色时出错")
	}

	if len(roles) != len(roleIDs) {
		return errors.New("部分角色不存在")
	}

	if err := ctxDB(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).Association("Roles").Replace(roles)
	}); err != nil {
		return err
	}

	// 角色变更后失效该用户的权限缓存
	InvalidateUserPermissionCache(ctx, userID)
	return nil
}

func (s *rbacService) AssignMenusToRole(ctx context.Context, roleID int, menuIDs []int) error {
	var role models.RBACRole
	if err := ctxDB(ctx).First(&role, roleID).Error; err != nil {
		return errors.New("角色不存在")
	}

	var menus []models.Menu
	if len(menuIDs) > 0 {
		if err := ctxDB(ctx).Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return errors.New("查询菜单时出错")
		}
		if len(menus) != len(menuIDs) {
			return errors.New("部分菜单不存在")
		}
	}

	if err := ctxDB(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&role).Association("Menus").Replace(menus)
	}); err != nil {
		return err
	}

	// 菜单变更后失效拥有该角色的所有用户的权限缓存
	InvalidateRolePermissionCache(ctx, roleID)
	return nil
}

func (s *rbacService) GetRoleByID(ctx context.Context, id int) (*models.RBACRole, error) {
	var role models.RBACRole
	if err := ctxSDB(ctx).Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}

	role.MenuIDs = menuIDsOf(role.Menus)
	return &role, nil
}

func (s *rbacService) GetUserRoles(ctx context.Context, userID int) ([]models.RBACRole, error) {
	var user models.User
	if err := ctxSDB(ctx).Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	for i := range user.Roles {
		user.Roles[i].MenuIDs = menuIDsOf(user.Roles[i].Menus)
	}

	return user.Roles, nil
}

func (s *rbacService) GetRoles(ctx context.Context) ([]models.RBACRole, error) {
	var roles []models.RBACRole
	if err := ctxSDB(ctx).Preload("Menus").Find(&roles).Error; err != nil {
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

// PermissionPoint 实时权限点（路由投影），非持久化实体。
type PermissionPoint struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
}

// GetPermissions 实时返回所有权限点（路由投影），不再查库，直接来自内存路由描述表。
func (s *rbacService) GetPermissions(ctx context.Context) ([]PermissionPoint, error) {
	points := make([]PermissionPoint, 0, len(models.RouteDescriptions))
	for key, desc := range models.RouteDescriptions {
		// key = path#method（小写 path + 大写 method）
		path, method := key, ""
		if idx := strings.LastIndex(key, "#"); idx >= 0 {
			path, method = key[:idx], key[idx+1:]
		}
		points = append(points, PermissionPoint{
			Key:         key,
			Name:        key,
			Description: desc,
			Path:        path,
			Method:      method,
		})
	}
	// 按 key 排序，保证前端树渲染稳定
	sort.Slice(points, func(i, j int) bool {
		return points[i].Key < points[j].Key
	})
	return points, nil
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleID int) ([]models.MenuPermission, error) {
	var role models.RBACRole
	if err := ctxSDB(ctx).Preload("Menus.PermissionKeys").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	// 角色权限 = 所绑菜单下所有权限键的去重集合
	permMap := make(map[string]models.MenuPermission)
	for _, menu := range role.Menus {
		for _, perm := range menu.PermissionKeys {
			permMap[perm.PermKey] = perm
		}
	}
	permissions := make([]models.MenuPermission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

// GetMenuPermissions 获取菜单绑定的权限键列表
func (s *rbacService) GetMenuPermissions(ctx context.Context, menuID int) ([]models.MenuPermission, error) {
	var menu models.Menu
	if err := ctxSDB(ctx).Preload("PermissionKeys").First(&menu, menuID).Error; err != nil {
		return nil, err
	}
	return menu.PermissionKeys, nil
}

// DeletePermission 解绑单个菜单权限键（按 menu_id + perm_key 删除 rbac_menu_permissions 行）。
func (s *rbacService) DeletePermission(ctx context.Context, menuID int, permKey string) error {
	if err := ctxDB(ctx).Where("menu_id = ? AND perm_key = ?", menuID, permKey).
		Delete(&models.MenuPermission{}).Error; err != nil {
		return err
	}
	// 菜单-权限变更后，失效绑定了该菜单的角色下所有用户的权限缓存
	InvalidateMenuPermissionCache(ctx, menuID)
	return nil
}

// InvalidateUserPermissionCache 失效指定用户的权限缓存
func InvalidateUserPermissionCache(ctx context.Context, userID int) {
	cache.GetCache().Delete(cache.PermissionPrefix + ":" + strconv.Itoa(userID))
}

// InvalidateRolePermissionCache 失效拥有指定角色的所有用户的权限缓存
func InvalidateRolePermissionCache(ctx context.Context, roleID int) {
	var userIDs []int
	ctxDB(ctx).Table("rbac_user_roles").Where("rbac_role_id = ?", roleID).Pluck("user_id", &userIDs)
	for _, uid := range userIDs {
		InvalidateUserPermissionCache(ctx, uid)
	}
}

// InvalidateMenuPermissionCache 失效绑定了指定菜单的角色下所有用户的权限缓存
func InvalidateMenuPermissionCache(ctx context.Context, menuID int) {
	var roleIDs []int
	ctxDB(ctx).Table("rbac_role_menus").Where("menu_id = ?", menuID).Pluck("rbac_role_id", &roleIDs)
	for _, rid := range roleIDs {
		InvalidateRolePermissionCache(ctx, rid)
	}
}
