package migrate

import (
	"webgos/internal/config"
	"webgos/internal/database"
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
	return database.DB.AutoMigrate(
		&models.Product{},
		&models.InventoryRecord{},
		&models.User{},
		&models.RBACRole{},
		&models.RBACPermission{},
		&models.RBACUserRole{},
		&models.RBACRolePermission{},
		&models.Menu{})
}
