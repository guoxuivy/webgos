package handlers

import (
	"hserp/internal/dto"
	"hserp/internal/models"
	"hserp/internal/services"
	"hserp/internal/utils"
	"hserp/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary 创建角色
// @Description 创建新角色
// @Tags RBAC
// @Accept json
// @Produce json
// @Param data body dto.AddRoleDTO true "角色参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/add/role [post]
// AddRole 创建角色
// @Security BearerAuth
func AddRole(c *gin.Context) {
	var dtoModel dto.AddRoleDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	role := &models.RBACRole{
		Name:   dtoModel.Name,
		Remark: dtoModel.Remark,
		Status: dtoModel.Status,
	}
	role.SetMenuIDs(dtoModel.Menus)
	role.Create(role)

	// // 创建角色服务
	// rbacService := services.NewRBACService()
	// role, err := rbacService.CreateRole(dtoModel.Name, dtoModel.Remark, dtoModel.Menus)
	// if err != nil {
	// 	response.Error(c, "创建角色失败: "+err.Error())
	// 	return
	// }

	// // 更新角色状态
	// role.Status = dtoModel.Status
	// database.DB.Save(role)

	response.Success(c, "角色创建成功", role)
}

// @Summary 编辑角色
// @Description 编辑角色信息
// @Tags RBAC
// @Accept json
// @Produce json
// @Param data body dto.EditRoleDTO true "编辑参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/edit/role [POST]
// EditRole 编辑角色
// @Security BearerAuth
func EditRole(c *gin.Context) {
	var dtoModel dto.EditRoleDTO
	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	role := &models.RBACRole{}
	role, err := role.Read(dtoModel.ID)
	if err != nil {
		response.Error(c, "角色不存在: "+err.Error())
		return
	}

	role.SetMenuIDs(dtoModel.Menus)
	role.Name = dtoModel.Name
	role.Remark = dtoModel.Remark
	role.Status = dtoModel.Status

	err = role.Select("*").Update(role)
	if err != nil {
		response.Error(c, "更新角色失败: "+err.Error())
		return
	}

	// // 创建角色服务
	// rbacService := services.NewRBACService()
	// role, err := rbacService.GetRoleByID(dtoModel.ID)
	// if err != nil {
	// 	response.Error(c, "角色不存在: "+err.Error())
	// 	return
	// }

	// // 更新角色信息
	// role.Name = dtoModel.Name
	// role.Remark = dtoModel.Remark
	// role.Status = dtoModel.Status
	// err = database.DB.Save(role).Error
	// if err != nil {
	// 	response.Error(c, "编辑角色失败: "+err.Error())
	// 	return
	// }
	response.Success(c, "编辑角色成功", nil)
}

// @Summary 分配角色给用户
// @Description 给用户分配角色
// @Tags RBAC
// @Accept json
// @Produce json
// @Param data body dto.AssignRolesDTO true "分配参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/assign/roles [post]
// AssignRoles 分配角色给用户
// @Security BearerAuth
func AssignRoles(c *gin.Context) {
	var dtoModel dto.AssignRolesDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建角色服务
	rbacService := services.NewRBACService()
	err := rbacService.AssignRolesToUser(dtoModel.UserID, dtoModel.RoleIDs)
	if err != nil {
		response.Error(c, "分配角色失败: "+err.Error())
		return
	}

	response.Success(c, "角色分配成功", nil)
}

// @Summary 分配权限给角色
// @Description 给角色分配权限
// @Tags RBAC
// @Accept json
// @Produce json
// @Param data body dto.AssignPermissionsDTO true "分配参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/assign/permissions [post]
// AssignPermissions 分配权限给角色
// @Security BearerAuth
func AssignPermissions(c *gin.Context) {
	var dtoModel dto.AssignPermissionsDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建角色服务
	rbacService := services.NewRBACService()
	err := rbacService.AssignPermissionsToRole(dtoModel.RoleID, dtoModel.PermissionIDs)
	if err != nil {
		response.Error(c, "分配权限失败: "+err.Error())
		return
	}

	response.Success(c, "权限分配成功", nil)
}

// @Summary 获取角色详情
// @Description 获取指定角色详情
// @Tags RBAC
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/role/{id} [get]
// GetRoleByID 根据ID获取角色
// @Security BearerAuth
func GetRoleByID(c *gin.Context) {
	var dtoModel dto.GetRoleDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建角色服务
	rbacService := services.NewRBACService()
	role, err := rbacService.GetRoleByID(dtoModel.ID)
	if err != nil {
		response.Error(c, "获取角色失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色成功", role)
}

// @Summary 获取用户角色
// @Description 获取指定用户的角色列表
// @Tags RBAC
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {array} models.RBACRole
// @Failure 400 {object} response.Response
// @Router /api/rbac/user/{id}/roles [get]
// GetUserRoles 获取用户的角色
// @Security BearerAuth
func GetUserRoles(c *gin.Context) {
	var dtoModel dto.GetUserRolesDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建角色服务
	rbacService := services.NewRBACService()
	roles, err := rbacService.GetUserRoles(dtoModel.UserID)
	if err != nil {
		response.Error(c, "获取用户角色失败: "+err.Error())
		return
	}

	response.Success(c, "获取用户角色成功", roles)
}

// GetRoles 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色列表
// @Tags RBAC
// @Produce json
// @Success 200 {array} models.RBACRole
// @Failure 400 {object} response.Response
// @Router /api/rbac/roles [get]
// @Security BearerAuth
func GetRoles(c *gin.Context) {
	// 创建角色服务
	rbacService := services.NewRBACService()
	roles, err := rbacService.GetRoles()
	if err != nil {
		response.Error(c, "获取角色列表失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色列表成功", gin.H{"items": roles, "total": 4})
}

// GetPermissions 获取权限项列表
// @Summary 获取权限项列表
// @Description 获取所有权限项列表
// @Tags RBAC
// @Produce json
// @Success 200 {array} models.RBACPermission
// @Failure 400 {object} response.Response
// @Router /api/rbac/permissions [get]
// @Security BearerAuth
func GetPermissions(c *gin.Context) {
	// 创建角色服务
	rbacService := services.NewRBACService()
	permissions, err := rbacService.GetPermissions()
	if err != nil {
		response.Error(c, "获取权限项列表失败: "+err.Error())
		return
	}
	response.Success(c, "获取权限项列表成功", permissions)
}

// GetRolePermissions 获取角色权限列表
// @Summary 获取角色权限列表
// @Description 获取指定角色的所有权限列表
// @Tags RBAC
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {array} models.RBACPermission
// @Failure 400 {object} response.Response
// @Router /api/rbac/role/{id}/permissions [get]
// @Security BearerAuth
func GetRolePermissions(c *gin.Context) {
	var dtoModel dto.GetRolePermissionsDTO

	if err := utils.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	// 创建角色服务
	rbacService := services.NewRBACService()
	permissions, err := rbacService.GetRolePermissions(dtoModel.RoleID)
	if err != nil {
		response.Error(c, "获取角色权限列表失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色权限列表成功", permissions)
}
