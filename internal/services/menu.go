package services

import (
	"errors"

	"webgos/internal/config"
	"webgos/internal/xdb"
	"webgos/internal/dto"
	"webgos/internal/models"
)

type MenuService interface {
	AddMenu(dtoModel dto.MenuDTO) (*models.Menu, error)
	EditMenu(id int, dtoModel dto.MenuDTO) error
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

type menuService struct{}

func NewMenuService() MenuService {
	return &menuService{}
}

func (s *menuService) AddMenu(dtoModel dto.MenuDTO) (*models.Menu, error) {
	menu := &models.Menu{
		Name:      dtoModel.Name,
		Path:      dtoModel.Path,
		AuthCode:  dtoModel.AuthCode,
		Component: dtoModel.Component,
		Type:      dtoModel.Type,
		Status:    dtoModel.Status,
		Pid:       dtoModel.Pid,
		Meta:      dtoModel.Meta,
	}

	if err := xdb.GetDB().Create(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *menuService) EditMenu(id int, dtoModel dto.MenuDTO) error {
	var menu models.Menu
	if err := xdb.GetDB().First(&menu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	menu.Name = dtoModel.Name
	menu.Path = dtoModel.Path
	menu.AuthCode = dtoModel.AuthCode
	menu.Component = dtoModel.Component
	menu.Type = dtoModel.Type
	menu.Status = dtoModel.Status
	menu.Pid = dtoModel.Pid
	menu.Meta = dtoModel.Meta

	return xdb.GetDB().Select("*").Updates(&menu).Error
}

func (s *menuService) CreateMenu(menu *models.Menu) error {
	return xdb.GetDB().Create(menu).Error
}

func (s *menuService) UpdateMenu(id int, menu *models.Menu) error {
	var existingMenu models.Menu
	if err := xdb.GetDB().First(&existingMenu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	menu.ID = id
	return xdb.GetDB().Select("*").Updates(menu).Error
}

func (s *menuService) DeleteMenu(id int) error {
	var existingMenu models.Menu
	if err := xdb.GetDB().First(&existingMenu, id).Error; err != nil {
		return errors.New("菜单不存在")
	}

	var childCount int64
	if err := xdb.GetDB().Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	return xdb.GetDB().Delete(&models.Menu{}, id).Error
}

func (s *menuService) GetMenuByID(id int) (*models.Menu, error) {
	var menu models.Menu
	err := xdb.GetDB().First(&menu, id).Error
	return &menu, err
}

func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	var menus []models.Menu
	err := xdb.GetDB().Find(&menus).Error
	return menus, err
}

func (s *menuService) GetMenuTree() ([]models.Menu, error) {
	var menus []models.Menu
	if err := xdb.GetDB().Order("`order` ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	menuTree := buildMenuTree(menus, 0)
	return menuTree, nil
}

func (s *menuService) IsNameExists(name string, id ...int) (bool, error) {
	var count int64
	db := xdb.GetDB().Model(&models.Menu{}).Where("name = ?", name)

	if len(id) > 0 {
		db = db.Where("id != ?", id[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *menuService) IsPathExists(path string, id ...int) (bool, error) {
	var count int64
	db := xdb.GetDB().Model(&models.Menu{}).Where("path = ?", path)

	if len(id) > 0 {
		db = db.Where("id != ?", id[0])
	}

	if err := db.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *menuService) GetUserMenus(userID int) ([]models.Menu, error) {
	var user models.User
	if err := xdb.GetDB().Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	isSuper := false
	if user.Username == config.GlobalConfig.SuperAccount {
		isSuper = true
	}

	menuIDMap := make(map[int]bool)
	for _, role := range user.Roles {
		menuIDs := role.GetMenuIDs()
		for _, menuID := range menuIDs {
			menuIDMap[menuID] = true
		}
	}

	if !isSuper && len(menuIDMap) == 0 {
		return []models.Menu{}, nil
	}

	menuIDs := make([]int, 0, len(menuIDMap))
	for menuID := range menuIDMap {
		menuIDs = append(menuIDs, menuID)
	}

	db := xdb.GetDB().Model(&models.Menu{}).Where("status = ? AND type != ?", 1, "button")
	if !isSuper {
		db = db.Where("id IN ?", menuIDs)
	}

	var menus []models.Menu
	if err := db.Order("`order` ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	menuTree := buildMenuTree(menus, 0)
	setRedirectForMenus(menuTree)
	return menuTree, nil
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
