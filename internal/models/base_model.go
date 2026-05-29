package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"webgos/internal/database"

	"gorm.io/gorm"
)

// IChainQuery 链式查询接口
type IChainQuery[T any] interface {
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

	// Scopes 添加多个WHERE条件
	// 封装可复用的数据库查询逻辑
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) IActiveRecord[T]

	// Table 指定表名
	// 参数 tableName: 表名
	// 参数 args: 表名参数
	// 返回值: 支持链式调用的接口
	Table(tableName string, args ...any) IActiveRecord[T]
}

// IReader 读操作接口
type IReader[T any] interface {
	// Read 根据ID读取一条记录
	// 使用gorm的First方法根据主键ID查询记录
	// 当记录不存在时，返回 gorm.ErrRecordNotFound 错误
	//
	// 示例:
	//    user, err := userModel.Read(1)
	//    if err != nil {
	//        if errors.Is(err, gorm.ErrRecordNotFound) {
	//            // 处理记录不存在的情况
	//            fmt.Println("User not found")
	//        } else {
	//            // 处理其他数据库错误
	//            log.Error("Database error:", err)
	//        }
	//        return
	//    }
	//    fmt.Printf("Found user: %+v\n", user)
	//
	// 参数 id: 要读取记录的ID
	// 返回值:
	//   - *T: 读取到的记录对象指针
	//   - error: 可能的错误，包括 gorm.ErrRecordNotFound
	Read(id int) (*T, error)

	// One 查询单条记录
	// 使用gorm的Take方法查询第一条匹配的记录
	//
	// 示例:
	//    // 查询第一条用户记录
	//    user, err := userModel.One()
	//    if err != nil {
	//        if errors.Is(err, gorm.ErrRecordNotFound) {
	//            fmt.Println("No users found")
	//        } else {
	//            log.Error("Failed to query user:", err)
	//        }
	//        return
	//    }
	//    fmt.Printf("First user: %+v\n", user)
	//
	//    // 查询满足条件的第一条用户记录
	//    adminUser, err := userModel.Where("role = ?", "admin").One()
	//    if err != nil {
	//        if errors.Is(err, gorm.ErrRecordNotFound) {
	//            fmt.Println("No admin users found")
	//        } else {
	//            log.Error("Failed to query admin user:", err)
	//        }
	//        return
	//    }
	//    fmt.Printf("First admin user: %+v\n", adminUser)
	//
	// 注意事项:
	//   - Take 方法查询到一条记录后就会停止查询
	//   - 如果没有满足条件的记录，会返回 gorm.ErrRecordNotFound 错误
	//
	// 返回值: 查询到的记录对象和可能发生的错误
	One() (*T, error)

	// More 查询多条记录
	// 使用gorm的Find方法查询所有匹配的记录
	//
	// 示例:
	//    // 查询所有用户
	//    users, err := userModel.More()
	//    if err != nil {
	//        log.Error("Failed to query users:", err)
	//        return
	//    }
	//    fmt.Printf("Found %d users\n", len(users))
	//
	//    // 查询满足条件的用户
	//    activeUsers, err := userModel.Where("active = ?", true).More()
	//    if err != nil {
	//        log.Error("Failed to query active users:", err)
	//        return
	//    }
	//    fmt.Printf("Found %d active users\n", len(activeUsers))
	//
	// 注意事项:
	//   - 会查询所有满足条件的记录，大量数据时注意性能
	//   - 对于大量数据，建议使用 Page 方法分页查询
	//
	// 返回值: 查询到的记录列表和可能发生的错误
	More() ([]T, error)

	// Count 统计记录总数
	// 返回值: 记录总数和可能发生的错误
	Count() (int, error)

	// Exist 检查记录是否存在
	// 使用gorm的Take方法检查是否存在匹配的记录
	// 不会返回 gorm.ErrRecordNotFound 错误，而是将其转换为 false 返回值
	//
	// 示例:
	//    // 检查特定ID的用户是否存在
	//    exists, err := userModel.Where("id = ?", 1).Exist()
	//    if err != nil {
	//        log.Error("Database error:", err)
	//        return
	//    }
	//
	//    if exists {
	//        fmt.Println("User exists")
	//    } else {
	//        fmt.Println("User does not exist")
	//    }
	//
	// 边界情况:
	//   - 当查询条件匹配多条记录时，只检查是否存在至少一条记录
	//   - 数据库连接错误会返回错误
	//
	// 返回值:
	//   - bool: 是否存在匹配的记录
	//   - error: 数据库查询过程中可能发生的错误（不包括记录不存在）
	Exist() (bool, error)

	// Page 分页查询记录
	// 使用gorm的Offset和Limit方法实现分页查询
	// 自动处理页数和页面大小的边界情况
	//
	// 示例:
	//    // 查询第一页，每页10条记录
	//    users, total, err := userModel.Page(1, 10)
	//    if err != nil {
	//        log.Error("Failed to query users:", err)
	//        return
	//    }
	//
	//    fmt.Printf("Total users: %d\n", total)
	//    fmt.Printf("Current page users: %d\n", len(users))
	//    for _, user := range users {
	//        fmt.Printf("User: %+v\n", user)
	//    }
	//
	// 边界情况:
	//   - page < 1 时自动设置为 1
	//   - pageSize < 1 时自动设置为 10
	//   - pageSize > 1000 时自动设置为 1000
	//   - 当 offset 超过总记录数时，返回空列表
	//
	// 参数:
	//   - page: 页数，从1开始
	//   - pageSize: 每页记录数
	//
	// 返回值:
	//   - []T: 查询到的记录列表
	//   - int: 总记录数
	//   - error: 查询过程中可能发生的错误
	Page(page, pageSize int) ([]T, int, error)

	// Pluck 查询单个字段的值
	// 使用gorm的Pluck方法查询指定字段的值
	//
	// 示例:
	// var names []string
	// err := model.Pluck("name", &names)
	//
	// 参数 column: 要查询的字段名
	// 参数 dest: 接收查询结果的目标变量指针
	// 返回值: 字段值列表和可能发生的错误
	Pluck(column string, dest any) error

	// Find 查询多条记录
	// 参数 dest: 接收查询结果的目标变量指针
	// 参数 conds: 查询条件
	// 返回值: 查询过程中可能发生的错误
	Find(dest any, conds ...any) error

	// First 查询第一条记录
	// 参数 dest: 接收查询结果的目标变量指针
	// 参数 conds: 查询条件
	// 返回值: 查询过程中可能发生的错误
	First(dest any, conds ...any) error

	// FirstOrCreate 获取或创建记录
	// 如果记录存在则获取第一条匹配的记录，否则创建新记录
	//
	// 示例:
	//    // 查找或创建用户
	//    user := &User{Name: "John"}
	//    err := userModel.Where("name = ?", "John").FirstOrCreate(user)
	//    if err != nil {
	//        log.Error("Failed to find or create user:", err)
	//        return
	//    }
	//
	//    if user.ID == 0 {
	//        fmt.Println("Created new user")
	//    } else {
	//        fmt.Println("Found existing user")
	//    }
	//
	// 注意事项:
	//   - 会根据查询条件查找记录，如果不存在则创建
	//   - 创建时会使用传入的结构体作为默认值
	//
	// 参数 dest: 要创建的记录对象
	// 参数 conds: 查询条件，用于查找记录是否存在
	// 返回值: 查询到或创建的记录对象和可能发生的错误
	FirstOrCreate(dest any, conds ...any) error
}

