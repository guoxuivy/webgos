## ADDED Requirements

### Requirement: Department tree displays direct members
部门树展开时，系统 SHALL 显示该部门的直接成员列表。

#### Scenario: Expand department node
- **WHEN** 用户点击部门树的展开箭头
- **THEN** 系统显示该部门的下级部门和直接成员列表

#### Scenario: Collapse department node
- **WHEN** 用户点击部门树的收起箭头
- **THEN** 系统隐藏下级部门和成员列表

### Requirement: Members are displayed in department tree
成员列表 SHALL 显示用户名、昵称、邮箱等基本信息。

#### Scenario: Display member information
- **WHEN** 部门节点展开
- **THEN** 成员列表显示用户的昵称、邮箱、电话等信息