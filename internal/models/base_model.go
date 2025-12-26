package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"webgos/internal/database"

	"gorm.io/gorm"
)

// IActiveRecord 定义了基础模型的接口，提供通用的数据库操作方法
// T 是泛型参数，代表具体的模型类型
type IActiveRecord[T any] interface {
	// WithTx 将查询器与事务对象绑定，用于在事务中执行数据库操作
	// 用于在事务中执行数据库操作，创建一个新的模型实例并绑定事务对象
	//
	// 示例:
	//    // 在事务中执行操作
	//    err := database.DB.Transaction(func(tx *gorm.DB) error {
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
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)

	// Create 创建一条新记录
	// 参数 item: 要创建的记录对象
	// 返回值: 创建过程中可能发生的错误
	Create(item *T) error

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

	// Update 更新一条记录
	// 参数 item: 包含更新数据的记录对象
	// 返回值: 更新过程中可能发生的错误
	// 只会更新非零值字段，如果要更新全部字段，请使加上Select("*")
	Update(item *T) error

	// UpdateColumns 更新指定字段
	// 使用gorm的UpdateColumns方法更新指定字段
	// 参数 columns: 要更新的字段及其新值的映射
	// 返回值: 更新过程中可能发生的错误
	// 如果存在 WHERE 条件，则使用该条件更新记录
	// 否则，使用模型主键ID作为WHERE条件
	UpdateColumns(columns map[string]any) error

	// Delete 根据ID删除一条记录（软删除）
	// 参数 id: 要删除记录的ID
	// 返回值: 删除过程中可能发生的错误
	Delete(id int) error

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

	// Pluck 查询单个字段的值
	// 使用gorm的Pluck方法查询指定字段的值
	//
	// 示例:
	//    // 查询所有用户的名称
	//    names, err := userModel.Pluck("name")
	//    if err != nil {
	//        log.Error("Failed to pluck names:", err)
	//        return
	//    }
	//
	//    for _, name := range names {
	//        fmt.Println(name)
	//    }
	//
	//    // 查询满足条件的用户ID
	//    ids, err := userModel.Where("active = ?", true).Pluck("id")
	//    if err != nil {
	//        log.Error("Failed to pluck IDs:", err)
	//        return
	//    }
	//
	//    for _, id := range ids {
	//        fmt.Println(id)
	//    }
	//
	// 注意事项:
	//   - 只查询指定字段的值，减少数据传输和内存占用
	//   - 返回的是 interface{} 类型的切片，需要类型断言才能使用
	//
	// 参数 column: 要查询的字段名
	// 返回值: 字段值列表和可能发生的错误
	Pluck(column string) ([]interface{}, error)

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
	// 参数 item: 要创建的记录对象
	// 返回值: 查询到或创建的记录对象和可能发生的错误
	FirstOrCreate(item *T) error

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

// WithTx 将查询器与事务对象绑定，用于在事务中执行数据库操作
// 用于在事务中执行数据库操作，创建一个新的模型实例并绑定事务对象
//
// 示例:
//
//	// 在事务中执行操作
//	err := database.DB.Transaction(func(tx *gorm.DB) error {
//	    // 绑定事务对象到模型
//	    userTxModel := userModel.WithTx(tx)
//
//	    // 在事务中创建用户
//	    user := &User{Name: "John"}
//	    if err := userTxModel.Create(user); err != nil {
//	        return err // 回滚事务
//	    }
//
//	    // 在事务中创建用户配置
//	    profile := &Profile{UserID: user.ID, Bio: "Hello"}
//	    if err := userTxModel.WithTx(tx).Create(profile); err != nil {
//	        return err // 回滚事务
//	    }
//
//	    return nil // 提交事务
//	})
//
// 注意事项:
//   - 每次调用都会创建新的实例，保证并发安全
//   - 事务对象必须由 GORM 提供
//
// 参数 tx: 事务对象
// 返回值: 绑定了事务对象的新模型实例
func (c *BaseModel[T]) WithTx(tx *gorm.DB) IActiveRecord[T] {
	// 复制当前查询器的所有属性（仅替换 db 为事务对象 tx）
	newQuery := *c
	newQuery.queryHandler = tx
	return &newQuery
}