// IWriter 写操作接口
type IWriter[T any] interface {
	// Create 创建一条新记录
	// 参数 item: 要创建的记录对象
	// 返回值: 创建过程中可能发生的错误
	Create(item *T) error

	// Update 更新一条对象记录（不包含零值字段）
	// 参数 item: 包含更新数据的记录对象
	// 返回值: 更新过程中可能发生的错误
	// 只会更新非零值字段，如果要更新全部字段，请使加上Select("*")
	Update(item *T) error
	// Updates 更新一条对象记录（包含零值字段）
	Updates(item *T) error

	// Delete 根据ID删除一条记录（软删除）
	// 参数 id: 要删除记录的ID
	// 返回值: 删除过程中可能发生的错误
	Delete(id int) error

	// BatchCreate 批量创建记录
	// 使用gorm的CreateInBatches方法批量创建记录
	//
	// 示例:
	//    // 批量创建用户
	//    users := []User{
	//        {Name: "John", Email: "john@example.com"},
	//        {Name: "Jane", Email: "jane@example.com"},
	//        {Name: "Bob", Email: "bob@example.com"},
	//    }
	//
	//    // 每批创建100条记录
	//    err := userModel.BatchCreate(users, 100)
	//    if err != nil {
	//        log.Error("Failed to batch create users:", err)
	//        return
	//    }
	//    fmt.Println("Users created successfully")
	//
	// 注意事项:
	//   - 会按照指定的批次大小分批创建记录
	//   - 每一批次都在单独的事务中执行
	//   - 如果某一批次失败，该批次之前的数据仍然会被提交
	//
	// 参数 items: 要创建的记录对象列表
	// batchSize: 每批创建的记录数量，默认为100
	// 返回值: 创建过程中可能发生的错误
	BatchCreate(items []T, batchSize int) error

	// UpdateColumns 更新指定字段 多字段
	// 使用gorm的UpdateColumns方法更新指定字段
	// 参数 columns: 要更新的字段及其新值的映射
	// 返回值: 更新过程中可能发生的错误
	// 如果存在 WHERE 条件，则使用该条件更新记录
	// 否则，使用模型主键ID作为WHERE条件
	UpdateColumns(columns map[string]any) error

	// UpdateColumn 更新指定字段 单字段
	// 使用gorm的UpdateColumn方法更新指定字段
	// 参数 column: 要更新的字段名
	// 参数 value: 新值
	// 返回值: 更新过程中可能发生的错误
	UpdateColumn(column string, value any) error
}

