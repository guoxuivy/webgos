package services

import (
	"errors"

	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xdb"
	"webgos/internal/xlog"

	"gorm.io/gorm"
)

type DepartmentService interface {
	Create(dtoModel dto.AddDepartmentDTO) (*models.Department, error)
	Update(dtoModel dto.EditDepartmentDTO) error
	Delete(id int) error
	GetTree() ([]models.Department, error)
	AddUsers(departmentID int, userIDs []int) error
}

type departmentService struct{}

func NewDepartmentService() DepartmentService {
	return &departmentService{}
}

func (s *departmentService) Create(dtoModel dto.AddDepartmentDTO) (*models.Department, error) {
	if dtoModel.ParentID > 0 {
		var parent models.Department
		if err := xdb.GetDB().First(&parent, dtoModel.ParentID).Error; err != nil {
			return nil, errors.New("父部门不存在")
		}
	}

	if dtoModel.LeaderID != nil && *dtoModel.LeaderID > 0 {
		var leader models.User
		if err := xdb.GetDB().First(&leader, *dtoModel.LeaderID).Error; err != nil {
			return nil, errors.New("负责人不存在")
		}
	}

	department := dtoModel.ToModel()
	if department.Status == 0 {
		department.Status = 1
	}

	if err := xdb.GetDB().Create(&department).Error; err != nil {
		return nil, err
	}

	return &department, nil
}

func (s *departmentService) Update(dtoModel dto.EditDepartmentDTO) error {
	var department models.Department
	if err := xdb.GetDB().First(&department, dtoModel.ID).Error; err != nil {
		return errors.New("部门不存在")
	}

	if dtoModel.ParentID != nil && *dtoModel.ParentID > 0 {
		var parent models.Department
		if err := xdb.GetDB().First(&parent, *dtoModel.ParentID).Error; err != nil {
			return errors.New("父部门不存在")
		}
	}

	if dtoModel.LeaderID != nil && *dtoModel.LeaderID > 0 {
		var leader models.User
		if err := xdb.GetDB().First(&leader, *dtoModel.LeaderID).Error; err != nil {
			return errors.New("负责人不存在")
		}
	}

	if dtoModel.Name != nil {
		department.Name = *dtoModel.Name
	}
	if dtoModel.ParentID != nil {
		department.ParentID = *dtoModel.ParentID
	}
	if dtoModel.LeaderID != nil {
		department.LeaderID = dtoModel.LeaderID
	}
	if dtoModel.Remark != nil {
		department.Remark = *dtoModel.Remark
	}
	if dtoModel.Status != nil {
		department.Status = *dtoModel.Status
	}
	if dtoModel.Order != nil {
		department.Sort = *dtoModel.Order
	}

	return xdb.GetDB().Select("*").Updates(&department).Error
}

func (s *departmentService) Delete(id int) error {
	var department models.Department
	if err := xdb.GetDB().First(&department, id).Error; err != nil {
		return errors.New("部门不存在")
	}

	return xdb.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&models.Department{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("department_id = ?", id).Update("department_id", 0).Error; err != nil {
			return err
		}

		return tx.Delete(&department, id).Error
	})
}

func (s *departmentService) GetTree() ([]models.Department, error) {
	var departments []models.Department
	if err := xdb.GetDB().Preload("Leader").Order("parent_id ASC, sort ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return s.buildDepartmentTree(departments, 0), nil
}

func (s *departmentService) buildDepartmentTree(departments []models.Department, parentID int) []models.Department {
	var tree []models.Department
	for i := range departments {
		if departments[i].ParentID == parentID {
			children := s.buildDepartmentTree(departments, departments[i].ID)
			departments[i].Children = children

			if err := xdb.GetDB().Where("department_id = ?", departments[i].ID).Find(&departments[i].Users).Error; err != nil {
				xlog.Error("加载部门成员失败: %v", err)
			}
			tree = append(tree, departments[i])
		}
	}
	return tree
}

func (s *departmentService) AddUsers(departmentID int, userIDs []int) error {
	var department models.Department
	if err := xdb.GetDB().First(&department, departmentID).Error; err != nil {
		return errors.New("部门不存在")
	}

	if len(userIDs) == 0 {
		return nil
	}

	return xdb.GetDB().Model(&models.User{}).Where("id IN ?", userIDs).Update("department_id", departmentID).Error
}