// WithCtx 将查询器与上下文绑定
// 支持超时控制和取消操作
// 创建一个新的模型实例并绑定上下文，不影响原实例
//
// 示例:
//
//	// 创建一个5秒超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// 使用带超时的查询
//	user, err := userModel.WithCtx(ctx).Read(1)
//	if err != nil {
//	    if errors.Is(err, context.DeadlineExceeded) {
//	        fmt.Println("Query timeout")
//	    } else {
//	        fmt.Printf("Query error: %v\n", err)
//	    }
//	    return
//	}
//
//	fmt.Printf("Found user: %+v\n", user)
//
// 注意事项:
//   - 每次调用都会创建新的实例，保证并发安全
//   - 上下文的取消操作会影响后续的数据库操作
//
// 参数 ctx: 上下文对象
// 返回值: 绑定了上下文的新模型实例
func (c *BaseModel[T]) WithCtx(ctx context.Context) IActiveRecord[T] {
	newQuery := *c
	newQuery.queryHandler = c.getQuery().WithContext(ctx)
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
//
// 示例:
//
//	user := &User{Name: "John", Email: "john@example.com"}
//	err := userModel.Create(user)
//	if err != nil {
//	    log.Error("Failed to create user:", err)
//	    return
//	}
//	fmt.Printf("Created user with ID: %d\n", user.ID)
//
// 注意事项:
//   - 创建成功后，会自动填充记录的ID和时间戳字段
//   - 如果违反数据库约束（如唯一索引），会返回相应错误
//
// 参数 item: 要创建的记录对象
// 返回值: 创建过程中可能发生的错误
func (c *BaseModel[T]) Create(item *T) error {
	return c.getQuery().Create(item).Error
}

// BatchCreate 批量创建记录
// 使用gorm的CreateInBatches方法批量创建记录
//
// 示例:
//
//	// 批量创建用户
//	users := []User{
//	    {Name: "John", Email: "john@example.com"},
//	    {Name: "Jane", Email: "jane@example.com"},
//	    {Name: "Bob", Email: "bob@example.com"},
//	}
//
//	// 每批创建100条记录
//	err := userModel.BatchCreate(users, 100)
//	if err != nil {
//	    log.Error("Failed to batch create users:", err)
//	    return
//	}
//	fmt.Println("Users created successfully")
//
// 注意事项:
//   - 会按照指定的批次大小分批创建记录
//   - 每一批次都在单独的事务中执行
//   - 如果某一批次失败，该批次之前的数据仍然会被提交
//
// 参数 items: 要创建的记录对象列表
// batchSize: 每批创建的记录数量，默认为100
// 返回值: 创建过程中可能发生的错误
func (c *BaseModel[T]) BatchCreate(items []T, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100 // 默认批次大小
	}
	return c.getQuery().CreateInBatches(items, batchSize).Error
}

// Read 根据ID读取一条记录
// 使用gorm的First方法根据主键ID查询记录
// 当记录不存在时，返回 gorm.ErrRecordNotFound 错误
//
// 示例:
//
//	user, err := userModel.Read(1)
//	if err != nil {
//	    if errors.Is(err, gorm.ErrRecordNotFound) {
//	        // 处理记录不存在的情况
//	        fmt.Println("User not found")
//	    } else {
//	        // 处理其他数据库错误
//	        log.Error("Database error:", err)
//	    }
//	    return
//	}
//	fmt.Printf("Found user: %+v\n", user)
//
// 参数 id: 要读取记录的ID
// 返回值:
//   - *T: 读取到的记录对象指针
//   - error: 可能的错误，包括 gorm.ErrRecordNotFound
func (c *BaseModel[T]) Read(id int) (*T, error) {
	var item T
	err := c.getQuery().First(&item, id).Error
	return &item, err
}

