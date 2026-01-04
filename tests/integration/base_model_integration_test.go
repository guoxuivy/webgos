package integration

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	"webgos/internal/database"
	"webgos/internal/models"
	"webgos/internal/xlog"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestMain 设置测试环境
func TestMain(m *testing.M) {
	// 初始化测试日志系统
	xlog.InitLogger()

	// 初始化测试数据库
	code := 0
	if setupTestDB() == nil {
		// 运行测试
		code = m.Run()

		// 清理测试数据
		teardownTestDB()
	}

	// 关闭日志系统
	xlog.Xlogger.Close()

	os.Exit(code)
}

// setupTestDB 初始化测试数据库
func setupTestDB() error {
	// 使用与主应用相同的数据库配置进行测试
	// 在实际项目中，应该使用独立的测试数据库

	// 构建DSN字符串 - 使用现有数据库进行测试
	dsn := "root:123456@tcp(localhost:3306)/hserp?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化测试数据库
	_, err := database.InitDB(dsn)
	if err != nil {
		// 如果连接失败，打印错误但不中断测试
		// 因为有些测试不需要数据库连接
		return err
	}

	return nil
}

// teardownTestDB 清理测试数据库
func teardownTestDB() {
	// 关闭数据库连接
	database.CloseDB()
}

// generateTestUsername 生成测试用户名
func generateTestUsername(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), rand.Intn(10000))
}

