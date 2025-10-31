package models

import (
	"strconv"
	"strings"
)

// RBACRole 角色模型
type RBACRole struct {
	BaseModel[RBACRole]
	Name        string           `gorm:"size:50;unique" json:"name"`
	Remark      string           `gorm:"size:200" json:"remark"`
	Status      int              `gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-启用" json:"status"` // 状态 0-禁用 1-启用
	Permissions []RBACPermission `gorm:"many2many:rbac_role_permissions;" json:"permissions"`
	Users       []User           `gorm:"many2many:rbac_user_roles;" json:"-"`
	Menus       []int            `gorm:"-" json:"menus"`              // 菜单ID数组，不直接存储在数据库中
	MenuIDs     string           `gorm:"column:menu_ids;size:500" json:"-"` // 以逗号分隔的形式存储菜单ID
}

// TableName 指定表名
func (RBACRole) TableName() string {
	return "rbac_roles"
}

// 实现RBACRole的Setter方法，将菜单ID数组转换为逗号分隔的字符串
func (r *RBACRole) SetMenuIDs(menus []int) {
	if len(menus) == 0 {
		r.MenuIDs = ""
		return
	}
	
	strMenus := make([]string, len(menus))
	for i, menu := range menus {
		strMenus[i] = strconv.Itoa(menu)
	}
	r.MenuIDs = strings.Join(strMenus, ",")
}

// 实现RBACRole的Getter方法，将逗号分隔的字符串转换为菜单ID数组
func (r *RBACRole) GetMenuIDs() []int {
	if r.MenuIDs == "" {
		return []int{}
	}
	
	strMenus := strings.Split(r.MenuIDs, ",")
	menus := make([]int, 0, len(strMenus))
	
	for _, strMenu := range strMenus {
		if strMenu != "" {
			if menuID, err := strconv.Atoi(strMenu); err == nil {
				menus = append(menus, menuID)
			}
		}
	}
	
	return menus
}

// RBACPermission 权限模型(对应权限点)
type RBACPermission struct {
	BaseModel[RBACPermission]
	Name        string     `gorm:"size:100;unique" json:"name"` // 权限名称，如：user:list
	Description string     `gorm:"size:200" json:"description"` // 权限描述
	Path        string     `gorm:"size:255" json:"path"`        // 路由路径
	Method      string     `gorm:"size:20" json:"method"`       // 请求方法
	Roles       []RBACRole `gorm:"many2many:rbac_role_permissions;" json:"-"`
}

// 关联表
type RBACUserRole struct {
	UserID int `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
}

type RBACRolePermission struct {
	RoleID       int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
	PermissionID int `gorm:"column:rbac_permission_id;primaryKey" json:"rbac_permission_id"`
}