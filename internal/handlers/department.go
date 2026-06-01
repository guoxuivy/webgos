package handlers

import (
	"strconv"

	"webgos/internal/dto"
	"webgos/internal/services"
	"webgos/internal/utils/param"
	"webgos/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// CreateDepartment 创建部门
// @Summary 创建部门
// @Description 创建新部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param body body dto.AddDepartmentDTO true "部门信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department [post]
// @Security BearerAuth
func CreateDepartment(c *gin.Context) {
	var dtoModel dto.AddDepartmentDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	departmentService := services.NewDepartmentService()
	department, err := departmentService.Create(dtoModel)
	if err != nil {
		response.Error(c, "创建部门失败: "+err.Error())
		return
	}

	response.Success(c, "部门创建成功", department)
}

// UpdateDepartment 更新部门
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param body body dto.EditDepartmentDTO true "部门信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department [put]
// @Security BearerAuth
func UpdateDepartment(c *gin.Context) {
	var dtoModel dto.EditDepartmentDTO

	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	departmentService := services.NewDepartmentService()
	err := departmentService.Update(dtoModel)
	if err != nil {
		response.Error(c, "更新部门失败: "+err.Error())
		return
	}

	response.Success(c, "部门更新成功", nil)
}

// DeleteDepartment 删除部门
// @Summary 删除部门
// @Description 删除指定部门（级联删除子部门）
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department/{id} [delete]
// @Security BearerAuth
func DeleteDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(c, "无效的部门ID")
		return
	}

	departmentService := services.NewDepartmentService()
	err = departmentService.Delete(id)
	if err != nil {
		response.Error(c, "删除部门失败: "+err.Error())
		return
	}

	response.Success(c, "部门删除成功", nil)
}

// GetDepartment 获取部门详情
// @Summary 获取部门详情
// @Description 根据ID获取部门详情
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department/{id} [get]
// @Security BearerAuth
func GetDepartment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(c, "无效的部门ID")
		return
	}

	departmentService := services.NewDepartmentService()
	department, err := departmentService.GetByID(id)
	if err != nil {
		response.Error(c, "获取部门失败: "+err.Error())
		return
	}

	response.Success(c, "获取部门成功", department)
}

// GetDepartmentTree 获取部门树
// @Summary 获取部门树形结构
// @Description 获取所有部门的树形结构
// @Tags 部门管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department/tree [get]
// @Security BearerAuth
func GetDepartmentTree(c *gin.Context) {
	departmentService := services.NewDepartmentService()
	tree, err := departmentService.GetTree()
	if err != nil {
		response.Error(c, "获取部门树失败: "+err.Error())
		return
	}

	response.Success(c, "获取部门树成功", tree)
}

// GetDepartmentUsers 获取部门用户
// @Summary 获取部门用户列表
// @Description 获取指定部门下的用户列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department/{id}/users [get]
// @Security BearerAuth
func GetDepartmentUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(c, "无效的部门ID")
		return
	}

	var query dto.DepartmentUserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, err.Error())
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	departmentService := services.NewDepartmentService()
	users, total, err := departmentService.GetUsers(id, query.Page, query.PageSize)
	if err != nil {
		response.Error(c, "获取部门用户失败: "+err.Error())
		return
	}

	response.Success(c, "获取部门用户成功", gin.H{"items": users, "total": total})
}

// SetDepartmentLeader 设置部门负责人
// @Summary 设置部门负责人
// @Description 设置指定部门的负责人
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Param body body dto.SetLeaderDTO true "负责人信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/department/{id}/leader [put]
// @Security BearerAuth
func SetDepartmentLeader(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Error(c, "无效的部门ID")
		return
	}

	var dtoModel dto.SetLeaderDTO
	if err := param.Validate(c, &dtoModel); err != nil {
		response.Error(c, err.Error())
		return
	}

	departmentService := services.NewDepartmentService()
	err = departmentService.SetLeader(id, dtoModel.LeaderID)
	if err != nil {
		response.Error(c, "设置负责人失败: "+err.Error())
		return
	}

	response.Success(c, "设置负责人成功", nil)
}