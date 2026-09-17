package models

import "strings"

// RouteDescriptions 路由描述表，key = path#method（小写 path + 大写 method），value = 路由描述。
// 由 routes 包在注册路由时写入，由 services 包在 GetPermissions 时读取。权限点不再持久化到数据库，
// 始终等于真实注册的路由，杜绝代码与数据库漂移。
var RouteDescriptions = map[string]string{}

// BuildPermKey 统一归一化权限键，注册端/校验端/绑定端共用，确保两端键完全一致。
func BuildPermKey(path, method string) string {
	return strings.ToLower(path) + "#" + strings.ToUpper(method)
}

type RBACRole struct {
	BaseFields
	Name    string `gorm:"size:50;unique" json:"name"`
	Remark  string `gorm:"size:200" json:"remark"`
	Status  int    `gorm:"default:1;comment:状态 0-禁用 1-启用" json:"status"`
	Users   []User `gorm:"many2many:rbac_user_roles;" json:"-"`
	Menus   []Menu `gorm:"many2many:rbac_role_menus" json:"-"`
	MenuIDs []int  `gorm:"-" json:"menu_ids"`
}

func (RBACRole) TableName() string {
	return "rbac_roles"
}

type RBACUserRole struct {
	UserID int `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
}

type RBACRoleMenu struct {
	RoleID int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
	MenuID int `gorm:"column:menu_id;primaryKey" json:"menu_id"`
}

func (RBACRoleMenu) TableName() string {
	return "rbac_role_menus"
}

type MenuPermission struct {
	MenuID      int    `gorm:"column:menu_id;primaryKey" json:"menu_id"`
	PermKey     string `gorm:"column:perm_key;primaryKey" json:"perm_key"`
	Description string `gorm:"column:description" json:"description"`
}

func (MenuPermission) TableName() string {
	return "rbac_menu_permissions"
}