// ITransaction 事务与上下文接口
type ITransaction[T any] interface {
	// WithTx 将查询器与事务对象绑定，用于在事务中执行数据库操作
	// 用于在事务中执行数据库操作，创建一个新的模型实例并绑定事务对象
	//
	// 示例:
	//    // 在事务中执行操作
	//    err := database.GetDB().Transaction(func(tx *gorm.DB) error {
	//        // 绑定事务对象到模型
	//        userTxModel := userModel.WithTx(tx)
	//
	//        // 在事务中创建用户
	//        user := &User{Name: "John"}
	//        if err := userTxModel.Create(user); err != nil {
	//            return err // 回滚事务
	//        }
	//
	//        // 在事务中创建用户配置
	//        profile := &Profile{UserID: user.ID, Bio: "Hello"}
	//        if err := userTxModel.WithTx(tx).Create(profile); err != nil {
	//            return err // 回滚事务
	//        }
	//
	//        return nil // 提交事务
	//    })
	//
	// 注意事项:
	//   - 每次调用都会创建新的实例，保证并发安全
	//   - 事务对象必须由 GORM 提供
	//
	// 参数 tx: 事务对象
	// 返回值: 绑定了事务对象的新模型实例
	WithTx(tx *gorm.DB) IActiveRecord[T]

	// WithCtx 将查询器与上下文绑定
	// 支持超时控制和取消操作
	// 参数 ctx: 上下文对象
	// 返回值: 绑定了上下文的新模型实例
	WithCtx(ctx context.Context) IActiveRecord[T]

	// Transaction 在事务中执行数据库操作
	// 使用gorm的Transaction方法包装事务操作
	// 支持事务嵌套，如果已经在事务中，则使用当前事务
	//
	// 示例:
	//    // 在事务中执行多个操作
	//    err := userModel.Transaction(func(tx *gorm.DB) error {
	//        // 创建用户
	//        user := &User{Name: "John"}
	//        if err := tx.Create(user).Error; err != nil {
	//            return err // 回滚事务
	//        }
	//        user.WithTx(tx).More(); // 使用事务查询用户列表
	//
	//        // 创建用户配置
	//        profile := &Profile{UserID: user.ID, Bio: "Hello"}
	//        if err := tx.Create(profile).Error; err != nil {
	//            return err // 回滚事务
	//        }
	//
	//        return nil // 提交事务
	//    })
	//
	//    if err != nil {
	//        log.Error("Transaction failed:", err)
	//    } else {
	//        fmt.Println("Transaction succeeded")
	//    }
	//
	// 边界情况:
	//   - 如果在事务回调函数中返回 nil，则提交事务
	//   - 如果在事务回调函数中返回错误，则回滚事务
	//   - 支持事务嵌套，内层事务的回滚不会影响外层事务
	//
	// 参数:
	//   - fc: 事务执行函数，包含在事务中执行的业务逻辑
	//   - opts: 事务选项
	//
	// 返回值: 执行事务过程中可能发生的错误
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error

	// Unscoped 查询包含软删除的记录
	// 默认情况下，查询会自动过滤已软删除的记录
	// 调用此方法后，可以查询到包括软删除记录在内的所有记录
	//
	// 示例:
	//    // 查询已删除的用户
	//    deletedUser, err := userModel.Unscoped().Where("id = ?", 1).One()
	//    if err != nil {
	//        log.Error("Failed to query deleted user:", err)
	//        return
	//    }
	//    fmt.Printf("Deleted user: %+v\n", deletedUser)
	//
	// 注意事项:
	//   - 谨慎使用，避免误操作已删除数据
	//   - 常用于数据恢复、审计等场景
	Unscoped() IActiveRecord[T]
}

