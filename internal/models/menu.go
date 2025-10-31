package models

// Menu 菜单模型
type Menu struct {
	BaseModel[Menu]

	Name      string   `gorm:"size:50;not null;comment:菜单名称" json:"name"`
	Path      string   `gorm:"size:255;comment:路由路径" json:"path"`
	Component string   `gorm:"size:255;comment:组件路径" json:"component"`
	AuthCode  string   `gorm:"size:255;comment:权限标识" json:"authCode"`
	Type      string   `gorm:"size:20;not null;comment:菜单类型" json:"type"` // catalog, menu, button, embedded, link
	Status    int      `gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-启用" json:"status"`
	Redirect  string   `gorm:"-" json:"redirect"` // 重定向路径，不存储在数据库中
	Meta      MenuMeta `gorm:"embedded" json:"meta"`
	Pid       int      `gorm:"comment:父级菜单ID" json:"pid"`
	Children  []Menu   `gorm:"-" json:"children,omitempty"` // 子菜单，不存储在数据库中
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title                 string   `gorm:"size:100;comment:菜单标题" json:"title"`
	Icon                  string   `gorm:"size:50;comment:菜单图标" json:"icon"`
	AffixTab              bool     `gorm:"comment:固定标签页" json:"affixTab"`
	Authority             []string `gorm:"-" json:"authority,omitempty"` // 权限标识，不存储在数据库中
	MenuVisibleWithForbidden bool  `gorm:"-" json:"menuVisibleWithForbidden,omitempty"` // 是否在无权限时显示菜单，不存储在数据库中
	HideChildrenInMenu    bool     `gorm:"comment:隐藏子菜单" json:"hideChildrenInMenu"`
	HideInBreadcrumb      bool     `gorm:"comment:在面包屑中隐藏" json:"hideInBreadcrumb"`
	HideInMenu            bool     `gorm:"comment:在菜单中隐藏" json:"hideInMenu"`
	HideInTab             bool     `gorm:"comment:在标签页中隐藏" json:"hideInTab"`
	KeepAlive             bool     `gorm:"comment:保持活跃状态" json:"keepAlive"`
	Order                 int      `gorm:"comment:排序" json:"order"`
	Badge                 string   `gorm:"size:20;comment:徽标文本" json:"badge"`
	BadgeType             string   `gorm:"size:20;comment:徽标类型" json:"badgeType"`
	BadgeVariants         string   `gorm:"size:20;comment:徽标样式" json:"badgeVariants"`
	IframeSrc             string   `gorm:"size:255;comment:iframe地址" json:"iframeSrc"`
	Link                  string   `gorm:"size:255;comment:外链地址" json:"link"`
}

// TableName 设置菜单表名
func (Menu) TableName() string {
	return "sys_menus"
}