// TestBaseModelCRUDIntegration 测试 BaseModel 的实际 CRUD 操作
func TestBaseModelCRUDIntegration(t *testing.T) {
	// 检查数据库是否连接成功
	if database.DB == nil {
		t.Skip("数据库未连接，跳过集成测试")
	}

	t.Run("TestCreateAndRead", func(t *testing.T) {
		// 创建用户模型实例
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername := generateTestUsername("testuser_integration")

		// 创建测试用户
		testUser := &models.User{
			Username: testUsername,
			Nickname: "Test User",
			Email:    "test@example.com",
			Phone:    "13800138000",
		}
		testUser.SetPassword("password123")

		// 测试创建操作
		err := userModel.Create(testUser)
		assert.NoError(t, err)
		assert.True(t, testUser.ID > 0)

		// 测试读取操作
		readUser, err := userModel.Read(testUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, readUser)
		assert.Equal(t, testUser.Username, readUser.Username)
		assert.Equal(t, testUser.Nickname, readUser.Nickname)

		// 清理测试数据
		err = userModel.Delete(testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("TestUpdate", func(t *testing.T) {
		// 创建用户模型实例
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername := generateTestUsername("testuser_update")

		// 创建测试用户
		testUser := &models.User{
			Username: testUsername,
			Nickname: "Test User",
			Email:    "test@example.com",
		}
		testUser.SetPassword("password123")

		// 创建用户
		err := userModel.Create(testUser)
		assert.NoError(t, err)

		// 更新用户信息
		testUser.Nickname = "Updated User"
		testUser.Email = "updated@example.com"
		err = userModel.Update(testUser)
		assert.NoError(t, err)

		// 验证更新
		updatedUser, err := userModel.Read(testUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated User", updatedUser.Nickname)
		assert.Equal(t, "updated@example.com", updatedUser.Email)

		// 清理测试数据
		err = userModel.Delete(testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("TestDelete", func(t *testing.T) {
		// 创建用户模型实例
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername := generateTestUsername("testuser_delete")

		// 创建测试用户
		testUser := &models.User{
			Username: testUsername,
			Nickname: "Test User",
			Email:    "test@example.com",
		}
		testUser.SetPassword("password123")

		// 创建用户
		err := userModel.Create(testUser)
		assert.NoError(t, err)

		// 删除用户
		err = userModel.Delete(testUser.ID)
		assert.NoError(t, err)

		// 验证用户已被删除（软删除）
		_, err = userModel.Read(testUser.ID)
		assert.Error(t, err)
	})
}

// TestBaseModelQueryIntegration 测试 BaseModel 的查询操作
func TestBaseModelQueryIntegration(t *testing.T) {
	// 检查数据库是否连接成功
	if database.DB == nil {
		t.Skip("数据库未连接，跳过集成测试")
	}

	t.Run("TestWhereAndMore", func(t *testing.T) {
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername1 := generateTestUsername("testuser_query1")
		testUsername2 := generateTestUsername("testuser_query2")

		// 创建测试用户
		testUser1 := &models.User{
			Username: testUsername1,
			Nickname: "Query Test 1",
			Email:    "query1@example.com",
			Age:      25,
		}
		testUser1.SetPassword("password123")

		testUser2 := &models.User{
			Username: testUsername2,
			Nickname: "Query Test 2",
			Email:    "query2@example.com",
			Age:      30,
		}
		testUser2.SetPassword("password123")

		// 创建用户
		err := userModel.Create(testUser1)
		assert.NoError(t, err)
		err = userModel.Create(testUser2)
		assert.NoError(t, err)

		// 测试 Where 查询
		users, err := userModel.Where("age > ?", 20).More()
		assert.NoError(t, err)
		assert.NotEmpty(t, users)

		// 测试 Count
		count, err := userModel.Where("age > ?", 20).Count()
		assert.NoError(t, err)
		assert.True(t, count > 0)

		// 测试 Order 和 Limit
		users, err = userModel.Where("username LIKE ?", "testuser_query%").
			Order("age DESC").
			Limit(2).
			More()
		assert.NoError(t, err)
		xlog.Info("Queried Users: %v", users)
		// 注意：这里可能无法准确获取到我们创建的用户，因为LIKE匹配可能包含其他测试数据

		// 清理测试数据
		userModel.Delete(testUser1.ID)
		userModel.Delete(testUser2.ID)
	})

	t.Run("TestPage", func(t *testing.T) {
		userModel := &models.User{}

		// 生成唯一的测试用户名前缀
		usernamePrefix := generateTestUsername("testuser_page")

		// 创建多个测试用户
		var createdUsers []models.User
		for i := 0; i < 5; i++ {
			testUser := &models.User{
				Username: fmt.Sprintf("%s_%d", usernamePrefix, i),
				Nickname: fmt.Sprintf("Page Test %d", i),
				Email:    fmt.Sprintf("page%d@example.com", i),
			}
			testUser.SetPassword("password123")
			err := userModel.Create(testUser)
			assert.NoError(t, err)
			createdUsers = append(createdUsers, *testUser)
		}

		// 使用特定条件进行分页查询，确保只查询我们创建的用户
		users, total, err := userModel.Where("username LIKE ?", fmt.Sprintf("%s%%", usernamePrefix)).Page(1, 3)
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, users, 3)

		users, total, err = userModel.Where("username LIKE ?", fmt.Sprintf("%s%%", usernamePrefix)).Page(2, 3)
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, users, 2)

		// 清理测试数据
		for _, user := range createdUsers {
			userModel.Delete(user.ID)
		}
	})
}

// TestBaseModelChainableIntegration 测试 BaseModel 的链式查询方法
func TestBaseModelChainableIntegration(t *testing.T) {
	// 检查数据库是否连接成功
	if database.DB == nil {
		t.Skip("数据库未连接，跳过集成测试")
	}

	t.Run("TestChainableMethodsWork", func(t *testing.T) {
		userModel := &models.User{}

		// 测试链式调用是否正常工作
		query := userModel.Where("age > ?", 18).
			Order("created_at DESC").
			Limit(10)

		// 验证返回的是正确的接口类型
		assert.NotNil(t, query)

		// 执行查询
		users, err := query.More()
		assert.NoError(t, err)
		assert.NotNil(t, users)
	})
}

// TestBaseModelTransactionIntegration 测试 BaseModel 的事务操作
func TestBaseModelTransactionIntegration(t *testing.T) {
	// 检查数据库是否连接成功
	if database.DB == nil {
		t.Skip("数据库未连接，跳过集成测试")
	}

	t.Run("TestTransactionSuccess", func(t *testing.T) {
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername1 := generateTestUsername("testuser_tx1")
		testUsername2 := generateTestUsername("testuser_tx2")

		// 使用事务创建两个用户
		err := userModel.Transaction(func(tx *gorm.DB) error {
			// 创建第一个用户
			user1 := &models.User{
				Username: testUsername1,
				Nickname: "Transaction User 1",
				Email:    "tx1@example.com",
			}
			user1.SetPassword("password123")
			if err := tx.Create(user1).Error; err != nil {
				return err
			}

			// 创建第二个用户
			user2 := &models.User{
				Username: testUsername2,
				Nickname: "Transaction User 2",
				Email:    "tx2@example.com",
			}
			user2.SetPassword("password123")
			if err := tx.Create(user2).Error; err != nil {
				return err
			}

			return nil
		})

		assert.NoError(t, err)

		// 验证用户是否创建成功
		users, err := userModel.Where("username IN ?", []string{testUsername1, testUsername2}).More()
		assert.NoError(t, err)
		assert.Len(t, users, 2)

		// 清理测试数据
		for _, user := range users {
			userModel.Delete(user.ID)
		}
	})

	t.Run("TestTransactionRollback", func(t *testing.T) {
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername1 := generateTestUsername("testuser_tx3")
		testUsername2 := generateTestUsername("testuser_tx4") // 这个会重复导致错误

		// 创建一个已存在的用户来制造冲突
		existingUser := &models.User{
			Username: testUsername2,
			Nickname: "Existing User",
			Email:    "existing@example.com",
		}
		existingUser.SetPassword("password123")
		err := userModel.Create(existingUser)
		assert.NoError(t, err)

		// 使用事务创建用户，但会因为重复用户名而失败
		err = userModel.Transaction(func(tx *gorm.DB) error {
			// 创建第一个用户
			user1 := &models.User{
				Username: testUsername1,
				Nickname: "Transaction User 3",
				Email:    "tx3@example.com",
			}
			user1.SetPassword("password123")
			if createErr := tx.Create(user1).Error; createErr != nil {
				return createErr
			}

			// 尝试创建第二个用户，会因为用户名重复而失败
			user2 := &models.User{
				Username: testUsername2, // 重复的用户名
				Nickname: "Transaction User 4",
				Email:    "tx4@example.com",
			}
			user2.SetPassword("password123")
			if createErr := tx.Create(user2).Error; createErr != nil {
				return createErr
			}

			return nil
		})

		// 事务应该回滚，返回错误
		assert.Error(t, err)

		// 验证第一个用户没有被创建（事务回滚）
		_, err = userModel.Where("username = ?", testUsername1).One()
		assert.Error(t, err)

		// 清理测试数据
		userModel.Delete(existingUser.ID)
	})

	t.Run("TestWithTxMethod", func(t *testing.T) {
		userModel := &models.User{}

		// 生成唯一的测试用户名
		testUsername := generateTestUsername("testuser_with_tx")

		// 使用事务创建用户
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			// 使用 WithTx 方法创建绑定到事务的模型实例
			txUserModel := userModel.WithTx(tx)

			// 创建用户
			user := &models.User{
				Username: testUsername,
				Nickname: "WithTx User",
				Email:    "withtx@example.com",
			}
			user.SetPassword("password123")

			// 使用绑定到事务的模型实例创建用户
			if err := txUserModel.Create(user); err != nil {
				return err
			}

			return nil
		})

		assert.NoError(t, err)

		// 验证用户是否创建成功
		createdUser, err := userModel.Where("username = ?", testUsername).One()
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)

		// 清理测试数据
		userModel.Delete(createdUser.ID)
	})

	// 事务嵌套功能测试
	t.Run("TestTransactionNestedFunctionality", func(t *testing.T) {
		// 创建用户模型实例
		userModel := &models.BaseModel[models.User]{}

		// 清理测试数据
		defer database.DB.Where("username LIKE ?", "test_transaction_%").Delete(&models.User{})

		t.Run("TestTransactionRespectsBoundTx", func(t *testing.T) {
			t.Run("InnerFailureDoesNotAffectOuter", func(t *testing.T) {
				// 测试场景1: 内层事务失败，外层事务可以选择继续（忽略内层错误）
				outerErr := database.DB.Transaction(func(tx *gorm.DB) error {
					// 创建测试用户
					testUser := &models.User{
						Username: "test_transaction_outer_continue",
						Password: "password123",
						Email:    "outer_continue@test.com",
					}

					// 使用事务创建用户
					if createErr := tx.Create(testUser).Error; createErr != nil {
						return createErr
					}

					// 使用 WithTx 绑定到外层事务
					txUserModel := userModel.WithTx(tx)

					// 内层使用 Transaction 方法
					innerErr := txUserModel.Transaction(func(innerTx *gorm.DB) error {
						// 创建内层测试用户
						innerUser := &models.User{
							Username: "test_transaction_inner_fail",
							Password: "password123",
							Email:    "inner_fail@test.com",
						}

						// 使用内层事务创建用户
						if createErr := innerTx.Create(innerUser).Error; createErr != nil {
							return createErr
						}

						// 故意返回错误，触发内层事务回滚
						return assert.AnError
					})

					assert.Error(t, innerErr)
					return nil // 关键：忽略内层错误，外层继续执行
				})

				assert.NoError(t, outerErr)

				// 检查数据库状态
				var outerUserCount int64
				var innerUserCount int64
				database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_outer_continue").Count(&outerUserCount)
				database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_inner_fail").Count(&innerUserCount)

				// 验证结果
				assert.Equal(t, int64(1), outerUserCount, "外层事务的用户应该被成功创建")
				assert.Equal(t, int64(0), innerUserCount, "内层事务的用户应该被回滚")
			})

			t.Run("InnerFailurePropagatesToOuter", func(t *testing.T) {
				// 测试场景2: 内层事务失败，外层事务也应该回滚（传播错误）
				outerErr := database.DB.Transaction(func(tx *gorm.DB) error {
					// 创建测试用户
					testUser := &models.User{
						Username: "test_transaction_outer_rollback",
						Password: "password123",
						Email:    "outer_rollback@test.com",
					}

					// 使用事务创建用户
					if createErr := tx.Create(testUser).Error; createErr != nil {
						return createErr
					}

					// 使用 WithTx 绑定到外层事务
					txUserModel := userModel.WithTx(tx)

					// 内层使用 Transaction 方法，并且将错误传播到外层
					return txUserModel.Transaction(func(innerTx *gorm.DB) error {
						// 创建内层测试用户
						innerUser := &models.User{
							Username: "test_transaction_inner_propagate",
							Password: "password123",
							Email:    "inner_propagate@test.com",
						}

						// 使用内层事务创建用户
						if createErr := innerTx.Create(innerUser).Error; createErr != nil {
							return createErr
						}

						// 故意返回错误，触发内层事务回滚，并且传播到外层
						return assert.AnError
					})
				})

				// 外层事务应该失败
				assert.Error(t, outerErr)

				// 检查数据库状态 - 两个用户都应该被回滚
				var outerUserCount int64
				var innerUserCount int64
				database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_outer_rollback").Count(&outerUserCount)
				database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_inner_propagate").Count(&innerUserCount)

				// 验证结果
				assert.Equal(t, int64(0), outerUserCount, "外层事务的用户应该被回滚")
				assert.Equal(t, int64(0), innerUserCount, "内层事务的用户应该被回滚")
			})
		})

		t.Run("TestTransactionWithoutBoundTx", func(t *testing.T) {
			// 不绑定事务，直接使用Transaction方法
			err := userModel.Transaction(func(tx *gorm.DB) error {
				// 创建测试用户
				testUser := &models.User{
					Username: "test_transaction_direct",
					Password: "password123",
					Email:    "direct@test.com",
				}

				// 使用事务创建用户
				return tx.Create(testUser).Error
			})

			assert.NoError(t, err)

			// 验证用户创建成功
			var userCount int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_direct").Count(&userCount)
			assert.Equal(t, int64(1), userCount, "直接使用Transaction方法应该成功创建用户")
		})

		t.Run("TestNestedTransactionProperRollback", func(t *testing.T) {
			// 测试多层嵌套事务的回滚
			outerErr := database.DB.Transaction(func(outerTx *gorm.DB) error {
				// 外层创建用户1
				user1 := &models.User{
					Username: "test_transaction_outer1",
					Password: "password123",
					Email:    "outer1@test.com",
				}
				if err := outerTx.Create(user1).Error; err != nil {
					return err
				}

				// 第一层嵌套
				middleErr := userModel.WithTx(outerTx).Transaction(func(middleTx *gorm.DB) error {
					// 中层创建用户2
					user2 := &models.User{
						Username: "test_transaction_middle2",
						Password: "password123",
						Email:    "middle2@test.com",
					}
					if err := middleTx.Create(user2).Error; err != nil {
						return err
					}

					// 第二层嵌套
					return userModel.WithTx(middleTx).Transaction(func(innerTx *gorm.DB) error {
						// 内层创建用户3
						user3 := &models.User{
							Username: "test_transaction_inner3",
							Password: "password123",
							Email:    "inner3@test.com",
						}
						if err := innerTx.Create(user3).Error; err != nil {
							return err
						}

						// 内层故意失败，应该回滚所有嵌套事务
						return assert.AnError
					})
				})

				return middleErr
			})

			// 外层事务应该失败（因为内层失败传播）
			assert.Error(t, outerErr)

			// 验证所有用户都被回滚
			var user1Count int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_outer1").Count(&user1Count)

			var user2Count int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_middle2").Count(&user2Count)

			var user3Count int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_inner3").Count(&user3Count)

			assert.Equal(t, int64(0), user1Count, "外层用户应该被回滚")
			assert.Equal(t, int64(0), user2Count, "中层用户应该被回滚")
			assert.Equal(t, int64(0), user3Count, "内层用户应该被回滚")
		})

		t.Run("TestNestedTransactionPartialRollback", func(t *testing.T) {
			// 测试部分嵌套事务回滚
			outerErr := database.DB.Transaction(func(outerTx *gorm.DB) error {
				// 外层创建用户1
				user1 := &models.User{
					Username: "test_transaction_outer4",
					Password: "password123",
					Email:    "outer4@test.com",
				}
				if err := outerTx.Create(user1).Error; err != nil {
					return err
				}

				// 第一层嵌套 - 应该成功
				middleErr := userModel.WithTx(outerTx).Transaction(func(middleTx *gorm.DB) error {
					// 中层创建用户2
					user2 := &models.User{
						Username: "test_transaction_middle5",
						Password: "password123",
						Email:    "middle5@test.com",
					}
					return middleTx.Create(user2).Error
				})

				if middleErr != nil {
					return middleErr
				}

				// 第二个独立的第一层嵌套 - 应该失败并只回滚自己
				secondMiddleErr := userModel.WithTx(outerTx).Transaction(func(secondMiddleTx *gorm.DB) error {
					// 创建冲突的用户
					conflictUser := &models.User{
						Username: "test_transaction_outer4", // 与外层用户冲突
						Password: "password123",
						Email:    "conflict@test.com",
					}
					return secondMiddleTx.Create(conflictUser).Error
				})

				// 第二个中层事务应该失败，但不影响外层和第一个中层
				assert.Error(t, secondMiddleErr)
				return nil // 外层事务应该成功
			})

			// 外层事务应该成功
			assert.NoError(t, outerErr)

			// 验证结果
			var user1Count int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_outer4").Count(&user1Count)
			assert.Equal(t, int64(1), user1Count, "外层用户应该成功创建")

			var user2Count int64
			database.DB.Model(&models.User{}).Where("username = ?", "test_transaction_middle5").Count(&user2Count)
			assert.Equal(t, int64(1), user2Count, "第一个中层用户应该成功创建")

			var conflictCount int64
			database.DB.Model(&models.User{}).Where("email = ?", "conflict@test.com").Count(&conflictCount)
			assert.Equal(t, int64(0), conflictCount, "冲突的用户应该被回滚")
		})
	})
}
