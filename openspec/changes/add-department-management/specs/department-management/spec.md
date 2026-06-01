## Overview

部门管理功能，支持树形层级结构的部门组织管理。

## Requirements

### 1. 部门模型

**字段定义:**
- `id`: 主键，自增整数
- `parent_id`: 父部门ID，支持树形结构（根部门为0或null）
- `name`: 部门名称，必填，唯一
- `leader_id`: 部门负责人ID，关联用户表
- `remark`: 部门备注，可选
- `status`: 状态（0-禁用，1-启用），默认1
- `order`: 排序字段，用于同级部门排序

### 2. 用户模型扩展

**新增字段:**
- `department_id`: 用户所属部门ID

### 3. API接口

#### 3.1 创建部门

**请求:**
```
POST /api/department
Content-Type: application/json

{
  "name": "string (必填，部门名称)",
  "parent_id": "integer (可选，父部门ID)",
  "leader_id": "integer (可选，负责人ID)",
  "remark": "string (可选，备注)",
  "status": "integer (可选，状态，默认1)",
  "order": "integer (可选，排序，默认0)"
}
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "技术部",
    "parent_id": 0,
    "leader_id": 1,
    "remark": "技术研发部门",
    "status": 1,
    "order": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3.2 更新部门

**请求:**
```
PUT /api/department/{id}
Content-Type: application/json

{
  "name": "string (可选，部门名称)",
  "parent_id": "integer (可选，父部门ID)",
  "leader_id": "integer (可选，负责人ID)",
  "remark": "string (可选，备注)",
  "status": "integer (可选，状态)",
  "order": "integer (可选，排序)"
}
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 3.3 删除部门

**请求:**
```
DELETE /api/department/{id}
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 3.4 获取部门详情

**请求:**
```
GET /api/department/{id}
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "技术部",
    "parent_id": 0,
    "leader_id": 1,
    "leader": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员"
    },
    "remark": "技术研发部门",
    "status": 1,
    "order": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 3.5 获取部门树

**请求:**
```
GET /api/department/tree
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "技术部",
      "parent_id": 0,
      "leader_id": 1,
      "children": [
        {
          "id": 2,
          "name": "前端组",
          "parent_id": 1,
          "leader_id": 2,
          "children": []
        },
        {
          "id": 3,
          "name": "后端组",
          "parent_id": 1,
          "leader_id": 3,
          "children": []
        }
      ]
    }
  ]
}
```

#### 3.6 获取部门用户

**请求:**
```
GET /api/department/{id}/users?page=1&pageSize=20
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "email": "admin@example.com"
      }
    ],
    "total": 1
  }
}
```

#### 3.7 设置部门负责人

**请求:**
```
PUT /api/department/{id}/leader
Content-Type: application/json

{
  "leader_id": "integer (必填，负责人ID)"
}
```

**响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 4. 业务规则

1. 部门名称必须唯一
2. 删除部门时需级联删除所有子部门
3. 删除部门时需将部门下的用户department_id置为null
4. 设置负责人时需验证用户存在
5. 创建/更新用户时需验证部门存在

### 5. 错误码

| 错误码 | 描述 |
|--------|------|
| 10001 | 部门名称已存在 |
| 10002 | 部门不存在 |
| 10003 | 父部门不存在 |
| 10004 | 用户不存在 |
| 10005 | 无法删除有子部门的部门（需先删除子部门） |