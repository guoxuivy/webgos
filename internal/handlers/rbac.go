package handlers

import (
	"webgos/internal/dto"
	"webgos/internal/services"
	"webgos/internal/utils/param"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// AddRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.AddRoleDTO true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/role [post]
// @Security BearerAuth
func AddRole(c *gin.Context) {
	var dtoModel dto.AddRoleDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	role, err := rbacService.AddRole(dtoModel)
	if err != nil {
		response.Error(c, "创建角色失败: "+err.Error())
		return
	}

	response.Success(c, "角色创建成功", role)
}

// EditRole 编辑角色
// @Summary 编辑角色
// @Description 更新角色信息
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.EditRoleDTO true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/edit_role [post]
// @Security BearerAuth
func EditRole(c *gin.Context) {
	var dtoModel dto.EditRoleDTO
	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	err := rbacService.EditRole(dtoModel)
	if err != nil {
		response.Error(c, "编辑角色失败: "+err.Error())
		return
	}
	response.Success(c, "编辑角色成功", nil)
}

// AssignRoles 分配角色
// @Summary 分配角色给用户
// @Description 给用户分配角色
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.AssignRolesDTO true "角色分配信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/assign_roles [post]
// @Security BearerAuth
func AssignRoles(c *gin.Context) {
	var dtoModel dto.AssignRolesDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	err := rbacService.AssignRolesToUser(dtoModel.UserID, dtoModel.RoleIDs)
	if err != nil {
		response.Error(c, "分配角色失败: "+err.Error())
		return
	}

	response.Success(c, "角色分配成功", nil)
}

// AssignPermissions 分配权限
// @Summary 分配权限给角色
// @Description 给角色分配权限
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.AssignPermissionsDTO true "权限分配信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/assign_permissions [post]
// @Security BearerAuth
func AssignPermissions(c *gin.Context) {
	var dtoModel dto.AssignPermissionsDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	err := rbacService.AssignPermissionsToRole(dtoModel.RoleID, dtoModel.PermissionIDs)
	if err != nil {
		response.Error(c, "分配权限失败: "+err.Error())
		return
	}

	response.Success(c, "权限分配成功", nil)
}

// GetRoleByID 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.GetRoleDTO true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/role/{id} [get]
// @Security BearerAuth
func GetRoleByID(c *gin.Context) {
	var dtoModel dto.GetRoleDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	role, err := rbacService.GetRoleByID(dtoModel.ID)
	if err != nil {
		response.Error(c, "获取角色失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色成功", role)
}

// GetUserRoles 获取用户角色
// @Summary 获取用户的角色列表
// @Description 获取指定用户的角色列表
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.GetUserRolesDTO true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/user_roles/{id} [get]
// @Security BearerAuth
func GetUserRoles(c *gin.Context) {
	var dtoModel dto.GetUserRolesDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

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
// @Tags 角色权限
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/roles [get]
// @Security BearerAuth
func GetRoles(c *gin.Context) {
	rbacService := services.NewRBACService()
	roles, err := rbacService.GetRoles()
	if err != nil {
		response.Error(c, "获取角色列表失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色列表成功", gin.H{"items": roles, "total": len(roles)})
}

// GetPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 获取所有权限列表
// @Tags 角色权限
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/permissions [get]
// @Security BearerAuth
func GetPermissions(c *gin.Context) {
	rbacService := services.NewRBACService()
	permissions, err := rbacService.GetPermissions()
	if err != nil {
		response.Error(c, "获取权限项列表失败: "+err.Error())
		return
	}
	response.Success(c, "获取权限项列表成功", permissions)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除指定权限
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/permission/{id} [delete]
// @Security BearerAuth
func DeletePermission(c *gin.Context) {
	var dtoModel dto.DeletePermissionDTO
	if err := param.ValidateUri(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}
	rbacService := services.NewRBACService()
	err := rbacService.DeletePermission(dtoModel.ID)
	if err != nil {
		response.Error(c, "删除权限失败: "+err.Error())
		return
	}
	response.Success(c, "删除权限成功", nil)
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色的权限列表
// @Description 获取指定角色的权限列表
// @Tags 角色权限
// @Accept json
// @Produce json
// @Param body body dto.GetRolePermissionsDTO true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/rbac/role_permissions/{id} [get]
// @Security BearerAuth
func GetRolePermissions(c *gin.Context) {
	var dtoModel dto.GetRolePermissionsDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	rbacService := services.NewRBACService()
	permissions, err := rbacService.GetRolePermissions(dtoModel.RoleID)
	if err != nil {
		response.Error(c, "获取角色权限列表失败: "+err.Error())
		return
	}

	response.Success(c, "获取角色权限列表成功", permissions)
}