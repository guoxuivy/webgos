package dto

import "webgos/internal/models"

type AddDepartmentDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=50" label:"部门名称"`
	ParentID int    `json:"parent_id" validate:"omitempty,gte=0" label:"父部门ID"`
	LeaderID *int   `json:"leader_id" validate:"omitempty,gte=1" label:"负责人ID"`
	Remark   string `json:"remark" validate:"omitempty,max=200" label:"备注"`
	Status   int    `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
	Order    int    `json:"order" validate:"omitempty,gte=0" label:"排序"`
}

func (dto *AddDepartmentDTO) ToModel() models.Department {
	return models.Department{
		Name:     dto.Name,
		ParentID: dto.ParentID,
		LeaderID: dto.LeaderID,
		Remark:   dto.Remark,
		Status:   dto.Status,
		Sort:     dto.Order,
	}
}

type EditDepartmentDTO struct {
	ID       int     `json:"id" validate:"required,gte=1" label:"部门ID"`
	Name     *string `json:"name" validate:"omitempty,min=2,max=50" label:"部门名称"`
	ParentID *int    `json:"parent_id" validate:"omitempty,gte=0" label:"父部门ID"`
	LeaderID *int    `json:"leader_id" validate:"omitempty,gte=0" label:"负责人ID"`
	Remark   *string `json:"remark" validate:"omitempty,max=200" label:"备注"`
	Status   *int    `json:"status" validate:"omitempty,oneof=0 1" label:"状态"`
	Order    *int    `json:"order" validate:"omitempty,gte=0" label:"排序"`
}

type BatchUpdateDeptUsersDTO struct {
	UserIDs []int `json:"user_ids" validate:"required" label:"用户ID列表"`
}
