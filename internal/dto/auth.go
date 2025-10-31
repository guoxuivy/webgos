package dto

type Login struct {
	Username string `json:"username" form:"username" validate:"required,min=2,max=20" label:"用户名"` // 用户名
	Password string `json:"password" form:"password" validate:"required,min=6" label:"密码"`         // 密码
}
