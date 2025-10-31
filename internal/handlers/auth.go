package handlers

import (
	"hserp/internal/config"
	"hserp/internal/dto"
	"hserp/internal/services"
	"hserp/internal/utils"
	"hserp/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary 用户登录
// @Description 用户登录接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param data body dto.Login true "登录参数"
// @Success 200 {object} response.Response "data={accessToken: string}"
// @Failure 400 {object} response.Response
// @Router /auth/login [post]
// Login 用户登录
func Login(c *gin.Context) {
	var userLoginDTO dto.Login

	if err := utils.Validate(c, &userLoginDTO); err != nil {
		response.Error(c, err.Error())
		return
	}

	service := services.NewAuthService()
	token, err := service.Login(userLoginDTO.Username, userLoginDTO.Password)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "登录成功", gin.H{"accessToken": token})
}

// @Summary 用户登出
// @Description 用户登出接口
// @Tags 登录
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
// Logout 用户登出
func Logout(c *gin.Context) {
	service := services.NewAuthService()
	service.Logout(c.GetString("tokenString"))
	response.Success(c, "登出成功", nil)
}

// @Summary 重置密码 测试专用
// @Description 重置密码接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param data body dto.Login true "登录参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/reset-password [post]
// ResetPassword 重置密码
func ResetPassword(c *gin.Context) {

	if config.GlobalConfig.Server.Mode != "debug" {
		response.Error(c, "生产模式不允许重置密码")
		return
	}
	var userLoginDTO dto.Login

	if err := utils.Validate(c, &userLoginDTO); err != nil {
		response.Error(c, err.Error())
	}

	userService := services.NewUserService()
	err := userService.ResetPassword(userLoginDTO.Username, userLoginDTO.Password)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "重置密码成功", nil)
}

// @Summary 用户注册
// @Description 用户注册接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param data body dto.UserRegister true "用户注册参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var userRegisterDTO dto.UserRegister

	if err := utils.Validate(c, &userRegisterDTO); err != nil {
		response.Error(c, err.Error())
		return
	}

	user := userRegisterDTO.ToModel()
	userService := services.NewUserService()
	if err := userService.CreateOrUpdateUser(&user); err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "用户注册成功", nil)
}
