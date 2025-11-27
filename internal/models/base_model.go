package models

import (
	"database/sql"
	"time"
	"webgos/internal/database"

	"gorm.io/gorm"
)

// IActiveRecord 定义了基础模型的接口，提供通用的数据库操作方法
// T 是泛型参数，代表具体的模型类型
type IActiveRecord[T any] interface {
	// WithTx 将查询器与事务对象绑定，用于在事务中执行数据库操作
	// 参数 tx: 事务对象
	// 返回值: 绑定了事务对象的模型实例
	WithTx(tx *gorm.DB) *BaseModel[T]

	// Transaction 在事务中执行数据库操作
	// 使用gorm的Transaction方法包装事务操作
	// 参数 fc: 事务执行函数，包含在事务中执行的业务逻辑
	// 参数 opts: 事务选项
	// 返回值: 执行事务过程中可能发生的错误
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)

	// Create 创建一条新记录
	// 参数 item: 要创建的记录对象
	// 返回值: 创建过程中可能发生的错误
	Create(item *T) error

	// Read 根据ID读取一条记录
	// 参数 id: 要读取记录的ID
	// 返回值: 读取到的记录对象和可能发生的错误
	Read(id int) (*T, error)

	// Update 更新一条记录
	// 参数 item: 包含更新数据的记录对象
	// 返回值: 更新过程中可能发生的错误
	Update(item *T) error

	// Delete 根据ID删除一条记录（软删除）
	// 参数 id: 要删除记录的ID
	// 返回值: 删除过程中可能发生的错误
	Delete(id int) error

	// More 查询多条记录
	// 返回值: 查询到的记录列表和可能发生的错误
	More() ([]T, error)

	// One 查询单条记录
	// 返回值: 查询到的记录对象和可能发生的错误
	One() (*T, error)

	// Count 统计记录总数
	// 返回值: 记录总数和可能发生的错误
	Count() (int, error)

	// Page 分页查询记录
	// 使用gorm的Offset和Limit方法实现分页查询
	// page: 页数，从1开始
	// pageSize: 每页记录数
	// 返回值: 查询到的记录列表和总记录数
	Page(page, pageSize int) ([]T, int)

	// 链式查询方法start

	// Where 添加WHERE条件
	// 参数 query: 查询条件
	// 参数 args: 查询条件参数
	// 返回值: 支持链式调用的接口
	Where(query any, args ...any) IActiveRecord[T]

	// Select 指定要查询的字段
	// 参数 query: 要查询的字段
	// 参数 args: 查询字段参数
	// 返回值: 支持链式调用的接口
	Select(query any, args ...any) IActiveRecord[T]

	// Preload 预加载关联数据
	// 参数 query: 关联查询语句
	// 参数 args: 关联查询参数
	// 返回值: 支持链式调用的接口
	Preload(query string, args ...any) IActiveRecord[T]

	// Order 添加排序条件
	// 参数 value: 排序条件
	// 返回值: 支持链式调用的接口
	Order(value any) IActiveRecord[T]

	// Not 添加NOT条件
	// 参数 query: NOT条件
	// 参数 args: NOT条件参数
	// 返回值: 支持链式调用的接口
	Not(query any, args ...any) IActiveRecord[T]

	// Or 添加OR条件
	// 参数 query: OR条件
	// 参数 args: OR条件参数
	// 返回值: 支持链式调用的接口
	Or(query any, args ...any) IActiveRecord[T]

	// Limit 限制返回记录数
	// 参数 limit: 限制的记录数
	// 返回值: 支持链式调用的接口
	Limit(limit int) IActiveRecord[T]

	// Group 添加分组条件
	// 参数 query: 分组条件
	// 返回值: 支持链式调用的接口
	Group(query string) IActiveRecord[T]

	// Having 添加分组过滤条件
	// 参数 query: 分组过滤条件
	// 参数 args: 分组过滤条件参数
	// 返回值: 支持链式调用的接口
	Having(query any, args ...any) IActiveRecord[T]

	// Joins 添加JOIN连接查询
	// 参数 query: JOIN查询语句
	// 参数 args: JOIN查询参数
	// 返回值: 支持链式调用的接口
	Joins(query string, args ...any) IActiveRecord[T]

	// InnerJoins 添加INNER JOIN连接查询
	// 参数 query: INNER JOIN查询语句
	// 参数 args: INNER JOIN查询参数
	// 返回值: 支持链式调用的接口
	InnerJoins(query string, args ...any) IActiveRecord[T]

	// 链式查询方法end
}

