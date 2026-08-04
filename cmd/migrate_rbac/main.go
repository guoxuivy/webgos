// migrate_rbac 一次性迁移脚本：将 RBAC 从"角色直连权限 + MenuIDs 字符串"改造为
// "角色绑菜单、权限经菜单间接获得" 的多对多结构。
//
// 用法：
//   go run ./cmd/migrate_rbac -c ./config/config.yaml
//
// 执行步骤：
//   1. 加载配置 + 初始化 DB + AutoMigrate（自动建 rbac_role_menus / rbac_menu_permissions）
//   2. 注册路由并同步权限点（SyncPermissions 按路径前缀写入 rbac_menu_permissions）
//   3. 迁移历史：rbac_roles.menu_ids 字符串 -> rbac_role_menus 中间表
//   4. 清理：删除 rbac_role_permissions 表、rbac_roles.menu_ids 列
package main

import (
	"flag"
	"strconv"
	"strings"

	"webgos/internal/config"
	"webgos/internal/models"
	"webgos/internal/routes"
	"webgos/internal/xdb"
	"webgos/internal/xlog"

	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("c", "./config/config.yaml", "Specify the config file path")
	flag.Parse()

	// 1. 初始化配置与数据库
	if _, err := config.LoadConfig(*configPath); err != nil {
		panic("failed to load config: " + err.Error())
	}
	if err := xlog.InitLogger(); err != nil {
		panic("failed to init logger: " + err.Error())
	}
	if err := xdb.InitDB(); err != nil {
		panic("failed to init db: " + err.Error())
	}
	defer xdb.CloseDB()

	db := xdb.GetDB()

	// 确保中间表存在（AutoMigrate 会依据 many2many 标签自动建表）
	if err := db.AutoMigrate(&models.RBACRole{}, &models.RBACPermission{}, &models.Menu{}); err != nil {
		panic("failed to auto migrate: " + err.Error())
	}
	xlog.Access("[migrate] 中间表已就绪 (rbac_role_menus / rbac_menu_permissions)")

	// 2. 注册路由并同步权限点 -> 自动写入 rbac_menu_permissions（按路径前缀匹配，支持一权限多菜单）
	routes.New(config.GlobalConfig)
	if err := routes.SyncPermissions(db); err != nil {
		panic("failed to sync permissions: " + err.Error())
	}
	xlog.Access("[migrate] 权限点已同步并归属菜单 (rbac_menu_permissions)")

	// 3. 迁移历史：rbac_roles.menu_ids 字符串 -> rbac_role_menus
	if err := migrateRoleMenuIDs(db); err != nil {
		panic("failed to migrate role menu ids: " + err.Error())
	}

	// 4. 清理旧结构
	if err := cleanup(db); err != nil {
		panic("failed to cleanup: " + err.Error())
	}

	xlog.Access("[migrate] RBAC 迁移完成 ")
}

// migrateRoleMenuIDs 将 rbac_roles.menu_ids（逗号分隔）迁移到 rbac_role_menus 中间表。
// 已存在的关联会跳过，保证可重入。
func migrateRoleMenuIDs(db *gorm.DB) error {
	type legacyRole struct {
		ID      int
		MenuIDs string
	}
	var roles []legacyRole
	// 读取仍保留 menu_ids 列的角色（列被删除后此查询会失败，可安全跳过迁移）
	if err := db.Table("rbac_roles").Select("id, menu_ids").Scan(&roles).Error; err != nil {
		xlog.Access("[migrate] 跳过角色菜单迁移（rbac_roles.menu_ids 列可能已不存在）：%v", err)
		return nil
	}

	migrated, skipped := 0, 0
	for _, r := range roles {
		r.MenuIDs = strings.TrimSpace(r.MenuIDs)
		if r.MenuIDs == "" {
			continue
		}
		ids := strings.Split(r.MenuIDs, ",")
		for _, s := range ids {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			menuID, err := strconv.Atoi(s)
			if err != nil {
				continue
			}
			// 菜单是否存在
			var cnt int64
			if err := db.Table("sys_menus").Where("id = ?", menuID).Count(&cnt).Error; err != nil || cnt == 0 {
				continue
			}
			// 关联是否已存在
			var exist int64
			db.Table("rbac_role_menus").
				Where("rbac_role_id = ? AND menu_id = ?", r.ID, menuID).
				Count(&exist)
			if exist > 0 {
				skipped++
				continue
			}
			if err := db.Table("rbac_role_menus").
				Create(map[string]interface{}{"rbac_role_id": r.ID, "menu_id": menuID}).Error; err != nil {
				return err
			}
			migrated++
		}
	}
	xlog.Access("[migrate] 角色菜单迁移完成：新增 %d 条关联，跳过 %d 条已存在", migrated, skipped)
	return nil
}

// cleanup 删除旧结构：rbac_role_permissions 表、rbac_roles.menu_ids 列。
func cleanup(db *gorm.DB) error {
	// 删除 rbac_role_permissions（权限-角色直连关系已不再使用）
	if db.Migrator().HasTable("rbac_role_permissions") {
		if err := db.Migrator().DropTable("rbac_role_permissions"); err != nil {
			return err
		}
		xlog.Access("[migrate] 已删除旧表 rbac_role_permissions")
	}

	// 删除 rbac_roles.menu_ids 列（角色菜单改由中间表维护）
	if db.Migrator().HasColumn(&models.RBACRole{}, "menu_ids") {
		if err := db.Migrator().DropColumn(&models.RBACRole{}, "menu_ids"); err != nil {
			xlog.Access("[migrate] 删除列 menu_ids 失败（可忽略，手动 DROP 亦可）：%v", err)
		} else {
			xlog.Access("[migrate] 已删除 rbac_roles.menu_ids 列")
		}
	}

	// 删除 sys_menus.auth_code 列（权限标识功能已移除）
	if db.Migrator().HasColumn(&models.Menu{}, "auth_code") {
		if err := db.Migrator().DropColumn(&models.Menu{}, "auth_code"); err != nil {
			xlog.Access("[migrate] 删除列 auth_code 失败（可忽略，手动 DROP 亦可）：%v", err)
		} else {
			xlog.Access("[migrate] 已删除 sys_menus.auth_code 列")
		}
	}
	return nil
}
