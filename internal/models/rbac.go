package models

import (
	"strconv"
	"strings"
)

type RBACRole struct {
	BaseFields
	Name          string           `gorm:"size:50;unique" json:"name"`
	Remark        string           `gorm:"size:200" json:"remark"`
	MenuIDs       string           `gorm:"column:menu_ids;size:500" json:"-"`
	Status        int              `gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-启用" json:"status"`
	Permissions   []RBACPermission `gorm:"many2many:rbac_role_permissions;" json:"permissions"`
	Users         []User           `gorm:"many2many:rbac_user_roles;" json:"-"`
	PermissionIDs []int            `gorm:"-" json:"permission_ids"`
	Menus         []int            `gorm:"-" json:"menus"`
}

func (RBACRole) TableName() string {
	return "rbac_roles"
}

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

type RBACPermission struct {
	BaseFields
	Name        string     `gorm:"size:100;unique" json:"name"`
	Description string     `gorm:"size:200" json:"description"`
	Path        string     `gorm:"size:255" json:"path"`
	Method      string     `gorm:"size:20" json:"method"`
	Roles       []RBACRole `gorm:"many2many:rbac_role_permissions;" json:"-"`
}

type RBACUserRole struct {
	UserID int `gorm:"column:user_id;primaryKey" json:"user_id"`
	RoleID int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
}

type RBACRolePermission struct {
	RoleID       int `gorm:"column:rbac_role_id;primaryKey" json:"rbac_role_id"`
	PermissionID int `gorm:"column:rbac_permission_id;primaryKey" json:"rbac_permission_id"`
}