// Update 更新一条记录
// 使用gorm的Save方法更新对象，会更新所有字段
// 使用gorm的Updates方法更新对象，只更新非零值字段
//
// 示例:
//
//	// 先查询出要更新的记录
//	user, err := userModel.Read(1)
//	if err != nil {
//	    log.Error("Failed to find user:", err)
//	    return
//	}
//
//	// 修改记录字段
//	user.Name = "New Name"
//
//	// 执行更新
//	err = userModel.Update(user)
//	if err != nil {
//	    log.Error("Failed to update user:", err)
//	    return
//	}
//	fmt.Println("User updated successfully")
//
// 注意事项:
//   - 使用 Updates 方法，只会更新非零值字段
//   - 如果需要更新所有字段（包括零值），可以使用 Save 方法
//
// 参数 item: 包含更新数据的记录对象
// 返回值: 更新过程中可能发生的错误
func (c *BaseModel[T]) Update(item *T) error {
	return c.getQuery().Updates(item).Error
}

// UpdateColumns 更新指定字段
// 如果存在 WHERE 条件，则使用该条件更新记录
// 否则，使用模型主键ID作为WHERE条件
func (c *BaseModel[T]) UpdateColumns(columns map[string]any) error {
	if c.hasWhere() {
		return c.getQuery().Model((*T)(nil)).UpdateColumns(columns).Error
	} else {
		if c.ID == 0 {
			return errors.New("no id found")
		}
		return c.getQuery().Model((*T)(nil)).Where("id = ?", c.ID).UpdateColumns(columns).Error
	}
}
func (c *BaseModel[T]) hasWhere() bool {
	stmt := c.getQuery().Statement
	// 检查Statement是否为空
	if stmt == nil {
		return false
	}
	// 检查WHERE子句是否存在且有表达式
	if whereClause, ok := stmt.Clauses["WHERE"]; ok {
		return whereClause.Expression != nil
	}
	return false
}

// Delete 软删除：基于模型的主键，标记DeletedAt字段
// 使用gorm的Delete方法执行软删除操作
// 若需物理删除：改为 database.DB.Unscoped().Delete(model)
//
// 示例:
//
//	// 软删除ID为1的用户
//	err := userModel.Delete(1)
//	if err != nil {
//	    log.Error("Failed to delete user:", err)
//	    return
//	}
//	fmt.Println("User deleted successfully")
//
//	// 物理删除示例（谨慎使用）
//	// err := database.DB.Unscoped().Delete(&User{}, 1).Error
//
// 注意事项:
//   - 默认执行软删除，只标记 DeletedAt 字段
//   - 软删除的记录在普通查询中不可见
//   - 如需查询软删除的记录，需要使用 Unscoped() 方法
//
// 参数 id: 要删除记录的ID
// 返回值: 删除过程中可能发生的错误
func (c *BaseModel[T]) Delete(id int) error {
	return c.getQuery().Delete(new(T), id).Error
}

// More 查询多条记录
// 使用gorm的Find方法查询所有匹配的记录
//
// 示例:
//
//	// 查询所有用户
//	users, err := userModel.More()
//	if err != nil {
//	    log.Error("Failed to query users:", err)
//	    return
//	}
//	fmt.Printf("Found %d users\n", len(users))
//
//	// 查询满足条件的用户
//	activeUsers, err := userModel.Where("active = ?", true).More()
//	if err != nil {
//	    log.Error("Failed to query active users:", err)
//	    return
//	}
//	fmt.Printf("Found %d active users\n", len(activeUsers))
//
// 注意事项:
//   - 会查询所有满足条件的记录，大量数据时注意性能
//   - 对于大量数据，建议使用 Page 方法分页查询
//
// 返回值: 查询到的记录列表和可能发生的错误
func (c *BaseModel[T]) More() ([]T, error) {
	var items []T
	err := c.getQuery().Find(&items).Error
	return items, err
}

