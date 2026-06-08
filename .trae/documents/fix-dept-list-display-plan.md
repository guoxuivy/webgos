# 修复部门列表显示逻辑 - 计划

## 现状分析

### 已确认正常的部分
- **后端 API** `/api/department/list` 返回扁平部门列表（含 `parent_id` 字段），`getDeptList()` 工作正常（用户确认"数据有"）
- **响应拦截器**：`requestClient` 的 `code === 0` 检查正确，数据流：后端 `{code:0, data:[...]}` → 前端的 `getDeptList()` 返回 `[...]`
- **vxe-table 组件注册**：`initVxeTable()` 已注册 `VxeTable`, `VxeColumn`, `VxeGrid` 等组件
- **渲染器注册**：`CellTag`、`CellOperation` 等渲染器已通过 `vxeUI.renderer.add()` 注册

### 已修复的问题
- `vxe-column` 的 `#default` 插槽不传 slot props → 已替换为 `:cell-render="{ name: 'CellTag' }"` 和 `:cell-render="{ name: 'CellOperation', ... }"`
- 状态列和操作列不再使用 `#default` 插槽

### 待验证/可能存在的问题

1. **`#content` 插槽**：`vxe-column type="expand"` 使用 `<template #content="{ row }">` 渲染展开内容，需验证该插槽是否传 `{ row }`
2. **treeConfig 与 expandConfig 共存**：同时使用 `tree-config` 和 `expand-config` 时，展开箭头可能会冲突。tree 展开显示子部门，expand 展开显示员工。在根部门上两者都可用。
3. **treeConfig.transform 数据转换**：`parentField: 'parent_id'` 与后端 JSON 字段 `parent_id` 匹配，应该能正确构建树
4. **页面布局**：`Page auto-content-height` + 内部 `div` 布局可能导致表格高度为 0

## 实施步骤

### Step 1: 验证 `#content` 插槽是否正常
- **文件**: `d:\Goroot\webgos\frontend\apps\web-admin\src\views\system\dept\list.vue`
- **操作**: 在 `#content` 模板中添加 `{{ JSON.stringify(row) }}` 调试文本，观察展开时是否显示
- **预期**: 展开箭头点击后应显示 row 的数据 JSON

### Step 2: 处理 treeConfig 与 expandConfig 的共存
- **文件**: `d:\Goroot\webgos\frontend\apps\web-admin\src\views\system\dept\list.vue`
- **问题**: 当行同时有子部门（tree）和员工（expand）时，需要两个展开按钮同时工作
  - tree 展开通过 `tree-node` 列的箭头控制
  - expand 展开通过 `type="expand"` 列的箭头控制
- **操作**: 如果 tree 展开和 expand 展开互相干扰，使用 `tree-config` 的 `expandAll` 或 `accordion` 等配置
- **可能的方案**：
  - 方案 A：保留 `tree-config` 展示部门层级，`expand-config` 用于展开显示员工（当前方案）
  - 方案 B：去掉 `tree-config`，后端直接返回树结构数据（通过 `GetTree()`），简化显示逻辑

### Step 3: 可选 - 移除 treeConfig，改用后端树数据
- **文件**: 
  - `d:\Goroot\webgos\frontend\apps\web-admin\src\views\system\dept\list.vue`
  - `d:\Goroot\webgos\frontend\apps\web-admin\src\api\system\dept.ts`
- **操作**: 
  1. 修改 `getDeptList` 调 `/api/department/tree` 接口（或用 `/api/department/list` 已有 `GetAll` 扁平数据，前端 `transform: true`）
  2. 去掉 `treeConfig`（如果使用已构建好的树数据）
  3. 或保留 `treeConfig` 但添加 `lazy: true` 懒加载子部门

### Step 4: 添加控制台调试日志
- **文件**: `d:\Goroot\webgos\frontend\apps\web-admin\src\views\system\dept\list.vue`
- **操作**: 在 `refreshData` 和模板中添加调试信息
- **目的**: 即使 UI 正常，也保留日志以便后续排查

### Step 5: 验证页面布局正常
- **文件**: `d:\Goroot\webgos\frontend\apps\web-admin\src\views\system\dept\list.vue`
- **操作**: 确保表格容器有明确高度。`Page auto-content-height` 需要内部元素自适应
- **可能的修复**: 给 `div` 添加 `min-height` 或调整布局

## 决定点

1. **treeConfig + expandConfig 共存方案**：
   - 选项 A：继续用 `treeConfig.transform: true` 前端转换树 + `expandConfig` 展开员工（当前方案，需验证）
   - 选项 B：后端返回树结构数据，去掉 `treeConfig.transform`
   - 选项 C：去掉 `treeConfig`，部门用 `expandConfig` 展开显示子部门

2. **`expandConfig` 行为**：
   - `lazy: true` 配合 `loadMethod` 动态加载员工数据
   - 行对象上设置 `row.childCols` 和 `row.childData`

## 验证步骤
1. 刷新页面，检查 console 无报错
2. 确认部门列表以树形结构显示
3. 点击部门行的展开箭头，确认加载员工列表
4. 确认 tree 展开（子部门）正常工作
5. 确认操作按钮（新增下级、添加成员、编辑、删除）功能正常