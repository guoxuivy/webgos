--
-- PostgreSQL database dump
--

\restrict ynp2JQdPcrvqffXNgezmjdkoRkUGnj0AhKYcb5ez8TEpKdsIcAQ03Wm98ji9Q9T

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, created_at, updated_at, deleted_at, username, nickname, email, phone, password, gender, age, status, department_id) FROM stdin;
1	2025-11-27 11:08:28.385+08	2026-06-08 15:22:55.405+08	\N	admin	管理员			$2a$10$Wk1/3RbxTv.pMrwBtL/5O.44ojCwjIil7Ape063fiq.buUc17PEo.		0	1	5
2	2025-11-27 11:12:37.049+08	2026-06-04 17:59:04.036+08	\N	super	超管	guoxuwhoami@163.com	13800000000	$2a$10$oYVuKmUKXPODYNuoOqlu3e.gjzkd4AZusTaAYz3UfojVY2hn8kJXO		0	1	3
3	2025-11-27 11:13:51.698+08	2026-06-04 17:59:04.036+08	\N	guest	游客			$2a$10$0xRGARCbSKuaASaA/ynt/uG8YnXpHTdwlpA/b5V1rc0ehNjCYR8qO		0	1	3
4	2025-11-27 11:14:44.446+08	2026-08-05 15:02:14.251+08	\N	sean11	编辑+游客		13800000001	$2a$10$Xkgc/WchoqXJ.miskTjOH.WfkwLCQR9CSPJpK3QhRxBS.Em1X5fOq		0	1	0
558	2025-11-27 16:43:19.366+08	2026-08-05 15:45:41.06+08	\N	testuser_integration_1764232999308644500_8193	Test User	test@example.com	13800138000	$2a$10$HOFyUShUyrQuKwVF6o0lruMFh1cjhkFnFJ3XgWlLpflPEycDohgH.		0	1	0
559	2025-11-27 16:43:19.43+08	2026-08-05 15:46:05.954+08	\N	testuser_update_1764232999376371000_2423	Updated User	updated@example.com		$2a$10$iYOaBJ0lH406WZ9sS1hGqeCMhhbnWTycdy1kbZEIoHnNMB6hNLGqO		0	1	0
560	2025-11-27 16:43:19.486+08	2026-08-05 15:45:47.483+08	\N	testuser_delete_1764232999433267000_6098	Test User	test@example.com		$2a$10$7RkWrf/kGgnZ2G7pni.iwO45gsU72qMTsIbil/59wrEreBrtnKKSa		0	1	0
561	2025-11-27 16:43:19.592+08	2026-08-05 15:45:47.483+08	\N	testuser_query1	Query Test 1	query1@example.com		$2a$10$UPSKmAOy1X.QQHcAyd5P4eCufFc94jIBAnmvKwo8heVRXYj/UcqD6		25	1	8
562	2025-11-27 16:43:19.595+08	2026-08-05 15:45:47.483+08	\N	testuser_query2_1764232999488058500_1286	Query Test 2	query2@example.com		$2a$10$X/KvywsrRfeR39GE6j.fhuK8Q7BXoazcXAIgvrZnv8jGa4vExsxlK		30	1	8
563	2025-11-27 16:43:19.652+08	2025-11-27 16:43:19.652+08	\N	testuser_page_1764232999598424500_1458_0	Page Test 0	page0@example.com		$2a$10$scjx0L.ekAKcwfDR.2ji4e0iF.oge1Bm8ugDFig4mfdyZaDOFsgIS		0	1	0
564	2025-11-27 16:43:19.706+08	2025-11-27 16:43:19.706+08	\N	testuser_page_1764232999598424500_1458_1	Page Test 1	page1@example.com		$2a$10$30ONVw4RoixF9MV170a.wOwtXHmR13h/SMJ6cL06r/3eqwKtIRwbW		0	1	0
565	2025-11-27 16:43:19.76+08	2025-11-27 16:43:19.76+08	\N	testuser_page_1764232999598424500_1458_2	Page Test 2	page2@example.com		$2a$10$daLZDsRgXdOLEt79uvtHd.8HU808ZvrDvrosbc24H1i1HGLHm/1/i		0	1	0
566	2025-11-27 16:43:19.814+08	2025-11-27 16:43:19.814+08	\N	testuser_page_1764232999598424500_1458_3	Page Test 3	page3@example.com		$2a$10$1ImMDhcixwDX.D8nZFWvuew.5UiTSFI7DEEuysp72qptRKbez8iSO		0	1	0
567	2025-11-27 16:43:19.87+08	2025-11-27 16:43:19.87+08	\N	testuser_page_1764232999598424500_1458_4	Page Test 4	page4@example.com		$2a$10$gRUp8wSOT8rgaClmJ4GtbOwCRNE7kDGLiyQGDDVLp4.2NkuLLW9Cm		0	1	0
568	2025-11-27 16:43:19.925+08	2026-08-05 11:38:00.632+08	\N	testuser_tx1_1764232999873820900_5234	Transaction User 1	tx1@example.com		$2a$10$ar89CZuR1MUsBCAu2o3kKuSCt1AbRvkq3kM06jEKgmi/l4oVoTuJy		0	1	0
569	2025-11-27 16:43:19.979+08	2025-11-27 16:43:19.979+08	\N	testuser_tx2_1764232999873820900_1614	Transaction User 2	tx2@example.com		$2a$10$SEBcJUYYzMDLIp.WEuk2K.Q.bEthvFC8GEq2FGSND8hyfbQ6t.60O		0	1	0
570	2025-11-27 16:43:20.041+08	2025-11-27 16:43:20.041+08	\N	testuser_tx4_1764232999981688400_2376	Existing User	existing@example.com		$2a$10$sr.9/V2pElNvvXrSife.RO6ZL40IpmwS0kbp4sIes8wkw0Nx6v9Ly		0	1	0
573	2025-11-27 16:43:20.203+08	2025-11-27 16:43:20.203+08	\N	testuser_with_tx_1764233000151690300_7140	WithTx User	withtx@example.com		$2a$10$7ZeVl5qDBPdCqFhCIY76AOyd8CKRDcU1n.P1HRmRp8tC1iLqct9V6		0	1	0
828	2025-12-01 18:28:05.986+08	2026-08-05 15:41:09.184+08	\N	sean		admin@example.com		$2a$10$CWfJ2Mx/sovrlee50MfgkO8POTDRhMc5l8iMwyVZ2HDQje1To3q6q		0	1	0
833	2026-01-05 11:03:24.297+08	2026-01-05 11:03:24.297+08	\N	testuser_query1_1767582204196525600_127	Query Test 1	query1@example.com		$2a$10$gdvVt1UppyiAvcU3BwfXceebM.BDs27OVeCvW9v/2u4xJJfYjTJaq		25	1	0
834	2026-01-05 11:03:24.3+08	2026-01-05 11:03:24.3+08	\N	testuser_query2_1767582204196525600_8397	Query Test 2	query2@example.com		$2a$10$w.toSBmQxF/y57/mZzxPuOz0ASN03DmVE3oe5b3RXdg6w0icWtZZu		30	1	0
838	2026-01-05 11:12:43.926+08	2026-01-05 11:12:43.926+08	\N	testuser_query1_1767582763825391700_1641	Query Test 1	query1@example.com		$2a$10$/Og5zZR5HjECjr0WSaYHquSxVxjNQUlfXsgg1mqi0y4EzC5zVMTsu		25	1	0
839	2026-01-05 11:12:43.928+08	2026-01-05 11:12:43.928+08	\N	testuser_query2_1767582763825391700_487	Query Test 2	query2@example.com		$2a$10$TtPL1c0p.dH8vAwed7TnP./q5TG6l.L83.iPjVyC1Dh33V3AX9GRi		30	1	0
843	2026-01-05 11:13:51.555+08	2026-08-05 14:27:36.001+08	\N	testuser_query1_1767582831452586700_2019	Query Test 1	query1@example.com		$2a$10$0KnJ97zrZH1dYv0kqU7Ote5qo4h/Xkte3d9elHQc31bNx9Dw79By2		25	1	0
844	2026-01-05 11:13:51.557+08	2026-08-05 10:52:30.882+08	\N	testuser_query2_1767582831452586700_2233	Query Test 2	query2@example.com		$2a$10$yBJntLy/wsxQs5451Yx0LuQkG6VFZYxV9KL6mxjSfhDgWPPqgOQRC		30	1	0
\.


