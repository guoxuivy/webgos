## Why

当前系统缺少部门管理功能，无法组织和管理企业的组织架构。为了完善权限管理体系，需要实现部门的树形层级管理，支持用户与部门的关联，以及部门负责人的设置。

## What Changes

- 新增部门数据模型，支持树形层级结构（parent_id关联）
- 新增部门管理API接口（增删改查）
- 修改用户模型，添加部门ID字段关联
- 部门负责人通过user_id字段实现
- 部门列表以树形结构展示，无需分页

## Capabilities

### New Capabilities
- `department-management`: 部门管理功能，包括部门的增删改查、树形结构展示、负责人设置

### Modified Capabilities
- `user-management`: 用户模型新增部门关联字段

## Impact

- 修改文件：`internal/models/user.go` - 添加DepartmentID字段
- 新增文件：`internal/models/department.go` - 部门模型
- 新增文件：`internal/dto/department.go` - 部门DTO
- 新增文件：`internal/services/department.go` - 部门服务层
- 新增文件：`internal/handlers/department.go` - 部门处理器
- 新增文件：`internal/routes/department.go` - 部门路由