// IActiveRecord 组合接口，包含所有数据库操作能力
type IActiveRecord[T any] interface {
	IChainQuery[T]
	IReader[T]
	IWriter[T]
	ITransaction[T]

	// HasWhere 检查是否有WHERE条件
	// 返回值: 是否有WHERE条件
	HasWhere() bool
}

// BaseModel 基础模型，提供通用数据库操作
// 通过泛型 T 支持不同模型类型，实现 IActiveRecord 接口
// 链式操作每次创建新实例，保证并发安全
//
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

// chain 链式操作辅助函数，创建新实例并替换 queryHandler
// 参数 newDB: 新的 GORM 数据库查询实例
// 返回值: 绑定了新查询实例的模型
func (c *BaseModel[T]) chain(newDB *gorm.DB) IActiveRecord[T] {
	next := *c
	next.queryHandler = newDB
	return &next
}

// getQuery 返回 *gorm.DB 实例，优先使用已绑定的 queryHandler，否则使用全局数据库连接
// 返回值: GORM数据库查询对象
func (c *BaseModel[T]) getQuery() *gorm.DB {
	if c.queryHandler != nil {
		return c.queryHandler
	}
	return database.GetDB()
}

// ---------- ITransaction ----------

// WithTx 将查询器与事务对象绑定
// 复制当前查询器的所有属性，仅替换 db 为事务对象 tx
func (c *BaseModel[T]) WithTx(tx *gorm.DB) IActiveRecord[T] {
	return c.chain(tx)
}

// WithCtx 将查询器与上下文绑定
// 支持超时控制和取消操作
func (c *BaseModel[T]) WithCtx(ctx context.Context) IActiveRecord[T] {
	return c.chain(c.getQuery().WithContext(ctx))
}

