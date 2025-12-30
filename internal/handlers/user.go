package handlers

import (
	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/services"
	"webgos/internal/utils/param"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary 当前登录用户信息
// @Description 用户接口
// @Tags 用户
// @Accept json
// @Produce json
// @Router /api/user/info [get]
// @Security BearerAuth
func UserInfo(c *gin.Context) {
	userModel := &models.User{}
	// 预加载关联的角色信息
	userModel, err := userModel.Preload("Roles").Read(c.GetInt("user_id"))
	if err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}
	response.Success(c, "获取用户信息成功", userModel)
}

// @Summary 用户列表
// @Description 获取用户列表接口
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body dto.UserQuery true "用户列表参数"
// @Success 200 {object} response.Response "data={items: []models.User, total: int}"
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

// @Summary 用户添加、修改
// @Description 添加、修改 传ID则为修改
// @Tags 用户
// @Accept json
// @Produce json
// @Param data body dto.UserRegister true "用户注册参数"
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