// BaseModel 基础模型结构体，提供通用的数据库操作功能
// 使用泛型T来支持不同的模型类型
// 实现 IActiveRecord 接口
// @property ID int `json:"id" example:"1"` 主键ID
// @property CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"` 创建时间
// @property UpdatedAt time.Time `json:"updated_at" example:"2023-01-02T00:00:00Z"` 更新时间
// @property DeletedAt *time.Time `json:"deleted_at,omitempty" example:"null"` 删除时间（软删除）
// 只有Unscoped()方法能查询到包括软删除记录在内的所有记录。
type BaseModel[T any] struct {
	// ID 主键ID
	ID int `gorm:"primarykey" json:"id"`

	// CreatedAt 记录创建时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt 记录更新时间
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt 记录删除时间，用于软删除功能
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// queryHandler GORM数据库查询处理器
	// 每次链式操作时，都会创建新实例：避免了状态共享问题 实现并发安全
	queryHandler *gorm.DB `gorm:"-"`
}

// WithTx 将查询器与事务对象绑定
// 用于在事务中执行数据库操作，创建一个新的模型实例并绑定事务对象
// 参数 tx: 事务对象
// 返回值: 绑定了事务对象的新模型实例
func (c *BaseModel[T]) WithTx(tx *gorm.DB) *BaseModel[T] {
	// 复制当前查询器的所有属性（仅替换 db 为事务对象 tx）
	newQuery := *c
	newQuery.queryHandler = tx
	return &newQuery
}

// getQuery 返回一个 *gorm.DB 实例，用于构建查询条件
// 如果已存在queryHandler则直接返回，否则使用全局数据库连接
// 返回值: GORM数据库查询对象
func (c *BaseModel[T]) getQuery() *gorm.DB {
	if c.queryHandler != nil {
		return c.queryHandler
	}
	// 默认使用全局的 database.DB
	return database.DB
}

// Create 创建一条新记录
// 使用gorm的Create方法将对象保存到数据库
// 参数 item: 要创建的记录对象
// 返回值: 创建过程中可能发生的错误
func (c *BaseModel[T]) Create(item *T) error {
	return c.getQuery().Create(item).Error
}

// Read 根据ID读取一条记录
// 使用gorm的First方法根据主键ID查询记录
// 参数 id: 要读取记录的ID
// 返回值: 读取到的记录对象和可能发生的错误
func (c *BaseModel[T]) Read(id int) (*T, error) {
	var item T
	err := c.getQuery().First(&item, id).Error
	return &item, err
}

// Update 更新一条记录
// 使用gorm的Save方法更新对象，会更新所有字段
// 使用gorm的Updates方法更新对象，只更新非零值字段
// 参数 item: 包含更新数据的记录对象
// 返回值: 更新过程中可能发生的错误
func (c *BaseModel[T]) Update(item *T) error {
	return c.getQuery().Updates(item).Error
}

// Delete 软删除：基于模型的主键，标记DeletedAt字段
// 使用gorm的Delete方法执行软删除操作
// 若需物理删除：改为 database.DB.Unscoped().Delete(model)
// 参数 id: 要删除记录的ID
// 返回值: 删除过程中可能发生的错误
func (c *BaseModel[T]) Delete(id int) error {
	return c.getQuery().Delete(new(T), id).Error
}

// More 查询多条记录
// 使用gorm的Find方法查询所有匹配的记录
// 返回值: 查询到的记录列表和可能发生的错误
func (c *BaseModel[T]) More() ([]T, error) {
	var items []T
	err := c.getQuery().Find(&items).Error
	return items, err
}

