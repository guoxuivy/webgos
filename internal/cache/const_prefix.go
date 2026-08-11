package cache

// 业务缓存键前缀
const (
	UserPagePrefix = "user:page"
	UserMenuPrefix = "user:menusByUserID"
	// PermissionPrefix 用户权限缓存键前缀，格式：permissions:<userID>
	PermissionPrefix = "permissions"
)