// Pluck 查询单个字段的值
// 使用gorm的Pluck方法查询指定字段的值
//
// 示例:
//
//	// 查询所有用户的名称
//	names, err := userModel.Pluck("name")
//	if err != nil {
//	    log.Error("Failed to pluck names:", err)
//	    return
//	}
//
//	for _, name := range names {
//	    fmt.Println(name)
//	}
//
//	// 查询满足条件的用户ID
//	ids, err := userModel.Where("active = ?", true).Pluck("id")
//	if err != nil {
//	    log.Error("Failed to pluck IDs:", err)
//	    return
//	}
//
//	for _, id := range ids {
//	    fmt.Println(id)
//	}
//
// 注意事项:
//   - 只查询指定字段的值，减少数据传输和内存占用
//   - 返回的是 interface{} 类型的切片，需要类型断言才能使用
//
// 参数 column: 要查询的字段名
// 返回值: 字段值列表和可能发生的错误
func (c *BaseModel[T]) Pluck(column string) ([]interface{}, error) {
	var values []interface{}
	err := c.getQuery().Pluck(column, &values).Error
	return values, err
}

// One 查询单条记录
// 使用gorm的Take方法查询第一条匹配的记录
//
// 示例:
//
//	// 查询第一条用户记录
//	user, err := userModel.One()
//	if err != nil {
//	    if errors.Is(err, gorm.ErrRecordNotFound) {
//	        fmt.Println("No users found")
//	    } else {
//	        log.Error("Failed to query user:", err)
//	    }
//	    return
//	}
//	fmt.Printf("First user: %+v\n", user)
//
//	// 查询满足条件的第一条用户记录
//	adminUser, err := userModel.Where("role = ?", "admin").One()
//	if err != nil {
//	    if errors.Is(err, gorm.ErrRecordNotFound) {
//	        fmt.Println("No admin users found")
//	    } else {
//	        log.Error("Failed to query admin user:", err)
//	    }
//	    return
//	}
//	fmt.Printf("First admin user: %+v\n", adminUser)
//
// 注意事项:
//   - Take 方法查询到一条记录后就会停止查询
//   - 如果没有满足条件的记录，会返回 gorm.ErrRecordNotFound 错误
//
// 返回值: 查询到的记录对象和可能发生的错误
func (c *BaseModel[T]) One() (*T, error) {
	var item T
	err := c.getQuery().Take(&item).Error
	return &item, err
}

// FirstOrCreate 获取或创建记录
// 如果记录存在则获取第一条匹配的记录，否则创建新记录
//
// 示例:
//
//	// 查找或创建用户
//	user := &User{Name: "John"}
//	err := userModel.Where("name = ?", "John").FirstOrCreate(user)
//	if err != nil {
//	    log.Error("Failed to find or create user:", err)
//	    return
//	}
//
//	if user.ID == 0 {
//	    fmt.Println("Created new user")
//	} else {
//	    fmt.Println("Found existing user")
//	}
//
// 注意事项:
//   - 会根据查询条件查找记录，如果不存在则创建
//   - 创建时会使用传入的结构体作为默认值
//
// 参数 item: 要创建的记录对象
// 返回值: 查询到或创建的记录对象和可能发生的错误
func (c *BaseModel[T]) FirstOrCreate(item *T) error {
	return c.getQuery().FirstOrCreate(item).Error
}

