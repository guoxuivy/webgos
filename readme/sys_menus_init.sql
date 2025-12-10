-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: webgos
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
-- Table structure for table `sys_menus`
--

DROP TABLE IF EXISTS `sys_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menus` (
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
  `order` bigint DEFAULT NULL COMMENT '排序',
  `badge` varchar(20) DEFAULT NULL COMMENT '徽标文本',
  `badge_type` varchar(20) DEFAULT NULL COMMENT '徽标类型',
  `badge_variants` varchar(20) DEFAULT NULL COMMENT '徽标样式',
  `iframe_src` varchar(255) DEFAULT NULL COMMENT 'iframe地址',
  `link` varchar(255) DEFAULT NULL COMMENT '外链地址',
  `pid` bigint DEFAULT NULL COMMENT '父级菜单ID',
  `auth_code` varchar(255) DEFAULT NULL COMMENT '权限标识',
  PRIMARY KEY (`id`),
  KEY `idx_sys_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menus`
--

/*!40000 ALTER TABLE `sys_menus` DISABLE KEYS */;
INSERT INTO `sys_menus` VALUES (1,'2025-12-09 16:03:29.422','2025-12-09 18:04:53.253',NULL,'Dashboard','/dashboard','','catalog',1,'page.dashboard.title','carbon:workspace',0,0,0,0,0,0,0,'','','','','',0,''),(2,'2025-10-20 11:39:43.874','2025-10-20 11:39:43.874',NULL,'System','/system','','catalog',1,'system.title','carbon:settings',0,0,0,0,0,0,0,'new','normal','primary','','',0,NULL),(3,'2025-10-20 11:27:52.890','2025-10-20 11:49:12.935',NULL,'SystemMenu','/system/menu','/system/menu/list','menu',1,'system.menu.title','carbon:menu',0,0,0,0,0,0,0,'','','','','',2,'System:Menu:List'),(4,'2025-10-20 11:27:52.890','2025-10-20 11:49:33.959',NULL,'SystemMenuCreate','','','button',1,'common.create','',0,0,0,0,0,0,0,'','','','','',3,'System:Menu:Create'),(5,'2025-10-20 11:27:52.890','2025-10-20 11:49:46.323',NULL,'SystemMenuEdit','','','button',1,'common.edit','',0,0,0,0,0,0,0,'','','','','',3,'System:Menu:Edit'),(6,'2025-10-20 11:50:11.061','2025-10-20 11:50:11.061',NULL,'SystemMenuDelete','','','button',1,'common.delete','',0,0,0,0,0,0,0,'','','','','',3,'System:Menu:Delete'),(7,'2025-10-20 11:52:21.536','2025-10-20 11:52:21.536',NULL,'SystemDept','/system/dept','/system/dept/list','menu',1,'system.dept.title','carbon:container-services',0,0,0,0,0,0,0,'','','','','',2,'System:Dept:List'),(8,'2025-10-20 11:53:23.690','2025-10-20 11:53:23.690',NULL,'SystemDeptCreate','','','button',1,'common.create','',0,0,0,0,0,0,0,'','','','','',7,'System:Dept:Create'),(9,'2025-10-20 11:53:51.087','2025-10-20 11:53:51.087',NULL,'SystemDeptEdit','','','button',1,'common.edit','',0,0,0,0,0,0,0,'','','','','',7,'System:Dept:Edit'),(10,'2025-10-20 11:54:13.676','2025-10-20 11:54:13.676',NULL,'SystemDeptDelete','','','button',1,'common.delete','',0,0,0,0,0,0,0,'','','','','',7,'System:Dept:Delete'),(11,'2025-10-20 11:54:53.702','2025-10-20 11:54:53.702',NULL,'Project','/vben-admin','','catalog',1,'demos.vben.title','carbon:data-center',0,0,0,0,0,0,0,'','dot','','','',0,''),(12,'2025-10-20 11:55:38.396','2025-10-20 11:55:38.396',NULL,'VbenDocument','/vben-admin/document','','embedded',1,'demos.vben.document','carbon:book',0,0,0,0,0,0,0,'','','','https://doc.vben.pro','',11,''),(13,'2025-10-20 11:56:52.357','2025-10-20 11:56:52.357',NULL,'VbenAntdv','','','link',1,'demos.vben.antdv','carbon:hexagon-vertical-solid',0,0,0,0,0,0,0,'','dot','','','https://ant.vben.pro',11,''),(14,'2025-10-20 11:57:36.632','2025-10-20 11:57:36.632',NULL,'About','/about','_core/about/index','menu',1,'demos.vben.about','lucide:copyright',0,0,0,0,0,0,0,'','','','','',0,''),(16,'2025-12-09 15:41:32.104','2025-12-09 15:41:32.104',NULL,'权限管理','/system/permission','/system/permission/list','menu',1,'权限管理','carbon:security',0,0,0,0,0,0,0,'','','','','',2,''),(17,'2025-12-09 15:46:59.121','2025-12-09 15:46:59.121',NULL,'角色管理','/system/role','/system/role/list','menu',1,'角色管理','carbon:group-security',0,0,0,0,0,0,0,'','','','','',2,''),(18,'2025-12-09 15:48:29.178','2025-12-09 15:57:35.678',NULL,'用户管理','/system/user','/system/user/list','menu',1,'system.user.title','carbon:user-avatar',0,0,0,0,0,0,0,'','','','','',2,''),(19,'2025-10-20 11:27:52.890','2025-12-09 18:18:08.456',NULL,'Workspace','/workspace','/dashboard/workspace/index','catalog',1,'page.dashboard.workspace','carbon:workspace',1,0,0,0,0,0,0,'','','','','',1,''),(20,'2025-12-09 16:06:07.747','2025-12-09 16:06:32.646',NULL,'Analytics','/analytics','/dashboard/analytics/index','menu',1,'page.dashboard.analytics','carbon:text-link-analysis',0,0,0,0,0,0,0,'','','','','',1,'');
/*!40000 ALTER TABLE `sys_menus` ENABLE KEYS */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-10 17:52:10
