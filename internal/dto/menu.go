package dto

// AddMenuDTO 添加菜单DTO
type AddMenuDTO struct {
	Name      string   `json:"name" validate:"required" label:"菜单名称"`
	Path      string   `json:"path" validate:"omitempty" label:"路由路径"`
	AuthCode  string   `json:"authCode" validate:"omitempty" label:"权限标识"`
	Component string   `json:"component" validate:"omitempty" label:"组件路径"`
	Type      string   `json:"type" validate:"required,oneof=catalog menu button embedded link" label:"菜单类型"`
	Status    int      `json:"status" validate:"required,oneof=0 1" label:"状态"` // 0-禁用 1-启用
	Meta      MenuMeta `json:"meta" validate:"required" label:"菜单元数据"`
	Pid       int      `json:"pid" validate:"omitempty" label:"父级菜单ID"`
}

// EditMenuDTO 编辑菜单DTO
type EditMenuDTO struct {
	ID        int      `json:"id" validate:"required" label:"菜单ID"`
	Name      string   `json:"name" validate:"required" label:"菜单名称"`
	Path      string   `json:"path" validate:"omitempty" label:"路由路径"`
	AuthCode  string   `json:"authCode" validate:"omitempty" label:"权限标识"`
	Component string   `json:"component" validate:"omitempty" label:"组件路径"`
	Type      string   `json:"type" validate:"required,oneof=catalog menu button embedded link" label:"菜单类型"`
	Status    int      `json:"status" validate:"required,oneof=0 1" label:"状态"` // 0-禁用 1-启用
	Meta      MenuMeta `json:"meta" validate:"required" label:"菜单元数据"`
	Pid       int      `json:"pid" validate:"omitempty" label:"父级菜单ID"`
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title              string `json:"title" validate:"required" label:"菜单标题"`
	Icon               string `json:"icon" validate:"omitempty" label:"菜单图标"`
	AffixTab           bool   `json:"affixTab" validate:"omitempty" label:"固定标签页"`
	HideChildrenInMenu bool   `json:"hideChildrenInMenu" validate:"omitempty" label:"隐藏子菜单"`
	HideInBreadcrumb   bool   `json:"hideInBreadcrumb" validate:"omitempty" label:"在面包屑中隐藏"`
	HideInMenu         bool   `json:"hideInMenu" validate:"omitempty" label:"在菜单中隐藏"`
	HideInTab          bool   `json:"hideInTab" validate:"omitempty" label:"在标签页中隐藏"`
	KeepAlive          bool   `json:"keepAlive" validate:"omitempty" label:"保持活跃状态"`
	Order              int    `json:"order" validate:"omitempty" label:"排序"`
	Badge              string `json:"badge" validate:"omitempty" label:"徽标文本"`
	BadgeType          string `json:"badgeType" validate:"omitempty" label:"徽标类型"`
	BadgeVariants      string `json:"badgeVariants" validate:"omitempty" label:"徽标样式"`
	IframeSrc          string `json:"iframeSrc" validate:"omitempty" label:"iframe地址"`
	Link               string `json:"link" validate:"omitempty" label:"外链地址"`
}

// GetMenuDTO 获取菜单DTO
type GetMenuDTO struct {
	ID int `uri:"id" validate:"required" label:"菜单ID"`
}

// DeleteMenuDTO 删除菜单DTO
type DeleteMenuDTO struct {
	ID int `uri:"id" validate:"required" label:"菜单ID"`
}