// Unscoped 查询包含软删除的记录
func (c *BaseModel[T]) Unscoped() IActiveRecord[T] {
	return c.chain(c.getQuery().Unscoped())
}

// Transaction 在事务中执行数据库操作
// 使用当前绑定的数据库连接，而不是直接使用全局连接
// 这样可以尊重通过WithTx绑定的事务，实现正确的事务嵌套
func (c *BaseModel[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return c.getQuery().Transaction(fc, opts...)
}

// ---------- IWriter ----------

// Create 创建一条新记录
func (c *BaseModel[T]) Create(item *T) error {
	return c.getQuery().Create(item).Error
}

// BatchCreate 批量创建记录
// 如果 batchSize <= 0，默认使用 100
func (c *BaseModel[T]) BatchCreate(items []T, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100
	}
	return c.getQuery().CreateInBatches(items, batchSize).Error
}

// Update 更新操作（不包含零值字段）
func (c *BaseModel[T]) Update(item *T) error {
	return c.getQuery().Updates(item).Error
}

// Updates 更新操作（包含零值字段）
func (c *BaseModel[T]) Updates(item *T) error {
	return c.getQuery().Select("*").Updates(item).Error
}

// UpdateColumns 更新指定字段
// 如果存在 WHERE 条件，则使用该条件更新记录
// 否则，使用模型主键ID作为WHERE条件
func (c *BaseModel[T]) UpdateColumns(columns map[string]any) error {
	if c.HasWhere() {
		return c.getQuery().Model((*T)(nil)).UpdateColumns(columns).Error
	}
	if c.ID == 0 {
		return errors.New("no id found")
	}
	return c.getQuery().Model((*T)(nil)).Where("id = ?", c.ID).UpdateColumns(columns).Error
}

// UpdateColumn 更新单个字段
// 如果存在 WHERE 条件，则使用该条件更新记录
// 否则，使用模型主键ID作为WHERE条件
func (c *BaseModel[T]) UpdateColumn(column string, value any) error {
	if c.HasWhere() {
		return c.getQuery().Model((*T)(nil)).UpdateColumn(column, value).Error
	}
	if c.ID == 0 {
		return errors.New("no id found")
	}
	return c.getQuery().Model((*T)(nil)).Where("id = ?", c.ID).UpdateColumn(column, value).Error
}

// Delete 根据ID删除一条记录（软删除）
func (c *BaseModel[T]) Delete(id int) error {
	return c.getQuery().Delete(new(T), id).Error
}

// ---------- IReader ----------

// Read 根据ID读取一条记录
func (c *BaseModel[T]) Read(id int) (*T, error) {
	var item T
	err := c.getQuery().First(&item, id).Error
	return &item, err
}

// One 查询单条记录，未找到时返回 gorm.ErrRecordNotFound
func (c *BaseModel[T]) One() (*T, error) {
	var item T
	err := c.getQuery().Take(&item).Error
	return &item, err
}

// More 查询多条记录，大量数据时建议使用 Page 分页查询
func (c *BaseModel[T]) More() ([]T, error) {
	var items []T
	err := c.getQuery().Find(&items).Error
	return items, err
}

// Count 统计记录总数
// 使用 `(*T)(nil)` 作为模型类型占位，避免分配一个实际的零值对象。
// GORM 仅使用传入值的类型信息来确定表名/模型元信息，因此传入一个类型为 `*T` 的 nil 指针
// 能够达到同样目的并略微减少一次内存分配。
func (c *BaseModel[T]) Count() (int, error) {
	var count int64
	err := c.getQuery().Model((*T)(nil)).Count(&count).Error
	return int(count), err
}

