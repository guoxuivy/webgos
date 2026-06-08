## Why

当前部门管理功能缺少用户部门关系绑定的完整支持，需要实现部门与用户的关联管理，包括在部门列表中展示成员信息，以及在部门编辑中批量添加/移除员工的能力。

## What Changes

- 部门列表页面支持展开查看下级部门和直接成员列表
- 部门编辑页面支持批量添加员工（从用户表选择）
- 部门编辑页面支持批量移除员工
- 新增后端API支持批量用户操作

## Capabilities

### New Capabilities
- `dept-user-list-display`: 在部门树中展示部门直接成员列表
- `dept-user-batch-add`: 批量添加用户到部门
- `dept-user-batch-remove`: 批量从部门移除用户

### Modified Capabilities
- `department-management`: 扩展部门管理功能，增加用户绑定能力

## Impact

- 后端：新增 `AddUsers`、`RemoveUsers` 服务方法和对应的API接口
- 前端：修改部门列表页面展示成员，修改部门编辑页面支持批量用户操作
- 数据库：用户表的 `department_id` 字段用于关联部门