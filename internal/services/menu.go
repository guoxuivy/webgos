package services

import (
	"errors"

	"webgos/internal/database"
	"webgos/internal/models"
)

// MenuService 菜单服务接口
type MenuService interface {
	CreateMenu(menu *models.Menu) error
	UpdateMenu(id int, menu *models.Menu) error
	DeleteMenu(id int) error
	GetMenuByID(id int) (*models.Menu, error)
	GetAllMenus() ([]models.Menu, error)
	GetMenuTree() ([]models.Menu, error)
	IsNameExists(name string, id ...int) (bool, error)
	IsPathExists(path string, id ...int) (bool, error)
	GetUserMenus(userID int) ([]models.Menu, error)
}

// menuService 实现 MenuService 接口
type menuService struct{}

// NewMenuService 创建菜单服务实例
func NewMenuService() MenuService {
	return &menuService{}
}

// CreateMenu 创建菜单
func (s *menuService) CreateMenu(menu *models.Menu) error {
	if err := menu.Create(menu); err != nil {
		return err
	}
	return nil
}

// UpdateMenu 更新菜单
func (s *menuService) UpdateMenu(id int, menu *models.Menu) error {
	existingMenu, err := menu.Read(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 更新菜单信息
	//  Select("*") 来选择所有字段进行更新
	menu.ID = id
	if err := existingMenu.Select("*").Update(menu); err != nil {
		return err
	}
	return nil
}

// DeleteMenu 删除菜单
func (s *menuService) DeleteMenu(id int) error {
	menu := &models.Menu{}
	existingMenu, err := menu.Read(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否有子菜单
	var childCount int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	// 删除菜单
	if err := existingMenu.Delete(id); err != nil {
		return err
	}
	return nil
}

// GetMenuByID 根据ID获取菜单
func (s *menuService) GetMenuByID(id int) (*models.Menu, error) {
	menu := &models.Menu{}
	existingMenu, err := menu.Read(id)
	if err != nil {
		return nil, err
	}
	return existingMenu, nil
}

// GetAllMenus 获取所有菜单
func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	menu := &models.Menu{}
	menus, err := menu.More()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// GetMenuTree 获取菜单树
func (s *menuService) GetMenuTree() ([]models.Menu, error) {
	menu := &models.Menu{}
	menus, err := menu.More()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuTree := buildMenuTree(menus, 0)
	return menuTree, nil
}

// IsNameExists 检查菜单名称是否存在
func (s *menuService) IsNameExists(name string, id ...int) (bool, error) {
	menu := &models.Menu{}

	// 如果提供了ID参数，则排除该ID进行检查
	if len(id) > 0 {
		count, err := menu.Where("name = ? AND id != ?", name, id[0]).Count()
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	// 默认检查逻辑
	count, err := menu.Where("name = ?", name).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsPathExists 检查菜单路径是否存在
func (s *menuService) IsPathExists(path string, id ...int) (bool, error) {
	menu := &models.Menu{}

	// 如果提供了ID参数，则排除该ID进行检查
	if len(id) > 0 {
		count, err := menu.Where("path = ? AND id != ?", path, id[0]).Count()
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	// 默认检查逻辑
	count, err := menu.Where("path = ?", path).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserMenus 获取用户菜单
func (s *menuService) GetUserMenus(userID int) ([]models.Menu, error) {
	// 获取用户及其角色信息
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	// 检查是否为超管角色(拥有所有权限)
	isAdmin := false
	// 收集用户所有角色的菜单ID
	menuIDMap := make(map[int]bool)
	for _, role := range user.Roles {
		if role.Name == "Super" {
			isAdmin = true
		}
		// 获取角色的菜单ID列表
		menuIDs := role.GetMenuIDs()
		for _, menuID := range menuIDs {
			menuIDMap[menuID] = true
		}
	}

	// 如果没有关联的菜单ID，则返回空列表
	if !isAdmin && len(menuIDMap) == 0 {
		return []models.Menu{}, nil
	}

	// 将map中的key转换为slice
	menuIDs := make([]int, 0, len(menuIDMap))
	for menuID := range menuIDMap {
		menuIDs = append(menuIDs, menuID)
	}

	// 根据菜单ID获取菜单列表
	var menus []models.Menu
	// button类型的菜单不包含在内，因为它们通常不显示在菜单树中
	query := database.DB.Where("status = ? AND type != ?", 1, "button")
	if !isAdmin {
		query = query.Where("id IN ?", menuIDs)
	}
	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}
	// 构建菜单树
	menuTree := buildMenuTree(menus, 0)
	// 为有子菜单的菜单项设置重定向路径
	setRedirectForMenus(menuTree)
	return menuTree, nil
}

// setRedirectForMenus 为有子菜单的菜单项设置重定向路径
func setRedirectForMenus(menus []models.Menu) {
	for i := range menus {
		// 如果有子菜单且Redirect为空，则设置为第一个子菜单的路径
		if len(menus[i].Children) > 0 && menus[i].Redirect == "" {
			menus[i].Redirect = menus[i].Children[0].Path
		}

		// 递归处理子菜单
		setRedirectForMenus(menus[i].Children)
	}
}

// buildMenuTree 构建菜单树
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