// Exist 检查记录是否存在，将 gorm.ErrRecordNotFound 转换为 false
func (c *BaseModel[T]) Exist() (bool, error) {
	var item T
	err := c.getQuery().Take(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Page 分页查询，page 从 1 开始，pageSize 范围 [1, 1000]
func (c *BaseModel[T]) Page(page, pageSize int) ([]T, int, error) {
	var items []T
	var total int64
	// 确保页数从1开始
	if page < 1 {
		page = 1
	}

	// 确保每页记录数大于0且不超过1000
	if pageSize < 1 {
		pageSize = 10 // 默认每页10条记录
	} else if pageSize > 1000 {
		pageSize = 1000 // 最大每页1000条记录
	}
	// 计算偏移量 (页数-1)*每页记录数
	offset := (page - 1) * pageSize
	// 链式调用Count() - 需要明确指定Model
	err := c.getQuery().Model((*T)(nil)).Count(&total).Offset(offset).Limit(pageSize).Find(&items).Error
	return items, int(total), err
}

// Pluck 查询单个字段的值
func (c *BaseModel[T]) Pluck(column string, dest any) error {
	return c.getQuery().Model((*T)(nil)).Pluck(column, dest).Error
}

// Find 查询多条记录
func (c *BaseModel[T]) Find(dest any, conds ...any) error {
	return c.getQuery().Find(dest, conds...).Error
}

// First 查询第一条记录
func (c *BaseModel[T]) First(dest any, conds ...any) error {
	return c.getQuery().First(dest, conds...).Error
}

// FirstOrCreate 获取或创建记录
func (c *BaseModel[T]) FirstOrCreate(dest any, conds ...any) error {
	return c.getQuery().FirstOrCreate(dest, conds...).Error
}

// ---------- IChainQuery ----------

// Where 添加WHERE条件
func (c *BaseModel[T]) Where(query any, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Where(query, args...))
}

// Select 指定要查询的字段
func (c *BaseModel[T]) Select(query any, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Select(query, args...))
}

// Preload 预加载关联数据
func (c *BaseModel[T]) Preload(query string, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Preload(query, args...))
}

// Order 添加排序条件
func (c *BaseModel[T]) Order(value any) IActiveRecord[T] {
	return c.chain(c.getQuery().Order(value))
}

// Not 添加NOT条件
func (c *BaseModel[T]) Not(query any, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Not(query, args...))
}

// Or 添加OR条件
func (c *BaseModel[T]) Or(query any, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Or(query, args...))
}

// Limit 限制返回记录数
func (c *BaseModel[T]) Limit(limit int) IActiveRecord[T] {
	return c.chain(c.getQuery().Limit(limit))
}

// Group 添加分组条件
func (c *BaseModel[T]) Group(query string) IActiveRecord[T] {
	return c.chain(c.getQuery().Group(query))
}

// Having 添加分组过滤条件
func (c *BaseModel[T]) Having(query any, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Having(query, args...))
}

// Joins 添加JOIN连接查询
func (c *BaseModel[T]) Joins(query string, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Joins(query, args...))
}

// InnerJoins 添加INNER JOIN连接查询
func (c *BaseModel[T]) InnerJoins(query string, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().InnerJoins(query, args...))
}

// Scopes 添加多个WHERE条件，封装可复用的数据库查询逻辑
func (c *BaseModel[T]) Scopes(funcs ...func(*gorm.DB) *gorm.DB) IActiveRecord[T] {
	return c.chain(c.getQuery().Scopes(funcs...))
}

// Table 指定表名
func (c *BaseModel[T]) Table(tableName string, args ...any) IActiveRecord[T] {
	return c.chain(c.getQuery().Table(tableName, args...))
}

// ---------- 其他 ----------

// HasWhere 检查是否有 WHERE 条件
// 检查Statement是否为空，以及WHERE子句是否存在且有表达式
func (c *BaseModel[T]) HasWhere() bool {
	stmt := c.getQuery().Statement
	if stmt == nil {
		return false
	}
	if whereClause, ok := stmt.Clauses["WHERE"]; ok {
		return whereClause.Expression != nil
	}
	return false
}
