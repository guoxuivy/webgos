package handlers

import (
	"strconv"
	"webgos/internal/database"
	"webgos/internal/models"
	"webgos/internal/utils/response"
	"webgos/internal/xlog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary 测试接口
// @Description 测试用接口
// @Tags test
// @Produce json
// @Success 200 {object} response.Response
// @Router /test/Test [get]
func Test(c *gin.Context) {
	mod := models.User{}
	// tx := database.DB.Begin()
	// mod.Delete(2)
	// mod.Where("age > ?", 20).Where("age < ?", 90).One()
	usres, err := mod.Where("age > ?", 20).Where("age < ?", 90).Where("id < ?", 10).More()

	// database.DB.Commit()
	if err != nil {
		response.Error(c, "获取用户列表失败")
		return
	}
	// 这里可以添加更多的业务逻辑
	response.Success(c, "Test", usres)
}

// @Summary 测试事务
// @Description 测试数据库事务功能
// @Tags test
// @Produce json
// @Success 200 {object} response.Response
// @Router /test/TestTransaction [get]
func TestTransaction(c *gin.Context) {
	// 创建用户模型实例
	user := models.User{}

	// 使用事务创建两个用户
	err := user.Transaction(func(tx *gorm.DB) error {
		// 在事务中创建第一个用户
		u1 := models.User{
			Username: "transaction_user5",
			Nickname: "Transaction User 5",
			Age:      25,
		}
		err := user.WithTx(tx).Create(&u1)
		if err != nil {
			return err
		}
		// 使用tx 可以查询到刚刚插入的数据，因为在同一个事务中，不使用tx是查询不到的
		res, _ := user.WithTx(tx).Where("username = ?", "transaction_user5").Count()
		xlog.Debug("Count is :%v", res)

		// 在事务中创建第二个用户 由于Username冲突，事务会回滚
		u2 := models.User{
			Username: "transaction_user4",
			Nickname: "Transaction User 4",
			Age:      30,
		}
		err = user.WithTx(tx).Create(&u2)
		if err != nil {
			return err
		}

		// 模拟可能的错误情况，可以取消下面的注释来测试回滚
		// return errors.New("模拟事务回滚")

		return nil
	})

	if err != nil {
		response.Error(c, "事务执行失败: "+err.Error())
		return
	}

	response.Success(c, "事务执行成功，用户创建完成", nil)
}

// @Summary 测试事务2
// @Description 测试数据库事务功能
// @Tags test
// @Produce json
// @Success 200 {object} response.Response
// @Router /test/TestTransaction2 [get]
func TestTransaction2(c *gin.Context) {

	tx := database.DB.Begin()

	// 创建用户模型实例
	user := models.User{}

	// 在事务中创建第一个用户
	u1 := models.User{
		Username: "transaction_user5",
		Nickname: "Transaction User 5",
		Age:      25,
	}
	err := user.WithTx(tx).Create(&u1)
	if err != nil {
		tx.Rollback()
		response.Error(c, "事务执行失败: "+err.Error())
	}
	// 使用tx 可以查询到刚刚插入的数据，因为在同一个事务中，不使用tx是查询不到的
	res, _ := user.WithTx(tx).Where("username = ?", "transaction_user5").Count()
	xlog.Debug("Count is :%v", res)

	// 在事务中创建第二个用户 由于Username冲突，事务会回滚
	u2 := models.User{
		Username: "transaction_user4",
		Nickname: "Transaction User 4",
		Age:      30,
	}
	err = user.WithTx(tx).Create(&u2)
	if err != nil {
		tx.Rollback() // 如果没有使用defer，记得每次出错都要手动回滚，不然上面的语句会被执行成功
		response.Error(c, "事务执行失败: "+err.Error())
	}
	tx.Commit()

	response.Success(c, "事务执行成功，用户创建完成", nil)
}

// @Summary TestCurlList
// @Description 测试用接口
// @Tags test
// @Produce json
// @Success 200 {object} response.Response
// @Router /test/TestCurlList [get]
func TestCurlList(c *gin.Context) {
	mod := models.User{}
	usres, err := mod.Where("1=1").One()
	if err != nil {
		response.Error(c, "获取用户列表失败")
		return
	}
	// 这里可以添加更多的业务逻辑
	response.Success(c, "Test", usres)
}

// @Summary 更新用户年龄
// @Description 根据用户ID更新用户年龄的测试接口
// @Tags test
// @Produce json
// @Param id path int true "用户ID"
// @Param age path int true "用户年龄"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /test/TestCurlUpdate/{id}/{age} [get]
func TestCurlUpdate(c *gin.Context) {
	id, _ := c.Params.Get("id")
	age, _ := c.Params.Get("age")

	// 将id从string转换为int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, "参数错误")
		return
	}

	// 将age从string转换为int
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		response.Error(c, "参数错误")
		return
	}

	mod := models.User{}
	user, err := mod.Read(idInt)
	if err != nil {
		response.Error(c, "用户不存在")
		return
	}

	user.Age = ageInt
	mod.Update(user)

	// 这里可以添加更多的业务逻辑
	response.Success(c, "Test", user)
}
