package dto

import (
	"webgos/internal/models"
)

type UserRegister struct {
	ID       int    `json:"id" validate:"omitempty,gte=0" label:"id"`                              // id
	Username string `json:"username" form:"username" validate:"required,min=2,max=20" label:"用户名"` // 用户名
	Phone    string `json:"phone" form:"phone" validate:"omitempty,phone" label:"手机号码"`            // 手机号码
	Password string `json:"password" validate:"omitempty,min=6" label:"密码"`                        // 明文，仅用于验证
	Nickname string `json:"nickname" validate:"omitempty,min=2,max=20" label:"昵称"`                 // 昵称
	Email    string `json:"email" validate:"omitempty,email" label:"邮箱"`                           //	邮箱
	Age      int    `json:"age" validate:"omitempty,gte=0,lte=150" label:"年龄"`                     // 年龄
	Gender   string `json:"gender" validate:"omitempty,oneof=male female" label:"性别"`              // 性别
	Status   int    `json:"status" validate:"oneof=0 1" label:"状态"`                                // 状态 0禁用 1启用
}

func (dto *UserRegister) ToModel() models.User {
	userModel := models.User{
		Username: dto.Username,
		Email:    dto.Email,
		Phone:    dto.Phone,
		Password: dto.Password,
		Age:      dto.Age,
		Nickname: dto.Nickname,
		Gender:   dto.Gender,
		Status:   dto.Status,
	}
	return userModel
}

type UserQuery struct {
	Page     int    `form:"page" validate:"omitempty,min=1" label:"页码"`
	PageSize int    `form:"pageSize" validate:"omitempty,min=1,max=100" label:"每页数量"`
	Username string `form:"username" validate:"omitempty,min=2,max=20" label:"用户名"`
}
