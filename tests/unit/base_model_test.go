package unit

import (
	"hserp/internal/models"
	"hserp/internal/xlog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 初始化测试日志系统
	xlog.InitLogger("./logs", true)

	// 初始化测试数据库
	code := m.Run()
	// 关闭日志系统
	xlog.Xlogger.Close()

	os.Exit(code) //确保所有测试任务完成后再退出
}

// TestBaseModelCRUD 测试 BaseModel 的基本 CRUD 操作
func TestBaseModelCRUD(t *testing.T) {
	// 由于测试环境没有数据库连接，这里只是验证方法是否存在
	// 实际的数据库测试需要在集成测试中进行

	t.Run("TestBaseModelStruct", func(t *testing.T) {
		// 测试 BaseModel 结构体是否正确初始化
		userModel := &models.BaseModel[models.User]{}
		assert.NotNil(t, userModel)
	})

	t.Run("TestWithTxMethod", func(t *testing.T) {
		// 测试 WithTx 方法是否存在
		userModel := &models.BaseModel[models.User]{}
		newModel := userModel.WithTx(nil)
		assert.NotNil(t, newModel)
	})

	t.Run("TestGetMethod", func(t *testing.T) {
		// 测试 getQuery 方法是否存在（私有方法，通过公共方法间接测试）
		userModel := &models.BaseModel[models.User]{}
		assert.NotNil(t, userModel)
	})

	t.Run("TestCRUDMethods", func(t *testing.T) {
		// 测试 CRUD 方法是否存在
		userModel := &models.BaseModel[models.User]{}

		// 验证所有 CRUD 方法都存在
		assert.NotNil(t, userModel.Create)
		assert.NotNil(t, userModel.Read)
		assert.NotNil(t, userModel.Update)
		assert.NotNil(t, userModel.Delete)
		assert.NotNil(t, userModel.More)
		assert.NotNil(t, userModel.One)
		assert.NotNil(t, userModel.Count)
		assert.NotNil(t, userModel.Page)
	})
}

// TestBaseModelChainableMethods 测试 BaseModel 的链式查询方法
func TestBaseModelChainableMethods(t *testing.T) {
	t.Run("TestChainableMethodsExist", func(t *testing.T) {
		userModel := &models.BaseModel[models.User]{}

		// 验证所有链式查询方法都存在
		assert.NotNil(t, userModel.Where)
		assert.NotNil(t, userModel.Select)
		assert.NotNil(t, userModel.Order)
		assert.NotNil(t, userModel.Not)
		assert.NotNil(t, userModel.Or)
		assert.NotNil(t, userModel.Limit)
		assert.NotNil(t, userModel.Group)
		assert.NotNil(t, userModel.Having)
		assert.NotNil(t, userModel.Joins)
		assert.NotNil(t, userModel.InnerJoins)
	})

	t.Run("TestChainableInterface", func(t *testing.T) {
		// 验证 BaseModel 实现了 IActiveRecode 接口
		var _ models.IActiveRecode[models.User] = &models.BaseModel[models.User]{}
	})
}

// TestBaseModelTransaction 测试 BaseModel 的事务方法
func TestBaseModelTransaction(t *testing.T) {
	t.Run("TestTransactionMethod", func(t *testing.T) {
		userModel := &models.BaseModel[models.User]{}
		assert.NotNil(t, userModel.Transaction)
	})
}

// TestUserModel 测试 User 模型是否正确嵌入 BaseModel
func TestUserModel(t *testing.T) {
	t.Run("TestUserModelStructure", func(t *testing.T) {
		user := &models.User{}

		// 验证 User 模型包含了 BaseModel 的所有方法
		assert.NotNil(t, user.Create)
		assert.NotNil(t, user.Read)
		assert.NotNil(t, user.Update)
		assert.NotNil(t, user.Delete)
		assert.NotNil(t, user.More)
		assert.NotNil(t, user.One)
		assert.NotNil(t, user.Count)
		assert.NotNil(t, user.Page)
		assert.NotNil(t, user.Where)
		assert.NotNil(t, user.Select)
		assert.NotNil(t, user.Order)
	})
}
func TestTT(t *testing.T) {
	// 初始化测试日志系统
	t.Run("TestTT", func(t *testing.T) {

		xlog.Debug("0: %v", ss(0))
		xlog.Debug("1: %v", ss(1))
		xlog.Debug("2: %v", ss(2)) //-2.2
		xlog.Debug("3: %v", ss(3)) //-2.2
		xlog.Debug("4: %v", ss(4)) //-2.2

	})

}

// 计算收益
// 本金100
// num 为上涨次数
// 每次上涨投入20，下跌10%时止损
func ss(num int) float64 {
	const rate = 0.1   // 上涨加仓幅度
	const drate = 0.07 // 下跌清仓幅度
	const invest = 20  // 每次投入资金
	// const base = 100   // 本金 仓位

	// 根据上涨次数计算总投入
	totalInvest := float64((num + 1) * invest)

	switch num {
	case 0: // 投入资金20
		currentGain := float64(invest)
		return currentGain*(1-drate) - totalInvest
	case 1: // 投入资金40
		gain := invest + invest*rate + invest
		return gain*(1-drate) - totalInvest
	case 2: // 投入资金60
		prevGain := invest + invest*rate + invest
		currentGain := prevGain + prevGain*rate + invest
		return currentGain*(1-drate) - totalInvest
	case 3: // 投入资金80
		prevGain1 := invest + invest*rate + invest
		prevGain2 := prevGain1 + prevGain1*rate + invest
		currentGain := prevGain2 + prevGain2*rate + invest
		return currentGain*(1-drate) - totalInvest
	case 4: // 投入资金100
		prevGain1 := invest + invest*rate + invest
		prevGain2 := prevGain1 + prevGain1*rate + invest
		prevGain3 := prevGain2 + prevGain2*rate + invest
		currentGain := prevGain3 + prevGain3*rate + invest
		return currentGain*(1-drate) - totalInvest
	default:
		return 0
	}
}
