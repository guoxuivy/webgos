## ADDED Requirements

### Requirement: Batch remove users from department
系统 SHALL 支持批量将用户从指定部门移除。

#### Scenario: Successful batch remove
- **WHEN** 用户在部门编辑页面选择多个成员并点击移除
- **THEN** 系统将选中用户从该部门移除，用户的部门ID置为0

#### Scenario: Remove users from non-existent department
- **WHEN** 尝试从不存在的部门移除用户
- **THEN** 系统返回错误提示"部门不存在"

#### Scenario: Remove users not in department
- **WHEN** 尝试移除不属于该部门的用户
- **THEN** 系统忽略这些用户，只移除属于该部门的用户

#### Scenario: Remove empty user list
- **WHEN** 用户未选择任何用户点击移除
- **THEN** 系统返回错误提示"请选择用户"

### Requirement: Leader protection
系统 SHALL 禁止移除部门负责人。

#### Scenario: Remove department leader
- **WHEN** 用户尝试移除部门负责人
- **THEN** 系统返回错误提示"不能移除部门负责人"