--
-- Data for Name: departments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.departments (id, created_at, updated_at, deleted_at, name, parent_id, leader_id, remark, status, sort) FROM stdin;
3	2026-06-01 16:29:23.737+08	2026-06-08 16:22:54.964+08	\N	研发中心	0	2	2	1	0
4	2026-06-01 16:38:16.27+08	2026-08-05 15:46:05.951+08	\N	前端组	3	\N		1	0
5	2026-06-01 16:38:46.454+08	2026-08-05 15:41:09.18+08	\N	后端组	3	\N	2	1	0
6	2026-06-01 16:46:33.829+08	2026-06-01 16:46:33.829+08	2026-06-02 09:59:45.994+08	flutter	4	\N		1	0
7	2026-06-02 10:00:01.63+08	2026-06-02 10:00:01.63+08	\N	运维	5	\N		1	0
8	2026-06-02 10:00:15.021+08	2026-08-05 15:08:06.143+08	\N	设计	4	\N		1	0
\.


--
-- Data for Name: inventory_records; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventory_records (id, created_at, updated_at, deleted_at, product_id, quantity, type) FROM stdin;
\.


--
-- Data for Name: menus; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.menus (id, created_at, updated_at, deleted_at, name, path, component, type, status, title, icon, affix_tab, hide_children_in_menu, hide_in_breadcrumb, hide_in_menu, hide_in_tab, keep_alive, sort, badge, badge_type, badge_variants, iframe_src, link, pid) FROM stdin;
1	2025-12-09 16:03:29.422+08	2025-12-30 11:29:37.586+08	\N	Dashboard	/dashboard		catalog	1	page.dashboard.title	carbon:workspace	f	f	f	f	f	f	0						0
3	2025-10-20 11:27:52.89+08	2025-10-20 11:49:12.935+08	\N	SystemMenu	/system/menu	/system/menu/list	menu	1	system.menu.title	carbon:menu	f	f	f	f	f	f	2						2
4	2025-10-20 11:27:52.89+08	2025-10-20 11:49:33.959+08	\N	SystemMenuCreate			button	1	common.create		f	f	f	f	f	f	3						3
5	2025-10-20 11:27:52.89+08	2025-10-20 11:49:46.323+08	\N	SystemMenuEdit			button	1	common.edit		f	f	f	f	f	f	3						3
6	2025-10-20 11:27:52.89+08	2025-10-20 11:49:13.676+08	\N	SystemMenuDelete			button	1	common.delete		f	f	f	f	f	f	3						3
7	2025-10-20 11:52:21.536+08	2026-08-05 15:45:35.853+08	\N	SystemDept	/system/dept	/system/dept/list	menu	1	system.dept.title	carbon:container-services	f	f	f	f	f	f	9						2
8	2025-10-20 11:53:23.69+08	2025-10-20 11:53:23.69+08	\N	SystemDeptCreate			button	1	common.create		f	f	f	f	f	f	7						7
9	2025-10-20 11:53:51.087+08	2025-10-20 11:53:51.087+08	\N	SystemDeptEdit			button	1	common.edit		f	f	f	f	f	f	7						7
10	2025-10-20 11:54:13.676+08	2025-10-20 11:54:13.676+08	\N	SystemDeptDelete			button	1	common.delete		f	f	f	f	f	f	7						7
12	2025-10-20 11:55:38.396+08	2025-10-20 11:55:38.396+08	\N	VbenDocument	/vben-admin/document		embedded	1	demos.vben.document	carbon:book	f	f	f	f	f	f	0				https://doc.vben.pro		11
13	2025-10-20 11:56:52.357+08	2025-10-20 11:56:52.357+08	\N	VbenAntdv			link	1	demos.vben.antdv	carbon:hexagon-vertical-solid	f	f	f	f	f	f	0		dot			https://ant.vben.pro	11
16	2025-12-09 15:41:32.104+08	2025-12-09 15:41:32.104+08	\N	权限管理	/system/permission	/system/permission/list	menu	1	权限管理	carbon:security	f	f	f	f	f	f	0						2
17	2025-12-09 15:46:59.121+08	2025-12-09 15:46:59.121+08	\N	角色管理	/system/role	/system/role/list	menu	1	角色管理	carbon:group-security	f	f	f	f	f	f	0						2
18	2025-12-09 15:48:29.178+08	2025-12-09 15:57:35.678+08	\N	用户管理	/system/user	/system/user/list	menu	1	system.user.title	carbon:user-avatar	f	f	f	f	f	f	0						2
19	2025-10-20 11:27:52.89+08	2026-08-04 18:59:23.467+08	\N	Workspace	/workspace	/dashboard/workspace/index	menu	1	page.dashboard.workspace	carbon:workspace	t	f	f	f	f	f	0						1
20	2025-10-20 11:57:36.632+08	2025-12-09 16:06:07.747+08	\N	Analytics	/analytics	/dashboard/analytics/index	menu	1	page.dashboard.analytics	carbon:text-link-analysis	f	f	f	f	f	f	0						1
26	2026-06-01 13:43:28.046+08	2026-06-01 14:44:26.314+08	\N	test1	/test13	/system/dept/list	menu	1	测试		f	f	f	f	f	f	0						11
2	2025-10-20 11:39:43.874+08	2026-08-05 11:01:35.784+08	\N	System	/system		catalog	1	system.title	carbon:settings	f	f	f	f	f	f	1	new	normal	primary			0
11	2025-10-20 11:54:53.702+08	2025-10-20 11:54:53.702+08	\N	Project	/vben-admin		catalog	1	demos.vben.title	carbon:data-center	f	f	f	f	f	f	99		dot				0
14	2025-10-20 11:57:36.632+08	2025-12-26 15:43:03.74+08	\N	About	/about	_core/about/index	menu	1	demos.vben.about	lucide:copyright	f	f	f	f	f	f	100						0
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, created_at, updated_at, deleted_at, name, stock) FROM stdin;
\.


