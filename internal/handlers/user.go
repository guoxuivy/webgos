package handlers

import (
	"webgos/internal/dto"
	"webgos/internal/services"
	"webgos/internal/utils/param"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// UserInfo 获取用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response
// @Router /api/user/info [get]
// @Security BearerAuth
func UserInfo(c *gin.Context) {
	userService := services.NewUserService()
	user, err := userService.GetUserInfo(c.GetInt("user_id"))
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}
	response.Success(c, "获取用户信息成功", user)
}

// UsersList 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param body body dto.UserQuery true "查询参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/user/list [post]
// @Security BearerAuth
func UsersList(c *gin.Context) {
	var queryDTO dto.UserQuery
	if err := param.Validate(c, &queryDTO); err != nil {
		response.Error(c, err.Error())
		return
	}
	userService := services.NewUserService()
	items, total := userService.UsersPage(queryDTO)
	response.Success(c, "获取用户列表成功", gin.H{"items": items, "total": total})
}

// UserEdit 编辑用户
// @Summary 编辑用户信息
// @Description 创建或更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param body body dto.UserRegister true "用户信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/user/edit [post]
// @Security BearerAuth
func UserEdit(c *gin.Context) {
	var userRegisterDTO dto.UserRegister

	if err := param.Validate(c, &userRegisterDTO); err != nil {
		response.Error(c, err.Error())
		return
	}

	user := userRegisterDTO.ToModel()
	user.ID = userRegisterDTO.ID
	userService := services.NewUserService()
	if err := userService.CreateOrUpdateUser(&user); err != nil {
		response.Error(c, err.Error())
		return
	}
	roleIds := userRegisterDTO.RoleIds
	if roleIds != nil {
		rbacService := services.NewRBACService()
		err := rbacService.AssignRolesToUser(user.ID, roleIds)
		if err != nil {
			response.Error(c, "分配角色失败: "+err.Error())
			return
		}
	}

	response.Success(c, "操作成功", nil)
}