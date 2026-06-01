## 1. 数据模型层

- [ ] 1.1 创建部门模型 `internal/models/department.go`
- [ ] 1.2 修改用户模型 `internal/models/user.go`，添加 DepartmentID 字段

## 2. DTO层

- [ ] 2.1 创建部门DTO `internal/dto/department.go`，包含 AddDepartmentDTO、EditDepartmentDTO、DepartmentQuery、SetLeaderDTO

## 3. 服务层

- [ ] 3.1 创建部门服务接口和实现 `internal/services/department.go`
- [ ] 3.2 实现 Create、Update、Delete、GetByID、GetTree、GetUsers、SetLeader 方法

## 4. 处理器层

- [ ] 4.1 创建部门处理器 `internal/handlers/department.go`
- [ ] 4.2 实现 Create、Update、Delete、Get、GetTree、GetUsers、SetLeader 处理器函数

## 5. 路由层

- [ ] 5.1 创建部门路由 `internal/routes/department.go`
- [ ] 5.2 在 `internal/routes/routes.go` 中注册部门路由

## 6. 数据库迁移

- [ ] 6.1 创建部门表迁移文件 `internal/xdb/migrate/migration_department.go`

## 7. 测试验证

- [ ] 7.1 编译验证
- [ ] 7.2 接口测试（创建、查询、更新、删除部门）