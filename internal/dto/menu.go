package dto

import "webgos/internal/models"

// 菜单DTO
type MenuDTO struct {
	Name      string          `json:"name" validate:"required" label:"菜单名称"`
	Path      string          `json:"path" validate:"omitempty" label:"路由路径"`
	AuthCode  string          `json:"authCode" validate:"omitempty" label:"权限标识"`
	Component string          `json:"component" validate:"omitempty" label:"组件路径"`
	Type      string          `json:"type" validate:"required,oneof=catalog menu button embedded link" label:"菜单类型"`
	Status    int             `json:"status" validate:"required,oneof=0 1" label:"状态"` // 0-禁用 1-启用
	Meta      models.MenuMeta `json:"meta" validate:"required" label:"菜单元数据"`
	Pid       int             `json:"pid" validate:"omitempty" label:"父级菜单ID"`
}
