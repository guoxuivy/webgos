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
			department.GET("/:id", "部门详情", handlers.GetDepartment)
			department.GET("/tree", "部门树", handlers.GetDepartmentTree)
			department.GET("/:id/users", "部门用户", handlers.GetDepartmentUsers)
			department.PUT("/:id/leader", "设置负责人", handlers.SetDepartmentLeader)
		}
	})
}