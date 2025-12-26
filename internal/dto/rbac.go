package dto

// AddRoleDTO 添加角色DTO
type AddRoleDTO struct {
	Name   string `json:"name" validate:"required,min=1,max=50" label:"角色名称"`
	Remark string `json:"remark" validate:"omitempty,max=200" label:"角色备注"`
	Status int    `json:"status" validate:"oneof=0 1" label:"状态"` // 0-禁用 1-启用
	Menus  []int  `json:"menus" validate:"omitempty" label:"菜单ID列表"`
}

type EditRoleDTO struct {
	ID     int     `uri:"id" validate:"required" label:"角色ID"`
	Name   *string `json:"name" validate:"omitempty,min=1,max=50" label:"角色名称"`
	Remark *string `json:"remark" validate:"omitempty,max=200" label:"角色备注"`
	Status *int    `json:"status" validate:"omitempty,oneof=0 1" label:"状态"` // 0-禁用 1-启用
	Menus  []int   `json:"menus" validate:"omitempty" label:"菜单ID列表"`
}

// AssignRolesDTO 分配角色给用户DTO
type AssignRolesDTO struct {
	UserID  int   `json:"user_id" validate:"required" label:"用户ID"`
	RoleIDs []int `json:"role_ids" validate:"required,min=1" label:"角色ID列表"`
}

// AssignPermissionsDTO 分配权限给角色DTO
type AssignPermissionsDTO struct {
	RoleID        int   `json:"role_id" validate:"required" label:"角色ID"`
	PermissionIDs []int `json:"permission_ids" validate:"required" label:"权限ID列表"`
}

// AssignMenusDTO 分配菜单给角色DTO
// type AssignMenusDTO struct {
// 	RoleID  int   `json:"role_id" validate:"required" label:"角色ID"`
// 	MenuIDs []int `json:"menu_ids" validate:"required,min=1" label:"菜单ID列表"`
// }

// GetRoleDTO 获取角色DTO
type GetRoleDTO struct {
	ID int `uri:"id" validate:"required" label:"角色ID"`
}

// GetUserRolesDTO 获取用户角色DTO
type GetUserRolesDTO struct {
	UserID int `uri:"id" validate:"required" label:"用户ID"`
}

// GetRolePermissionsDTO 获取角色权限DTO
type GetRolePermissionsDTO struct {
	RoleID int `uri:"id" validate:"required" label:"角色ID"`
}

// DeletePermissionDTO 删除权限DTO
type DeletePermissionDTO struct {
	ID int `uri:"id" validate:"required" label:"权限ID"`
}