// Count 统计记录数
// 使用gorm的Count方法统计匹配条件的记录总数
//
// 示例:
//
//	// 统计所有用户数量
//	count, err := userModel.Count()
//	if err != nil {
//	    log.Error("Failed to count users:", err)
//	    return
//	}
//	fmt.Printf("Total users: %d\n", count)
//
//	// 统计满足条件的用户数量
//	count, err = userModel.Where("active = ?", true).Count()
//	if err != nil {
//	    log.Error("Failed to count active users:", err)
//	    return
//	}
//	fmt.Printf("Active users: %d\n", count)
//
// 注意事项:
//   - 使用 `(*T)(nil)` 作为模型类型占位，避免分配一个实际的零值对象
//   - GORM 仅使用传入值的类型信息来确定表名/模型元信息
//
// 返回值: 记录总数和可能发生的错误
func (c *BaseModel[T]) Count() (int, error) {
	var count int64
	// 使用 `(*T)(nil)` 作为模型类型占位，避免分配一个实际的零值对象。
	// GORM 仅使用传入值的类型信息来确定表名/模型元信息，因此传入一个类型为 `*T` 的 nil 指针
	// 能够达到同样目的并略微减少一次内存分配。
	err := c.getQuery().Model((*T)(nil)).Count(&count).Error
	return int(count), err
}

// Exist 检查记录是否存在
// 使用gorm的Take方法检查是否存在匹配的记录
// 不会返回 gorm.ErrRecordNotFound 错误，而是将其转换为 false 返回值
//
// 示例:
//
//	// 检查特定ID的用户是否存在
//	exists, err := userModel.Where("id = ?", 1).Exist()
//	if err != nil {
//	    log.Error("Database error:", err)
//	    return
//	}
//
//	if exists {
//	    fmt.Println("User exists")
//	} else {
//	    fmt.Println("User does not exist")
//	}
//
// 边界情况:
//   - 当查询条件匹配多条记录时，只检查是否存在至少一条记录
//   - 数据库连接错误会返回错误
//
// 返回值:
//   - bool: 是否存在匹配的记录
//   - error: 数据库查询过程中可能发生的错误（不包括记录不存在）
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

// Page 分页查询记录
// 使用gorm的Offset和Limit方法实现分页查询
// 自动处理页数和页面大小的边界情况
//
// 示例:
//
//	// 查询第一页，每页10条记录
//	users, total, err := userModel.Page(1, 10)
//	if err != nil {
//	    log.Error("Failed to query users:", err)
//	    return
//	}
//
//	fmt.Printf("Total users: %d\n", total)
//	fmt.Printf("Current page users: %d\n", len(users))
//	for _, user := range users {
//	    fmt.Printf("User: %+v\n", user)
//	}
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
func (c *BaseModel[T]) Page(page, pageSize int) ([]T, int, error) {
	var items []T
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

	total, err := c.Count()
	if err != nil {
		return items, total, err
	}

	if total == 0 {
		return items, total, nil
	}

	// 计算偏移量 (页数-1)*每页记录数
	offset := (page - 1) * pageSize
	err = c.getQuery().Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
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
// 使用gorm的Transaction方法包装事务操作
// 支持事务嵌套，如果已经在事务中，则使用当前事务
//
// 示例:
//
//	// 在事务中执行多个操作
//	err := userModel.Transaction(func(tx *gorm.DB) error {
//	    // 创建用户
//	    user := &User{Name: "John"}
//	    if err := tx.Create(user).Error; err != nil {
//	        return err // 回滚事务
//	    }
//
//	    // 创建用户配置
//	    profile := &Profile{UserID: user.ID, Bio: "Hello"}
//	    if err := tx.Create(profile).Error; err != nil {
//	        return err // 回滚事务
//	    }
//
//	    return nil // 提交事务
//	})
//
//	if err != nil {
//	    log.Error("Transaction failed:", err)
//	} else {
//	    fmt.Println("Transaction succeeded")
//	}
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
func (c *BaseModel[T]) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	// 使用当前绑定的数据库连接，而不是直接使用全局连接
	// 这样可以尊重通过WithTx绑定的事务，实现正确的事务嵌套
	return c.getQuery().Transaction(fc, opts...)
}
