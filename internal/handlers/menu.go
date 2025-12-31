package handlers

import (
	"fmt"
	"webgos/common/convert"
	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/services"
	"webgos/internal/utils/param"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary 创建菜单
// @Description 创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param data body dto.MenuDTO true "菜单参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu [post]
// AddMenu 创建菜单
// @Security BearerAuth
func AddMenu(c *gin.Context) {
	var dtoModel dto.MenuDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建菜单服务
	menuService := services.NewMenuService()

	// 构造菜单模型
	menu := &models.Menu{
		Name:      dtoModel.Name,
		Path:      dtoModel.Path,
		AuthCode:  dtoModel.AuthCode,
		Component: dtoModel.Component,
		Type:      dtoModel.Type,
		Status:    dtoModel.Status,
		Pid:       dtoModel.Pid,
		Meta: models.MenuMeta{
			Title:              dtoModel.Meta.Title,
			Icon:               dtoModel.Meta.Icon,
			AffixTab:           dtoModel.Meta.AffixTab,
			HideChildrenInMenu: dtoModel.Meta.HideChildrenInMenu,
			HideInBreadcrumb:   dtoModel.Meta.HideInBreadcrumb,
			HideInMenu:         dtoModel.Meta.HideInMenu,
			HideInTab:          dtoModel.Meta.HideInTab,
			KeepAlive:          dtoModel.Meta.KeepAlive,
			Order:              dtoModel.Meta.Order,
			Badge:              dtoModel.Meta.Badge,
			BadgeType:          dtoModel.Meta.BadgeType,
			BadgeVariants:      dtoModel.Meta.BadgeVariants,
			IframeSrc:          dtoModel.Meta.IframeSrc,
			Link:               dtoModel.Meta.Link,
		},
	}

	// 创建菜单
	if err := menuService.CreateMenu(menu); err != nil {
		response.Error(c, "创建菜单失败: "+err.Error())
		return
	}

	response.Success(c, "菜单创建成功", menu)
}

// @Summary 编辑菜单
// @Description 编辑菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Param data body dto.MenuDTO true "菜单参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/:id [put]
// EditMenu 编辑菜单
// @Security BearerAuth
func EditMenu(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required,min=1"`
	}
	// 绑定路径参数
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Error(c, "无效的菜单ID: "+err.Error())
		return
	}

	var dtoModel dto.MenuDTO
	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建菜单服务
	menuService := services.NewMenuService()

	menu := &models.Menu{}
	menu, err := menu.Read(uri.ID)
	if err != nil {
		response.Error(c, "菜单不存在: "+err.Error())
		return
	}
	menu.Name = dtoModel.Name
	menu.Path = dtoModel.Path
	menu.AuthCode = dtoModel.AuthCode
	menu.Component = dtoModel.Component
	menu.Type = dtoModel.Type
	menu.Status = dtoModel.Status
	menu.Pid = dtoModel.Pid
	menu.Meta = dtoModel.Meta

	// 更新菜单
	if err := menuService.UpdateMenu(uri.ID, menu); err != nil {
		response.Error(c, "编辑菜单失败: "+err.Error())
		return
	}
	response.Success(c, "编辑菜单成功", nil)
}

// @Summary 删除菜单
// @Description 删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/:id [delete]
// @Security BearerAuth
func DeleteMenu(c *gin.Context) {
	ID := convert.S2Int(c.Param("id"))
	// ID := utils.S2Int(c.Param("id"))
	// 创建菜单服务
	menuService := services.NewMenuService()
	// 删除菜单
	if err := menuService.DeleteMenu(ID); err != nil {
		response.Error(c, "删除菜单失败: "+err.Error())
		return
	}
	response.Success(c, "删除菜单成功", nil)
}

// @Summary 菜单详情
// @Description 获取菜单详情
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/:id [get]
// GetMenu 菜单详情
// @Security BearerAuth
func GetMenu(c *gin.Context) {
	ID := convert.S2Int(c.Param("id"))
	// 创建菜单服务
	menuService := services.NewMenuService()
	// 获取菜单
	menu, err := menuService.GetMenuByID(ID)
	if err != nil {
		response.Error(c, "获取菜单失败: "+err.Error())
		return
	}
	response.Success(c, "获取菜单成功", menu)
}

