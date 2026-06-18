package routes

import (
	"webgos/internal/handlers"

	"github.com/gin-gonic/gin"
)

func init() {
	Register(func(router *gin.Engine) {
		api := router.Group("/api")

		department := WrapRouter(api.Group("/department"))
		{
			department.POST("", "创建部门", handlers.CreateDepartment)
			department.PUT("", "更新部门", handlers.UpdateDepartment)
			department.DELETE("/:id", "删除部门", handlers.DeleteDepartment)
			department.GET("/tree", "部门树", handlers.GetDepartmentTree)
			department.POST("/:id/users", "批量添加用户", handlers.AddDepartmentUsers)
			department.DELETE("/user/:userID", "移除部门用户", handlers.RemoveDepartmentUser)
		}
	})
}