// One 查询单条记录
// 使用gorm的Take方法查询第一条匹配的记录
// 返回值: 查询到的记录对象和可能发生的错误
func (c *BaseModel[T]) One() (*T, error) {
	var item T
	err := c.getQuery().Take(&item).Error
	return &item, err
}

// Count 统计记录数
// 使用gorm的Count方法统计匹配条件的记录总数
// 返回值: 记录总数和可能发生的错误
func (c *BaseModel[T]) Count() (int, error) {
	var count int64
	err := c.getQuery().Model(new(T)).Count(&count).Error
	return int(count), err
}

// Page 分页查询记录
// 使用gorm的Offset和Limit方法实现分页查询
// page: 页数，从1开始
// pageSize: 每页记录数
// 返回值: 查询到的记录列表和总记录数
func (c *BaseModel[T]) Page(page, pageSize int) ([]T, int) {
	var items []T
	// 确保页数从1开始
	if page < 1 {
		page = 1
	}

	// 确保每页记录数大于0
	if pageSize < 1 {
		pageSize = 10 // 默认每页10条记录
	}

	total, err := c.Count()
	if (err != nil) || (total == 0) {
		return items, total
	}

	// 计算偏移量 (页数-1)*每页记录数
	offset := (page - 1) * pageSize
	c.getQuery().Offset(offset).Limit(pageSize).Find(&items)
	return items, total
}

// ------------------------------ gorm核心查询方法包装（链式） ------------------------------

// Where 添加WHERE条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: 查询条件
// 参数 args: 查询条件参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Where(query any, args ...any) IActiveRecord[T] {
	// 1. 基于原 db 构建新的 DB 实例（带 where 条件）
	newDB := c.getQuery().Where(query, args...)
	// 2. 复制原查询器，替换 db 为新实例
	newQuery := *c
	newQuery.queryHandler = newDB
	// 3. 返回新查询器
	return &newQuery
}

// Select 指定要查询的字段
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: 要查询的字段
// 参数 args: 查询字段参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Select(query any, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Select(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Order 添加排序条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 value: 排序条件
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Order(value any) IActiveRecord[T] {
	newDB := c.getQuery().Order(value)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Preload 预加载关联数据
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: 关联查询语句
// 参数 args: 关联查询参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Preload(query string, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Preload(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Not 添加NOT条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: NOT条件
// 参数 args: NOT条件参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Not(query any, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Not(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Or 添加OR条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: OR条件
// 参数 args: OR条件参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Or(query any, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Or(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Limit 限制返回记录数
// 通过创建新实例实现链式调用，确保并发安全
// 参数 limit: 限制的记录数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Limit(limit int) IActiveRecord[T] {
	newDB := c.getQuery().Limit(limit)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Group 添加分组条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: 分组条件
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Group(query string) IActiveRecord[T] {
	newDB := c.getQuery().Group(query)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Having 添加分组过滤条件
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: 分组过滤条件
// 参数 args: 分组过滤条件参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Having(query any, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Having(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Joins 添加JOIN连接查询
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: JOIN查询语句
// 参数 args: JOIN查询参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) Joins(query string, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().Joins(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// InnerJoins 添加INNER JOIN连接查询
// 通过创建新实例实现链式调用，确保并发安全
// 参数 query: INNER JOIN查询语句
// 参数 args: INNER JOIN查询参数
// 返回值: 支持链式调用的接口
func (c *BaseModel[T]) InnerJoins(query string, args ...any) IActiveRecord[T] {
	newDB := c.getQuery().InnerJoins(query, args...)
	newQuery := *c
	newQuery.queryHandler = newDB
	return &newQuery
}

// Transaction 执行事务操作
// 提供快捷入口，也可直接使用 database.DB.Transaction
// 修复：使用当前绑定的数据库连接，尊重已有的事务上下文
// 注意：错误传播控制需在回调函数中实现，内层事务回调返回nil可隔绝错误传递，外层事务不会回滚
func (c *BaseModel[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	// 使用当前绑定的数据库连接，而不是直接使用全局连接
	// 这样可以尊重通过WithTx绑定的事务，实现正确的事务嵌套
	return c.getQuery().Transaction(fc, opts...)
}
