package migrate

import (
	"webgos/internal/config"
	"webgos/internal/xdb"
	"webgos/internal/models"
	"webgos/internal/xlog"
)

// AutoMigrate 执行数据库迁移
func AutoMigrate() error {
	globalConfig := config.GlobalConfig
	if globalConfig.AutoMigrate {
		xlog.Access("Starting auto migration...")
		if err := migrate(); err != nil {
			return err
		}
		xlog.Access("Auto migration completed")
	} else {
		xlog.Access("Auto migration is disabled")
	}
	return nil

}

func migrate() error {
	// 3. 官方要求：先启用 PostGIS 扩展（必须执行）
	// if err := database.GetDB().Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error; err != nil {
	// 	xlog.Error("启用 PostGIS 失败: %v", err)
	// }

	// 执行自动迁移
	return xdb.GetDB().AutoMigrate(
		&models.Product{},
		&models.InventoryRecord{},
		&models.User{},
		&models.RBACRole{},
		&models.RBACPermission{},
		&models.RBACUserRole{},
		&models.RBACRolePermission{},
		&models.Menu{},
		&models.Department{})
}
