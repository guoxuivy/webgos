package services

import (
	"context"
	"errors"

	"webgos/internal/config"
	"webgos/internal/dto"
	"webgos/internal/models"

	"gorm.io/gorm"
)

type MenuService interface {
	AddMenu(ctx context.Context, dtoModel dto.MenuDTO) (*models.Menu, error)
	EditMenu(ctx context.Context, id int, dtoModel dto.MenuDTO) error
	CreateMenu(ctx context.Context, menu *models.Menu) error
	UpdateMenu(ctx context.Context, id int, menu *models.Menu) error
	DeleteMenu(ctx context.Context, id int) error
	GetMenuByID(ctx context.Context, id int) (*models.Menu, error)
	GetAllMenus(ctx context.Context) ([]models.Menu, error)
	GetMenuTree(ctx context.Context) ([]models.Menu, error)
	IsNameExists(ctx context.Context, name string, id ...int) (bool, error)
	IsPathExists(ctx context.Context, path string, id ...int) (bool, error)
	GetUserMenus(ctx context.Context, userID int) ([]models.Menu, error)
	AssignPermissionsToMenu(ctx context.Context, menuID int, permissionIDs []int) error
}

type menuService struct{}

func NewMenuService() MenuService {
	return &menuService{}
}

func (s *menuService) AddMenu(ctx context.Context, dtoModel dto.MenuDTO) (*models.Menu, error) {
	menu := &models.Menu{
		Name:      dtoModel.Name,
		Path:      dtoModel.Path,
		Component: dtoModel.Component,
		Type:      dtoModel.Type,
		Status:    dtoModel.Status,
		Pid:       dtoModel.Pid,
		Meta:      dtoModel.Meta,
	}

	if err := ctxDB(ctx).Create(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *menuService) EditMenu(ctx context.Context, id int, dtoModel dto.MenuDTO) error {
	var menu models.Menu
	if err := ctxDB(ctx).First(&menu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	menu.Name = dtoModel.Name
	menu.Path = dtoModel.Path
	menu.Component = dtoModel.Component
	menu.Type = dtoModel.Type
	menu.Status = dtoModel.Status
	menu.Pid = dtoModel.Pid
	menu.Meta = dtoModel.Meta

	return ctxDB(ctx).Select("*").Updates(&menu).Error
}

func (s *menuService) CreateMenu(ctx context.Context, menu *models.Menu) error {
	return ctxDB(ctx).Create(menu).Error
}

func (s *menuService) UpdateMenu(ctx context.Context, id int, menu *models.Menu) error {
	var existingMenu models.Menu
	if err := ctxDB(ctx).First(&existingMenu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	menu.ID = id
	return ctxDB(ctx).Select("*").Updates(menu).Error
}

func (s *menuService) DeleteMenu(ctx context.Context, id int) error {
	var existingMenu models.Menu
	if err := ctxDB(ctx).First(&existingMenu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	var childCount int64
	if err := ctxDB(ctx).Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	return ctxDB(ctx).Delete(&models.Menu{}, id).Error
}

func (s *menuService) GetMenuByID(ctx context.Context, id int) (*models.Menu, error) {
	var menu models.Menu
	err := ctxSDB(ctx).First(&menu, id).Error
	return &menu, err
}

func (s *menuService) GetAllMenus(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu
	err := ctxSDB(ctx).Find(&menus).Error
	return menus, err
}

func (s *menuService) GetMenuTree(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu
	if err := ctxSDB(ctx).Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	menuTree := buildMenuTree(menus, 0)
	return menuTree, nil
}

func (s *menuService) IsNameExists(ctx context.Context, name string, id ...int) (bool, error) {
	var count int64
	db := ctxDB(ctx).Model(&models.Menu{}).Where("name = ?", name)

	if len(id) > 0 {
		db = db.Where("id != ?", id[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *menuService) IsPathExists(ctx context.Context, path string, id ...int) (bool, error) {
	var count int64
	db := ctxDB(ctx).Model(&models.Menu{}).Where("path = ?", path)

	if len(id) > 0 {
		db = db.Where("id != ?", id[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *menuService) GetUserMenus(ctx context.Context, userID int) ([]models.Menu, error) {
	var user models.User
	if err := ctxSDB(ctx).Preload("Roles.Menus").First(&user, userID).Error; err != nil {
		return nil, err
	}

	isSuper := false
	if user.Username == config.GlobalConfig.SuperAccount {
		isSuper = true
	}

	menuIDMap := make(map[int]bool)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			menuIDMap[menu.ID] = true
		}
	}

	if !isSuper && len(menuIDMap) == 0 {
		return []models.Menu{}, nil
	}

	menuIDs := make([]int, 0, len(menuIDMap))
	for menuID := range menuIDMap {
		menuIDs = append(menuIDs, menuID)
	}

	db := ctxSDB(ctx).Model(&models.Menu{}).Where("status = ? AND type != ?", 1, "button")
	if !isSuper {
		db = db.Where("id IN ?", menuIDs)
	}

	var menus []models.Menu
	if err := db.Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	menuTree := buildMenuTree(menus, 0)
	setRedirectForMenus(menuTree)
	return menuTree, nil
}

// AssignPermissionsToMenu 维护菜单-权限多对多关系，支持同一权限绑定多个菜单。
// 采用 Replace 语义，保证幂等。
func (s *menuService) AssignPermissionsToMenu(ctx context.Context, menuID int, permissionIDs []int) error {
	var menu models.Menu
	if err := ctxDB(ctx).First(&menu, menuID).Error; err != nil {
		return errors.New("菜单不存在")
	}

	var permissions []models.RBACPermission
	if len(permissionIDs) > 0 {
		if err := ctxDB(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return errors.New("查询权限时出错")
		}
		if len(permissions) != len(permissionIDs) {
			return errors.New("部分权限不存在")
		}
	}

	if err := ctxDB(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&menu).Association("Permissions").Replace(permissions)
	}); err != nil {
		return err
	}

	// 菜单-权限变更后，失效绑定了该菜单的角色下所有用户的权限缓存
	InvalidateMenuPermissionCache(ctx, menuID)
	return nil
}

func setRedirectForMenus(menus []models.Menu) {
	for i := range menus {
		if len(menus[i].Children) > 0 && menus[i].Redirect == "" {
			menus[i].Redirect = menus[i].Children[0].Path
		}
		setRedirectForMenus(menus[i].Children)
	}
}

func buildMenuTree(menus []models.Menu, parentID int) []models.Menu {
	var tree []models.Menu
	for i := range menus {
		if int(menus[i].Pid) == parentID {
			children := buildMenuTree(menus, int(menus[i].ID))
			menus[i].Children = children
			tree = append(tree, menus[i])
		}
	}
	return tree
}