--
-- Data for Name: rbac_menu_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rbac_menu_permissions (menu_id, rbac_permission_id) FROM stdin;
7	28
7	30
7	31
7	32
\.


--
-- Data for Name: rbac_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rbac_permissions (id, created_at, updated_at, deleted_at, name, description, path, method) FROM stdin;
18	2025-12-11 10:53:24.541+08	2026-05-20 16:37:10.328+08	\N	/api/rbac/assign_permissions#POST	分配权限给角色	/api/rbac/assign_permissions	POST
27	2026-08-05 10:59:29.806+08	2026-08-05 18:52:37.948946+08	\N	/api/department#POST	创建部门	/api/department	POST
28	2026-08-05 10:59:29.815+08	2026-08-05 18:52:37.950087+08	\N	/api/department#PUT	更新部门	/api/department	PUT
4	2025-12-11 10:53:24.501+08	2026-08-05 18:52:37.969056+08	\N	/api/products/:id#GET	获取商品详情	/api/products/:id	GET
5	2025-12-11 10:53:24.505+08	2026-08-05 18:52:37.970677+08	\N	/api/menu#POST	创建菜单	/api/menu	POST
6	2025-12-11 10:53:24.508+08	2026-08-05 18:52:37.974243+08	\N	/api/menu/:id#GET	菜单详情	/api/menu/:id	GET
7	2025-12-11 10:53:24.511+08	2026-08-05 18:52:37.975562+08	\N	/api/menu/:id#PUT	编辑菜单	/api/menu/:id	PUT
8	2025-12-11 10:53:24.513+08	2026-08-05 18:52:37.977372+08	\N	/api/menu/:id#DELETE	删除菜单	/api/menu/:id	DELETE
9	2025-12-11 10:53:24.514+08	2026-08-05 18:52:37.979387+08	\N	/api/menu/list#GET	获取菜单列表	/api/menu/list	GET
10	2025-12-11 10:53:24.516+08	2026-08-05 18:52:37.981208+08	\N	/api/menu/tree#GET	获取菜单树	/api/menu/tree	GET
11	2025-12-11 10:53:24.522+08	2026-08-05 18:52:37.982796+08	\N	/api/menu/name_exists#GET	检查菜单名称是否存在	/api/menu/name_exists	GET
12	2025-12-11 10:53:24.525+08	2026-08-05 18:52:37.983832+08	\N	/api/menu/path_exists#GET	检查菜单路径是否存在	/api/menu/path_exists	GET
13	2025-12-11 10:53:24.527+08	2026-08-05 18:52:37.985016+08	\N	/api/menu/user_menus#GET	获取当前用户目录	/api/menu/user_menus	GET
33	2026-08-05 10:59:29.963+08	2026-08-05 18:52:37.986215+08	\N	/api/menu/permissions/:id#GET	菜单权限项	/api/menu/permissions/:id	GET
34	2026-08-05 10:59:29.973+08	2026-08-05 18:52:37.989327+08	\N	/api/menu/permissions#POST	绑定菜单权限	/api/menu/permissions	POST
14	2025-12-11 10:53:24.53+08	2026-08-05 18:52:37.99128+08	\N	/api/rbac/role#POST	创建角色	/api/rbac/role	POST
15	2025-12-11 10:53:24.534+08	2026-08-05 18:52:37.994076+08	\N	/api/rbac/edit_role#POST	编辑角色	/api/rbac/edit_role	POST
16	2025-12-11 10:53:24.538+08	2026-08-05 18:52:37.995278+08	\N	/api/rbac/roles#GET	角色列表	/api/rbac/roles	GET
17	2025-12-11 10:53:24.542+08	2026-08-05 18:52:37.997256+08	\N	/api/rbac/assign_roles#POST	分配角色给用户	/api/rbac/assign_roles	POST
35	2026-08-05 10:59:30.017+08	2026-08-05 18:52:37.998298+08	\N	/api/rbac/assign_menus#POST	分配菜单给角色	/api/rbac/assign_menus	POST
19	2025-12-11 10:53:24.542+08	2026-08-05 18:52:37.998818+08	\N	/api/rbac/permission/:id#DELETE	删除权限	/api/rbac/permission/:id	DELETE
20	2025-12-11 10:53:24.544+08	2026-08-05 18:52:37.999854+08	\N	/api/rbac/permissions#GET	全部权限项	/api/rbac/permissions	GET
21	2025-12-11 10:53:24.547+08	2026-08-05 18:52:38.000368+08	\N	/api/rbac/role_permissions/:id#GET	角色权限项	/api/rbac/role_permissions/:id	GET
22	2025-12-11 10:53:24.553+08	2026-08-05 18:52:38.002046+08	\N	/api/rbac/role/:id#GET	角色详情	/api/rbac/role/:id	GET
23	2025-12-11 10:53:24.555+08	2026-08-05 18:52:38.003116+08	\N	/api/rbac/user_roles/:id#GET	用户角色	/api/rbac/user_roles/:id	GET
24	2025-12-11 10:53:24.557+08	2026-08-05 18:52:38.005067+08	\N	/api/user/info#GET	当前用户	/api/user/info	GET
25	2025-12-11 10:53:24.559+08	2026-08-05 18:52:38.007687+08	\N	/api/user/list#POST	获取用户列表	/api/user/list	POST
26	2025-12-11 10:53:24.562+08	2026-08-05 18:52:38.009787+08	\N	/api/user/edit#POST	修改用户	/api/user/edit	POST
29	2026-08-05 10:59:29.821+08	2026-08-05 18:52:37.951246+08	\N	/api/department/:id#DELETE	删除部门	/api/department/:id	DELETE
30	2026-08-05 10:59:29.826+08	2026-08-05 18:52:37.952296+08	\N	/api/department/tree#GET	部门树	/api/department/tree	GET
31	2026-08-05 10:59:29.831+08	2026-08-05 18:52:37.954169+08	\N	/api/department/:id/users#POST	批量添加用户	/api/department/:id/users	POST
32	2026-08-05 10:59:29.841+08	2026-08-05 18:52:37.960474+08	\N	/api/department/user/:userid#DELETE	移除部门用户	/api/department/user/:userid	DELETE
1	2025-12-11 10:53:24.494+08	2026-08-05 18:52:37.962865+08	\N	/api/inventory/in#POST	入库测试	/api/inventory/in	POST
2	2025-12-11 10:53:24.496+08	2026-08-05 18:52:37.966008+08	\N	/api/inventory/out#POST	出库测试	/api/inventory/out	POST
3	2025-12-11 10:53:24.499+08	2026-08-05 18:52:37.968016+08	\N	/api/products/add#POST	创建商品	/api/products/add	POST
\.


