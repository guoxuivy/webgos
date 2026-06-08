## 1. 后端 DTO 层

- [x] 1.1 添加 BatchUpdateDeptUsersDTO 数据传输对象

## 2. 后端服务层

- [x] 2.1 在 DepartmentService 接口添加 AddUsers 方法
- [x] 2.2 在 DepartmentService 接口添加 RemoveUsers 方法
- [x] 2.3 实现 AddUsers 服务方法
- [x] 2.4 实现 RemoveUsers 服务方法

## 3. 后端处理器层

- [x] 3.1 添加 AddDepartmentUsers 处理器函数
- [x] 3.2 添加 RemoveDepartmentUsers 处理器函数

## 4. 后端路由层

- [x] 4.1 注册批量添加用户路由 POST /api/department/{id}/users
- [x] 4.2 注册批量移除用户路由 DELETE /api/department/{id}/users

## 5. 前端 API 层

- [x] 5.1 添加 addDeptUsers API 函数
- [x] 5.2 添加 removeDeptUsers API 函数

## 6. 前端部门列表页面

- [ ] 6.1 修改列表页面支持展开显示成员
- [ ] 6.2 添加成员列表展示组件

## 7. 前端部门编辑页面

- [ ] 7.1 添加批量用户选择组件
- [ ] 7.2 实现批量添加用户功能
- [ ] 7.3 实现批量移除用户功能

## 8. 测试验证

- [ ] 8.1 编译验证后端代码
- [ ] 8.2 编译验证前端代码
- [ ] 8.3 测试批量添加用户接口
- [ ] 8.4 测试批量移除用户接口