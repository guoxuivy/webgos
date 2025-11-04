package migrate

import (
	"webgos/internal/database"
	"webgos/internal/models"
)

// AutoMigrate 执行数据库迁移
func AutoMigrate() error {
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
