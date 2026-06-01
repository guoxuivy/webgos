## 1. 数据模型层

- [x] 1.1 创建部门模型 `internal/models/department.go`
- [x] 1.2 修改用户模型 `internal/models/user.go`，添加 DepartmentID 字段

## 2. DTO层

- [x] 2.1 创建部门DTO `internal/dto/department.go`，包含 AddDepartmentDTO、EditDepartmentDTO、DepartmentQuery、SetLeaderDTO

## 3. 服务层

- [x] 3.1 创建部门服务接口和实现 `internal/services/department.go`
- [x] 3.2 实现 Create、Update、Delete、GetByID、GetTree、GetUsers、SetLeader 方法

## 4. 处理器层

- [x] 4.1 创建部门处理器 `internal/handlers/department.go`
- [x] 4.2 实现 Create、Update、Delete、Get、GetTree、GetUsers、SetLeader 处理器函数

## 5. 路由层

- [x] 5.1 创建部门路由 `internal/routes/department.go`
- [x] 5.2 在 `internal/routes/routes.go` 中注册部门路由

## 6. 数据库迁移

- [x] 6.1 更新自动迁移配置

## 7. 测试验证

- [x] 7.1 编译验证
- [ ] 7.2 接口测试

## 8. 前端API层

- [x] 8.1 更新部门API文件，匹配后端接口路径 `/api/department`
- [x] 8.2 添加获取部门树、获取部门用户、设置负责人接口

## 9. 前端页面

- [x] 9.1 更新部门列表页面
- [x] 9.2 更新部门表单页面，支持父部门选择、负责人选择