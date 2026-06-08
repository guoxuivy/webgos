## ADDED Requirements

### Requirement: Batch add users to department
系统 SHALL 支持批量将用户添加到指定部门。

#### Scenario: Successful batch add
- **WHEN** 用户在部门编辑页面选择多个用户并点击添加
- **THEN** 系统将选中用户添加到该部门

#### Scenario: Add users to non-existent department
- **WHEN** 尝试向不存在的部门添加用户
- **THEN** 系统返回错误提示"部门不存在"

#### Scenario: Add empty user list
- **WHEN** 用户未选择任何用户点击添加
- **THEN** 系统返回错误提示"请选择用户"

### Requirement: User list selection
系统 SHALL 从用户表获取可选用户列表供选择。

#### Scenario: Select users from list
- **WHEN** 用户打开部门编辑页面的用户选择器
- **THEN** 系统显示所有未禁用的用户列表