// @Summary 获取菜单列表
// @Description 获取所有菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/list [get]
// GetMenus 获取菜单列表
// @Security BearerAuth
func GetMenus(c *gin.Context) {
	// 创建菜单服务
	menuService := services.NewMenuService()

	// 获取所有菜单
	menus, err := menuService.GetAllMenus()
	if err != nil {
		response.Error(c, "获取菜单列表失败: "+err.Error())
		return
	}
	response.Success(c, "获取菜单列表成功", menus)
}

// @Summary 获取菜单树
// @Description 获取菜单树结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/tree [get]
// GetMenuTree 获取菜单树
// @Security BearerAuth
func GetMenuTree(c *gin.Context) {
	// 创建菜单服务
	menuService := services.NewMenuService()

	// 获取菜单树
	menuTree, err := menuService.GetMenuTree()
	if err != nil {
		response.Error(c, "获取菜单树失败: "+err.Error())
		return
	}

	response.Success(c, "获取菜单树成功", menuTree)
}

// @Summary 检查菜单名称是否存在
// @Description 检查菜单名称是否存在
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param name query string true "菜单名称"
// @Param id query int false "菜单ID（编辑时用于排除自身）"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/name_exists [get]
// NameExists 检查菜单名称是否存在
// @Security BearerAuth
func NameExists(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.Success(c, "跳过检查", true)
		return
	}

	// 创建菜单服务
	menuService := services.NewMenuService()

	// 检查是否有ID参数（编辑时排除自身）
	var exists bool
	var err error
	id := c.Query("id")
	if id != "" {
		// 转换ID为整数
		var menuID int
		_, err := fmt.Sscanf(id, "%d", &menuID)
		if err != nil {
			response.Error(c, "菜单ID格式错误")
			return
		}

		// 检查名称是否存在（排除指定ID）
		exists, _ = menuService.IsNameExists(name, menuID)
	} else {
		// 检查名称是否存在
		exists, err = menuService.IsNameExists(name)
	}

	if err != nil {
		response.Error(c, "检查菜单名称失败: "+err.Error())
		return
	}

	response.Success(c, "检查成功", exists)
}

// @Summary 检查菜单路径是否存在
// @Description 检查菜单路径是否存在
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param path query string true "菜单路径"
// @Param id query int false "菜单ID（编辑时用于排除自身）"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/menu/path_exists [get]
// PathExists 检查菜单路径是否存在
// @Security BearerAuth
func PathExists(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.Success(c, "跳过检查", true)
		return
	}

	// 创建菜单服务
	menuService := services.NewMenuService()

	// 检查是否有ID参数（编辑时排除自身）
	var exists bool
	var err error
	id := c.Query("id")
	if id != "" {
		// 转换ID为整数
		var menuID int
		_, err := fmt.Sscanf(id, "%d", &menuID)
		if err != nil {
			response.Error(c, "菜单ID格式错误")
			return
		}

		// 检查路径是否存在（排除指定ID）
		exists, _ = menuService.IsPathExists(path, menuID)
	} else {
		// 检查路径是否存在
		exists, err = menuService.IsPathExists(path)
	}

	if err != nil {
		response.Error(c, "检查菜单路径失败: "+err.Error())
		return
	}

	response.Success(c, "检查成功", exists)
}

// @Summary 获取当前用户的菜单树
// @Description 根据当前登录用户的ID，获取其拥有的角色所关联的菜单权限树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.Menu}
// @Failure 400 {object} response.Response
// @Router /api/menu/user_menus [get]
// GetUserMenus 获取当前用户的菜单树
// @Security BearerAuth
func GetUserMenus(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, "未获取到用户信息")
		return
	}

	// 创建菜单服务
	menuService := services.NewMenuService()

	// 获取用户菜单
	menus, err := menuService.GetUserMenus(userID.(int))
	if err != nil {
		response.Error(c, "获取用户菜单失败: "+err.Error())
		return
	}

	response.Success(c, "获取用户菜单成功", menus)
}