--
-- Data for Name: rbac_role_menus; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rbac_role_menus (rbac_role_id, menu_id) FROM stdin;
1	1
2	1
3	1
1	2
2	2
1	3
2	3
1	4
2	4
1	5
2	5
1	6
2	6
1	7
2	7
1	8
2	8
1	9
2	9
1	10
2	10
1	11
2	11
1	12
2	12
1	13
1	14
3	14
1	16
1	17
1	18
1	19
2	19
3	19
1	20
2	20
\.


--
-- Data for Name: rbac_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rbac_roles (id, created_at, updated_at, deleted_at, name, remark, status) FROM stdin;
1	2025-12-02 14:22:43.227+08	2025-12-11 10:51:24.579+08	\N	管理员	管理员	1
2	2025-12-05 13:41:33.689+08	2026-08-04 18:36:52.457+08	\N	编辑		1
3	2025-12-05 14:44:43.104+08	2026-08-04 18:55:06.216+08	\N	游客	游客2游客2游客2	1
\.


--
-- Data for Name: rbac_user_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rbac_user_roles (user_id, rbac_role_id) FROM stdin;
1	1
4	2
828	2
2	3
3	3
4	3
\.


--
-- Name: departments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.departments_id_seq', 8, true);


--
-- Name: inventory_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inventory_records_id_seq', 1, false);


--
-- Name: menus_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.menus_id_seq', 26, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 1, false);


--
-- Name: rbac_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rbac_permissions_id_seq', 35, true);


--
-- Name: rbac_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rbac_roles_id_seq', 3, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 844, true);


--
-- PostgreSQL database dump complete
--

\unrestrict ynp2JQdPcrvqffXNgezmjdkoRkUGnj0AhKYcb5ez8TEpKdsIcAQ03Wm98ji9Q9T

