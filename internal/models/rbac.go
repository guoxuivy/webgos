package models

type RBACRole struct {
	BaseFields
	Name        string    `gorm:"size:50;unique" json:"name"`
	Remark      string    `gorm:"size:200" json:"remark"`
	Status      int       `gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-启用" json:"status"`
	Users       []User    `gorm:"many2many:rbac_user_roles;" json:"-"`
	Menus       []Menu    `gorm:"many2many:rbac_role_menus" json:"-"`
	MenuIDs     []int     `gorm:"-" json:"menu_ids"`
}

func (RBACRole) TableName() string {
	return "rbac_roles"
}

type RBACPermission struct {
	BaseFields
	Name        string    `gorm:"size:100;unique" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	Path        string    `gorm:"size:255" json:"path"`
	Method      string    `gorm:"size:20" json:"method"`
	Menus       []Menu    `gorm:"many2many:rbac_menu_permissions" json:"menus"`
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

type RBACMenuPermission struct {
	MenuID       int `gorm:"column:menu_id;primaryKey" json:"menu_id"`
	PermissionID int `gorm:"column:rbac_permission_id;primaryKey" json:"rbac_permission_id"`
}

func (RBACMenuPermission) TableName() string {
	return "rbac_menu_permissions"
}