-- ============================================================================
-- RBAC 改造迁移脚本（SQL 版）
-- 背景：权限点从"角色直连(rbac_role_permissions)"改为"经菜单间接获得"。
--       角色菜单从 rbac_roles.menu_ids 字符串改为 rbac_role_menus 中间表。
--       权限点-菜单为多对多，使用 rbac_menu_permissions 中间表。
--
-- 适用：MySQL 8.0+ / PostgreSQL（已分别标注方言差异）
-- 说明：本脚本负责【建表 + 角色菜单数据搬运 + 清理旧结构】。
--       权限点-菜单的自动归属（按路由路径前缀匹配）由应用启动时
--       SyncPermissions 自动完成，无需在 SQL 中处理。
-- ============================================================================

-- ---------------------------------------------------------------------------
-- Step 1. 创建中间表
-- ---------------------------------------------------------------------------

-- 角色-菜单关联表
DROP TABLE IF EXISTS rbac_role_menus;
CREATE TABLE rbac_role_menus (
    rbac_role_id BIGINT NOT NULL COMMENT '角色ID',
    menu_id      BIGINT NOT NULL COMMENT '菜单ID',
    PRIMARY KEY (rbac_role_id, menu_id),
    KEY idx_rbac_role_menus_menu (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 菜单-权限关联表（多对多，同一权限可绑多个菜单）
DROP TABLE IF EXISTS rbac_menu_permissions;
CREATE TABLE rbac_menu_permissions (
    menu_id           BIGINT NOT NULL COMMENT '菜单ID',
    rbac_permission_id BIGINT NOT NULL COMMENT '权限点ID',
    PRIMARY KEY (menu_id, rbac_permission_id),
    KEY idx_rbac_menu_permissions_perm (rbac_permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ---------------------------------------------------------------------------
-- Step 2. 迁移历史数据：rbac_roles.menu_ids 字符串 -> rbac_role_menus
--         仅插入菜单真实存在的关联，已存在则跳过。
-- ---------------------------------------------------------------------------

-- MySQL 语法（支持 INSERT ... ON DUPLICATE KEY UPDATE 做幂等）
INSERT INTO rbac_role_menus (rbac_role_id, menu_id)
SELECT r.id, CAST(TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(r.menu_ids, ',', n.n), ',', -1)) AS UNSIGNED) AS mid
FROM rbac_roles r
JOIN (
    SELECT 1 AS n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5
    UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
    UNION ALL SELECT 11 UNION ALL SELECT 12 UNION ALL SELECT 13 UNION ALL SELECT 14 UNION ALL SELECT 15
    UNION ALL SELECT 16 UNION ALL SELECT 17 UNION ALL SELECT 18 UNION ALL SELECT 19 UNION ALL SELECT 20
) n
  ON CHAR_LENGTH(r.menu_ids) - CHAR_LENGTH(REPLACE(r.menu_ids, ',', '')) >= n.n - 1
WHERE r.menu_ids IS NOT NULL AND r.menu_ids <> ''
  AND CAST(TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(r.menu_ids, ',', n.n), ',', -1)) AS UNSIGNED) IN (SELECT id FROM sys_menus)
ON DUPLICATE KEY UPDATE rbac_role_id = r.id;

-- PostgreSQL 语法（将上面 MySQL 段落替换为以下存储过程式插入）：
-- INSERT INTO rbac_role_menus (rbac_role_id, menu_id)
-- SELECT r.id, m::bigint AS mid
-- FROM rbac_roles r,
--      unnest(string_to_array(r.menu_ids, ',')) AS m
-- WHERE r.menu_ids IS NOT NULL AND r.menu_ids <> ''
--   AND m::bigint IN (SELECT id FROM sys_menus)
-- ON CONFLICT (rbac_role_id, menu_id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Step 3. 清理旧结构
--         注意：权限点-菜单的自动归属由应用 SyncPermissions 完成，
--               之后即可安全删除以下旧对象。
-- ---------------------------------------------------------------------------

DROP TABLE IF EXISTS rbac_role_permissions;

-- MySQL：删除 rbac_roles.menu_ids 列
-- ALTER TABLE rbac_roles DROP COLUMN IF EXISTS menu_ids;

-- PostgreSQL：删除 rbac_roles.menu_ids 列
-- ALTER TABLE rbac_roles DROP COLUMN IF EXISTS menu_ids;

-- MySQL：删除 sys_menus.auth_code 列（权限标识功能已移除）
-- ALTER TABLE sys_menus DROP COLUMN IF EXISTS auth_code;

-- PostgreSQL：删除 sys_menus.auth_code 列
-- ALTER TABLE sys_menus DROP COLUMN IF EXISTS auth_code;
