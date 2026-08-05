-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: hserp
-- ------------------------------------------------------
-- Server version	8.0.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `departments`
--

DROP TABLE IF EXISTS `departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `departments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `parent_id` bigint DEFAULT '0',
  `leader_id` bigint DEFAULT NULL,
  `remark` varchar(200) DEFAULT NULL,
  `status` tinyint DEFAULT '1',
  `sort` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_departments_name` (`name`),
  KEY `idx_departments_deleted_at` (`deleted_at`),
  KEY `fk_departments_leader` (`leader_id`),
  CONSTRAINT `fk_departments_leader` FOREIGN KEY (`leader_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `departments`
--

LOCK TABLES `departments` WRITE;
/*!40000 ALTER TABLE `departments` DISABLE KEYS */;
INSERT INTO `departments` VALUES (3,'2026-06-01 16:29:23.737','2026-06-08 16:22:54.964',NULL,'研发中心',0,2,'2',1,0),(4,'2026-06-01 16:38:16.270','2026-08-05 15:46:05.951',NULL,'前端组',3,NULL,'',1,0),(5,'2026-06-01 16:38:46.454','2026-08-05 15:41:09.180',NULL,'后端组',3,NULL,'2',1,0),(6,'2026-06-01 16:46:33.829','2026-06-01 16:46:33.829','2026-06-02 09:59:45.994','flutter',4,NULL,'',1,0),(7,'2026-06-02 10:00:01.630','2026-06-02 10:00:01.630',NULL,'运维',5,NULL,'',1,0),(8,'2026-06-02 10:00:15.021','2026-08-05 15:08:06.143',NULL,'设计',4,NULL,'',1,0);
/*!40000 ALTER TABLE `departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `inventory_records`
--

DROP TABLE IF EXISTS `inventory_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `inventory_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint DEFAULT NULL,
  `quantity` bigint DEFAULT NULL,
  `type` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_inventory_records_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventory_records`
--

LOCK TABLES `inventory_records` WRITE;
/*!40000 ALTER TABLE `inventory_records` DISABLE KEYS */;
/*!40000 ALTER TABLE `inventory_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menus`
--

DROP TABLE IF EXISTS `menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menus` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `path` varchar(255) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `type` varchar(20) NOT NULL COMMENT '菜单类型',
  `status` tinyint DEFAULT '1' COMMENT '状态 0-禁用 1-启用',
  `title` varchar(100) DEFAULT NULL COMMENT '菜单标题',
  `icon` varchar(50) DEFAULT NULL COMMENT '菜单图标',
  `affix_tab` tinyint(1) DEFAULT NULL COMMENT '固定标签页',
  `hide_children_in_menu` tinyint(1) DEFAULT NULL COMMENT '隐藏子菜单',
  `hide_in_breadcrumb` tinyint(1) DEFAULT NULL COMMENT '在面包屑中隐藏',
  `hide_in_menu` tinyint(1) DEFAULT NULL COMMENT '在菜单中隐藏',
  `hide_in_tab` tinyint(1) DEFAULT NULL COMMENT '在标签页中隐藏',
  `keep_alive` tinyint(1) DEFAULT NULL COMMENT '保持活跃状态',
  `sort` int DEFAULT NULL COMMENT '排序',
  `badge` varchar(20) DEFAULT NULL COMMENT '徽标文本',
  `badge_type` varchar(20) DEFAULT NULL COMMENT '徽标类型',
  `badge_variants` varchar(20) DEFAULT NULL COMMENT '徽标样式',
  `iframe_src` varchar(255) DEFAULT NULL COMMENT 'iframe地址',
  `link` varchar(255) DEFAULT NULL COMMENT '外链地址',
  `pid` bigint DEFAULT NULL COMMENT '父级菜单ID',
  PRIMARY KEY (`id`),
  KEY `idx_sys_menus_deleted_at` (`deleted_at`),
  KEY `idx_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menus`
--

LOCK TABLES `menus` WRITE;
/*!40000 ALTER TABLE `menus` DISABLE KEYS */;
INSERT INTO `menus` VALUES (1,'2025-12-09 16:03:29.422','2025-12-30 11:29:37.586',NULL,'Dashboard','/dashboard','','catalog',1,'page.dashboard.title','carbon:workspace',0,0,0,0,0,0,0,'','','','','',0),(2,'2025-10-20 11:39:43.874','2026-08-05 11:01:35.784',NULL,'System','/system','','catalog',1,'system.title','carbon:settings',0,0,0,0,0,0,0,'new','normal','primary','','',0),(3,'2025-10-20 11:27:52.890','2025-10-20 11:49:12.935',NULL,'SystemMenu','/system/menu','/system/menu/list','menu',1,'system.menu.title','carbon:menu',0,0,0,0,0,0,0,'','','','','',2),(4,'2025-10-20 11:27:52.890','2025-10-20 11:49:33.959',NULL,'SystemMenuCreate','','','button',1,'common.create','',0,0,0,0,0,0,0,'','','','','',3),(5,'2025-10-20 11:27:52.890','2025-10-20 11:49:46.323',NULL,'SystemMenuEdit','','','button',1,'common.edit','',0,0,0,0,0,0,0,'','','','','',3),(6,'2025-10-20 11:50:11.061','2025-10-20 11:50:11.061',NULL,'SystemMenuDelete','','','button',1,'common.delete','',0,0,0,0,0,0,0,'','','','','',3),(7,'2025-10-20 11:52:21.536','2026-08-05 15:45:35.853',NULL,'SystemDept','/system/dept','/system/dept/list','menu',1,'system.dept.title','carbon:container-services',0,0,0,0,0,0,9,'','','','','',2),(8,'2025-10-20 11:53:23.690','2025-10-20 11:53:23.690',NULL,'SystemDeptCreate','','','button',1,'common.create','',0,0,0,0,0,0,0,'','','','','',7),(9,'2025-10-20 11:53:51.087','2025-10-20 11:53:51.087',NULL,'SystemDeptEdit','','','button',1,'common.edit','',0,0,0,0,0,0,0,'','','','','',7),(10,'2025-10-20 11:54:13.676','2025-10-20 11:54:13.676',NULL,'SystemDeptDelete','','','button',1,'common.delete','',0,0,0,0,0,0,0,'','','','','',7),(11,'2025-10-20 11:54:53.702','2025-10-20 11:54:53.702',NULL,'Project','/vben-admin','','catalog',1,'demos.vben.title','carbon:data-center',0,0,0,0,0,0,0,'','dot','','','',0),(12,'2025-10-20 11:55:38.396','2025-10-20 11:55:38.396',NULL,'VbenDocument','/vben-admin/document','','embedded',1,'demos.vben.document','carbon:book',0,0,0,0,0,0,0,'','','','https://doc.vben.pro','',11),(13,'2025-10-20 11:56:52.357','2025-10-20 11:56:52.357',NULL,'VbenAntdv','','','link',1,'demos.vben.antdv','carbon:hexagon-vertical-solid',0,0,0,0,0,0,0,'','dot','','','https://ant.vben.pro',11),(14,'2025-10-20 11:57:36.632','2025-12-26 15:43:03.740',NULL,'About','/about','_core/about/index','menu',1,'demos.vben.about','lucide:copyright',0,0,0,0,0,0,0,'','','','','',0),(16,'2025-12-09 15:41:32.104','2025-12-09 15:41:32.104',NULL,'权限管理','/system/permission','/system/permission/list','menu',1,'权限管理','carbon:security',0,0,0,0,0,0,0,'','','','','',2),(17,'2025-12-09 15:46:59.121','2025-12-09 15:46:59.121',NULL,'角色管理','/system/role','/system/role/list','menu',1,'角色管理','carbon:group-security',0,0,0,0,0,0,0,'','','','','',2),(18,'2025-12-09 15:48:29.178','2025-12-09 15:57:35.678',NULL,'用户管理','/system/user','/system/user/list','menu',1,'system.user.title','carbon:user-avatar',0,0,0,0,0,0,0,'','','','','',2),(19,'2025-10-20 11:27:52.890','2026-08-04 18:59:23.467',NULL,'Workspace','/workspace','/dashboard/workspace/index','menu',1,'page.dashboard.workspace','carbon:workspace',1,0,0,0,0,0,0,'','','','','',1),(20,'2025-12-09 16:06:07.747','2025-12-09 16:06:32.646',NULL,'Analytics','/analytics','/dashboard/analytics/index','menu',1,'page.dashboard.analytics','carbon:text-link-analysis',0,0,0,0,0,0,0,'','','','','',1),(26,'2026-06-01 13:43:28.046','2026-06-01 14:44:26.314',NULL,'test1','/test13','/system/dept/list','menu',1,'测试','',0,0,0,0,0,0,0,'','','','','',11);
/*!40000 ALTER TABLE `menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `stock` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_menu_permissions`
--

DROP TABLE IF EXISTS `rbac_menu_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rbac_menu_permissions` (
  `menu_id` bigint NOT NULL,
  `rbac_permission_id` bigint NOT NULL,
  PRIMARY KEY (`menu_id`,`rbac_permission_id`),
  KEY `fk_rbac_menu_permissions_rbac_permission` (`rbac_permission_id`),
  CONSTRAINT `fk_rbac_menu_permissions_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`),
  CONSTRAINT `fk_rbac_menu_permissions_rbac_permission` FOREIGN KEY (`rbac_permission_id`) REFERENCES `rbac_permissions` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_menu_permissions`
--

LOCK TABLES `rbac_menu_permissions` WRITE;
/*!40000 ALTER TABLE `rbac_menu_permissions` DISABLE KEYS */;
INSERT INTO `rbac_menu_permissions` VALUES (7,28),(7,30),(7,31),(7,32);
/*!40000 ALTER TABLE `rbac_menu_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_permissions`
--

DROP TABLE IF EXISTS `rbac_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rbac_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `method` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_rbac_permissions_name` (`name`),
  KEY `idx_rbac_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_permissions`
--

LOCK TABLES `rbac_permissions` WRITE;
/*!40000 ALTER TABLE `rbac_permissions` DISABLE KEYS */;
INSERT INTO `rbac_permissions` VALUES (1,'2025-12-11 10:53:24.494','2026-08-05 15:07:01.458',NULL,'/api/inventory/in#POST','入库测试','/api/inventory/in','POST'),(2,'2025-12-11 10:53:24.496','2026-08-05 15:07:01.460',NULL,'/api/inventory/out#POST','出库测试','/api/inventory/out','POST'),(3,'2025-12-11 10:53:24.499','2026-08-05 15:07:01.463',NULL,'/api/products/add#POST','创建商品','/api/products/add','POST'),(4,'2025-12-11 10:53:24.501','2026-08-05 15:07:01.466',NULL,'/api/products/:id#GET','获取商品详情','/api/products/:id','GET'),(5,'2025-12-11 10:53:24.505','2026-08-05 15:07:01.469',NULL,'/api/menu#POST','创建菜单','/api/menu','POST'),(6,'2025-12-11 10:53:24.508','2026-08-05 15:07:01.471',NULL,'/api/menu/:id#GET','菜单详情','/api/menu/:id','GET'),(7,'2025-12-11 10:53:24.511','2026-08-05 15:07:01.473',NULL,'/api/menu/:id#PUT','编辑菜单','/api/menu/:id','PUT'),(8,'2025-12-11 10:53:24.513','2026-08-05 15:07:01.475',NULL,'/api/menu/:id#DELETE','删除菜单','/api/menu/:id','DELETE'),(9,'2025-12-11 10:53:24.514','2026-08-05 15:07:01.478',NULL,'/api/menu/list#GET','获取菜单列表','/api/menu/list','GET'),(10,'2025-12-11 10:53:24.516','2026-08-05 15:07:01.483',NULL,'/api/menu/tree#GET','获取菜单树','/api/menu/tree','GET'),(11,'2025-12-11 10:53:24.522','2026-08-05 15:07:01.486',NULL,'/api/menu/name_exists#GET','检查菜单名称是否存在','/api/menu/name_exists','GET'),(12,'2025-12-11 10:53:24.525','2026-08-05 15:07:01.489',NULL,'/api/menu/path_exists#GET','检查菜单路径是否存在','/api/menu/path_exists','GET'),(13,'2025-12-11 10:53:24.527','2026-08-05 15:07:01.491',NULL,'/api/menu/user_menus#GET','获取当前用户目录','/api/menu/user_menus','GET'),(14,'2025-12-11 10:53:24.529','2026-08-05 15:07:01.500',NULL,'/api/rbac/role#POST','创建角色','/api/rbac/role','POST'),(15,'2025-12-11 10:53:24.530','2026-08-05 15:07:01.502',NULL,'/api/rbac/edit_role#POST','编辑角色','/api/rbac/edit_role','POST'),(16,'2025-12-11 10:53:24.534','2026-08-05 15:07:01.504',NULL,'/api/rbac/roles#GET','角色列表','/api/rbac/roles','GET'),(17,'2025-12-11 10:53:24.538','2026-08-05 15:07:01.507',NULL,'/api/rbac/assign_roles#POST','分配角色给用户','/api/rbac/assign_roles','POST'),(18,'2025-12-11 10:53:24.541','2026-05-20 16:37:10.328',NULL,'/api/rbac/assign_permissions#POST','分配权限给角色','/api/rbac/assign_permissions','POST'),(19,'2025-12-11 10:53:24.542','2026-08-05 15:07:01.514',NULL,'/api/rbac/permission/:id#DELETE','删除权限','/api/rbac/permission/:id','DELETE'),(20,'2025-12-11 10:53:24.544','2026-08-05 15:07:01.518',NULL,'/api/rbac/permissions#GET','全部权限项','/api/rbac/permissions','GET'),(21,'2025-12-11 10:53:24.547','2026-08-05 15:07:01.520',NULL,'/api/rbac/role_permissions/:id#GET','角色权限项','/api/rbac/role_permissions/:id','GET'),(22,'2025-12-11 10:53:24.553','2026-08-05 15:07:01.523',NULL,'/api/rbac/role/:id#GET','角色详情','/api/rbac/role/:id','GET'),(23,'2025-12-11 10:53:24.555','2026-08-05 15:07:01.525',NULL,'/api/rbac/user_roles/:id#GET','用户角色','/api/rbac/user_roles/:id','GET'),(24,'2025-12-11 10:53:24.557','2026-08-05 15:07:01.527',NULL,'/api/user/info#GET','当前用户','/api/user/info','GET'),(25,'2025-12-11 10:53:24.559','2026-08-05 15:07:01.530',NULL,'/api/user/list#POST','获取用户列表','/api/user/list','POST'),(26,'2025-12-11 10:53:24.562','2026-08-05 15:07:01.535',NULL,'/api/user/edit#POST','修改用户','/api/user/edit','POST'),(27,'2026-08-05 10:59:29.806','2026-08-05 15:07:01.442',NULL,'/api/department#POST','创建部门','/api/department','POST'),(28,'2026-08-05 10:59:29.815','2026-08-05 15:07:01.444',NULL,'/api/department#PUT','更新部门','/api/department','PUT'),(29,'2026-08-05 10:59:29.821','2026-08-05 15:07:01.447',NULL,'/api/department/:id#DELETE','删除部门','/api/department/:id','DELETE'),(30,'2026-08-05 10:59:29.826','2026-08-05 15:07:01.450',NULL,'/api/department/tree#GET','部门树','/api/department/tree','GET'),(31,'2026-08-05 10:59:29.831','2026-08-05 15:07:01.453',NULL,'/api/department/:id/users#POST','批量添加用户','/api/department/:id/users','POST'),(32,'2026-08-05 10:59:29.841','2026-08-05 15:07:01.456',NULL,'/api/department/user/:userid#DELETE','移除部门用户','/api/department/user/:userid','DELETE'),(33,'2026-08-05 10:59:29.963','2026-08-05 15:07:01.494',NULL,'/api/menu/permissions/:id#GET','菜单权限项','/api/menu/permissions/:id','GET'),(34,'2026-08-05 10:59:29.973','2026-08-05 15:07:01.497',NULL,'/api/menu/permissions#POST','绑定菜单权限','/api/menu/permissions','POST'),(35,'2026-08-05 10:59:30.017','2026-08-05 15:07:01.510',NULL,'/api/rbac/assign_menus#POST','分配菜单给角色','/api/rbac/assign_menus','POST');
/*!40000 ALTER TABLE `rbac_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_role_menus`
--

DROP TABLE IF EXISTS `rbac_role_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rbac_role_menus` (
  `rbac_role_id` bigint NOT NULL,
  `menu_id` bigint NOT NULL,
  PRIMARY KEY (`rbac_role_id`,`menu_id`),
  KEY `fk_rbac_role_menus_menu` (`menu_id`),
  CONSTRAINT `fk_rbac_role_menus_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`),
  CONSTRAINT `fk_rbac_role_menus_rbac_role` FOREIGN KEY (`rbac_role_id`) REFERENCES `rbac_roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_role_menus`
--

LOCK TABLES `rbac_role_menus` WRITE;
/*!40000 ALTER TABLE `rbac_role_menus` DISABLE KEYS */;
INSERT INTO `rbac_role_menus` VALUES (1,1),(2,1),(3,1),(1,2),(2,2),(1,3),(2,3),(1,4),(2,4),(1,5),(2,5),(1,6),(2,6),(1,7),(2,7),(1,8),(2,8),(1,9),(2,9),(1,10),(2,10),(1,11),(2,11),(1,12),(2,12),(1,13),(1,14),(3,14),(1,16),(1,17),(1,18),(1,19),(2,19),(3,19),(1,20),(2,20);
/*!40000 ALTER TABLE `rbac_role_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_roles`
--

DROP TABLE IF EXISTS `rbac_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rbac_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(50) DEFAULT NULL,
  `remark` varchar(200) DEFAULT NULL,
  `status` tinyint DEFAULT '1' COMMENT '状态 0-禁用 1-启用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_rbac_roles_name` (`name`),
  KEY `idx_rbac_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_roles`
--

LOCK TABLES `rbac_roles` WRITE;
/*!40000 ALTER TABLE `rbac_roles` DISABLE KEYS */;
INSERT INTO `rbac_roles` VALUES (1,'2025-12-02 14:22:43.227','2025-12-11 10:51:24.579',NULL,'管理员','管理员',1),(2,'2025-12-05 13:41:33.689','2026-08-04 18:36:52.457',NULL,'编辑','',1),(3,'2025-12-05 14:44:43.104','2026-08-04 18:55:06.216',NULL,'游客','游客2游客2游客2',1);
/*!40000 ALTER TABLE `rbac_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rbac_user_roles`
--

DROP TABLE IF EXISTS `rbac_user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `rbac_user_roles` (
  `user_id` bigint NOT NULL,
  `rbac_role_id` bigint NOT NULL,
  PRIMARY KEY (`user_id`,`rbac_role_id`),
  KEY `fk_rbac_user_roles_rbac_role` (`rbac_role_id`),
  CONSTRAINT `fk_rbac_user_roles_rbac_role` FOREIGN KEY (`rbac_role_id`) REFERENCES `rbac_roles` (`id`),
  CONSTRAINT `fk_rbac_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rbac_user_roles`
--

LOCK TABLES `rbac_user_roles` WRITE;
/*!40000 ALTER TABLE `rbac_user_roles` DISABLE KEYS */;
INSERT INTO `rbac_user_roles` VALUES (1,1),(4,2),(828,2),(2,3),(3,3),(4,3);
/*!40000 ALTER TABLE `rbac_user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(191) DEFAULT NULL,
  `nickname` longtext,
  `email` longtext,
  `phone` longtext,
  `password` longtext,
  `gender` longtext,
  `age` bigint DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `department_id` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=872 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2025-11-27 11:08:28.385','2026-06-08 15:22:55.405',NULL,'admin','管理员','','','$2a$10$Wk1/3RbxTv.pMrwBtL/5O.44ojCwjIil7Ape063fiq.buUc17PEo.','',0,1,5),(2,'2025-11-27 11:12:37.049','2026-06-04 17:59:04.036',NULL,'super','超管','guoxuwhoami@163.com','13800000000','$2a$10$oYVuKmUKXPODYNuoOqlu3e.gjzkd4AZusTaAYz3UfojVY2hn8kJXO','',0,1,3),(3,'2025-11-27 11:13:51.698','2026-06-04 17:59:04.036',NULL,'guest','游客','','','$2a$10$0xRGARCbSKuaASaA/ynt/uG8YnXpHTdwlpA/b5V1rc0ehNjCYR8qO','',0,1,3),(4,'2025-11-27 11:14:44.446','2026-08-05 15:02:14.251',NULL,'sean11','编辑+游客','','13800000001','$2a$10$Xkgc/WchoqXJ.miskTjOH.WfkwLCQR9CSPJpK3QhRxBS.Em1X5fOq','',0,1,0),(558,'2025-11-27 16:43:19.366','2026-08-05 15:45:41.060',NULL,'testuser_integration_1764232999308644500_8193','Test User','test@example.com','13800138000','$2a$10$HOFyUShUyrQuKwVF6o0lruMFh1cjhkFnFJ3XgWlLpflPEycDohgH.','',0,1,0),(559,'2025-11-27 16:43:19.430','2026-08-05 15:46:05.954',NULL,'testuser_update_1764232999376371000_2423','Updated User','updated@example.com','','$2a$10$iYOaBJ0lH406WZ9sS1hGqeCMhhbnWTycdy1kbZEIoHnNMB6hNLGqO','',0,1,0),(560,'2025-11-27 16:43:19.486','2026-08-05 15:45:47.483',NULL,'testuser_delete_1764232999433267000_6098','Test User','test@example.com','','$2a$10$7RkWrf/kGgnZ2G7pni.iwO45gsU72qMTsIbil/59wrEreBrtnKKSa','',0,1,8),(561,'2025-11-27 16:43:19.592','2026-08-05 15:45:47.483',NULL,'testuser_query1','Query Test 1','query1@example.com','','$2a$10$UPSKmAOy1X.QQHcAyd5P4eCufFc94jIBAnmvKwo8heVRXYj/UcqD6','',25,1,8),(562,'2025-11-27 16:43:19.595','2026-08-05 15:45:47.483',NULL,'testuser_query2_1764232999488058500_1286','Query Test 2','query2@example.com','','$2a$10$X/KvywsrRfeR39GE6j.fhuK8Q7BXoazcXAIgvrZnv8jGa4vExsxlK','',30,1,8),(563,'2025-11-27 16:43:19.652','2025-11-27 16:43:19.652',NULL,'testuser_page_1764232999598424500_1458_0','Page Test 0','page0@example.com','','$2a$10$scjx0L.ekAKcwfDR.2ji4e0iF.oge1Bm8ugDFig4mfdyZaDOFsgIS','',0,1,0),(564,'2025-11-27 16:43:19.706','2025-11-27 16:43:19.706',NULL,'testuser_page_1764232999598424500_1458_1','Page Test 1','page1@example.com','','$2a$10$30ONVw4RoixF9MV170a.wOwtXHmR13h/SMJ6cL06r/3eqwKtIRwbW','',0,1,0),(565,'2025-11-27 16:43:19.760','2025-11-27 16:43:19.760',NULL,'testuser_page_1764232999598424500_1458_2','Page Test 2','page2@example.com','','$2a$10$daLZDsRgXdOLEt79uvtHd.8HU808ZvrDvrosbc24H1i1HGLHm/1/i','',0,1,0),(566,'2025-11-27 16:43:19.814','2025-11-27 16:43:19.814',NULL,'testuser_page_1764232999598424500_1458_3','Page Test 3','page3@example.com','','$2a$10$1ImMDhcixwDX.D8nZFWvuew.5UiTSFI7DEEuysp72qptRKbez8iSO','',0,1,0),(567,'2025-11-27 16:43:19.870','2025-11-27 16:43:19.870',NULL,'testuser_page_1764232999598424500_1458_4','Page Test 4','page4@example.com','','$2a$10$gRUp8wSOT8rgaClmJ4GtbOwCRNE7kDGLiyQGDDVLp4.2NkuLLW9Cm','',0,1,0),(568,'2025-11-27 16:43:19.925','2026-08-05 11:38:00.632',NULL,'testuser_tx1_1764232999873820900_5234','Transaction User 1','tx1@example.com','','$2a$10$ar89CZuR1MUsBCAu2o3kKuSCt1AbRvkq3kM06jEKgmi/l4oVoTuJy','',0,1,0),(569,'2025-11-27 16:43:19.979','2025-11-27 16:43:19.979',NULL,'testuser_tx2_1764232999873820900_1614','Transaction User 2','tx2@example.com','','$2a$10$SEBcJUYYzMDLIp.WEuk2K.Q.bEthvFC8GEq2FGSND8hyfbQ6t.60O','',0,1,0),(570,'2025-11-27 16:43:20.041','2025-11-27 16:43:20.041',NULL,'testuser_tx4_1764232999981688400_2376','Existing User','existing@example.com','','$2a$10$sr.9/V2pElNvvXrSife.RO6ZL40IpmwS0kbp4sIes8wkw0Nx6v9Ly','',0,1,0),(573,'2025-11-27 16:43:20.203','2025-11-27 16:43:20.203',NULL,'testuser_with_tx_1764233000151690300_7140','WithTx User','withtx@example.com','','$2a$10$7ZeVl5qDBPdCqFhCIY76AOyd8CKRDcU1n.P1HRmRp8tC1iLqct9V6','',0,1,0),(828,'2025-12-01 18:28:05.986','2026-08-05 15:41:09.184',NULL,'sean','','admin@example.com','','$2a$10$CWfJ2Mx/sovrlee50MfgkO8POTDRhMc5l8iMwyVZ2HDQje1To3q6q','',0,1,0),(833,'2026-01-05 11:03:24.297','2026-01-05 11:03:24.297',NULL,'testuser_query1_1767582204196525600_127','Query Test 1','query1@example.com','','$2a$10$gdvVt1UppyiAvcU3BwfXceebM.BDs27OVeCvW9v/2u4xJJfYjTJaq','',25,1,0),(834,'2026-01-05 11:03:24.300','2026-01-05 11:03:24.300',NULL,'testuser_query2_1767582204196525600_8397','Query Test 2','query2@example.com','','$2a$10$w.toSBmQxF/y57/mZzxPuOz0ASN03DmVE3oe5b3RXdg6w0icWtZZu','',30,1,0),(838,'2026-01-05 11:12:43.926','2026-01-05 11:12:43.926',NULL,'testuser_query1_1767582763825391700_1641','Query Test 1','query1@example.com','','$2a$10$/Og5zZR5HjECjr0WSaYHquSxVxjNQUlfXsgg1mqi0y4EzC5zVMTsu','',25,1,0),(839,'2026-01-05 11:12:43.928','2026-01-05 11:12:43.928',NULL,'testuser_query2_1767582763825391700_487','Query Test 2','query2@example.com','','$2a$10$TtPL1c0p.dH8vAwed7TnP./q5TG6l.L83.iPjVyC1Dh33V3AX9GRi','',30,1,0),(843,'2026-01-05 11:13:51.555','2026-08-05 14:27:36.001',NULL,'testuser_query1_1767582831452586700_2019','Query Test 1','query1@example.com','','$2a$10$0KnJ97zrZH1dYv0kqU7Ote5qo4h/Xkte3d9elHQc31bNx9Dw79By2','',25,1,0),(844,'2026-01-05 11:13:51.557','2026-08-05 10:52:30.882',NULL,'testuser_query2_1767582831452586700_2233','Query Test 2','query2@example.com','','$2a$10$yBJntLy/wsxQs5451Yx0LuQkG6VFZYxV9KL6mxjSfhDgWPPqgOQRC','',30,1,0);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-08-05 19:03:49
