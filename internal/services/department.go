package services

import (
	"errors"

	"webgos/internal/xdb"
	"webgos/internal/dto"
	"webgos/internal/models"

	"gorm.io/gorm"
)

type DepartmentService interface {
	Create(dtoModel dto.AddDepartmentDTO) (*models.Department, error)
	Update(dtoModel dto.EditDepartmentDTO) error
	Delete(id int) error
	GetByID(id int) (*models.Department, error)
	GetTree() ([]models.Department, error)
	GetUsers(departmentID int, page, pageSize int) ([]models.User, int64, error)
	SetLeader(id, leaderID int) error
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

	if dtoModel.LeaderID > 0 {
		var leader models.User
		if err := xdb.GetDB().First(&leader, dtoModel.LeaderID).Error; err != nil {
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
		department.LeaderID = *dtoModel.LeaderID
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

func (s *departmentService) GetByID(id int) (*models.Department, error) {
	var department models.Department
	if err := xdb.GetDB().Preload("Leader").First(&department, id).Error; err != nil {
		return nil, errors.New("部门不存在")
	}

	return &department, nil
}

func (s *departmentService) GetTree() ([]models.Department, error) {
	var departments []models.Department
	if err := xdb.GetDB().Preload("Leader").Order("parent_id ASC, sort ASC").Find(&departments).Error; err != nil {
		return nil, err
	}

	return buildDepartmentTree(departments), nil
}

func buildDepartmentTree(departments []models.Department) []models.Department {
	departmentMap := make(map[int]*models.Department)
	var rootDepartments []models.Department

	for i := range departments {
		departmentMap[departments[i].ID] = &departments[i]
	}

	for _, dept := range departments {
		if dept.ParentID == 0 {
			rootDepartments = append(rootDepartments, dept)
		} else {
			if parent, ok := departmentMap[dept.ParentID]; ok {
				parent.Children = append(parent.Children, dept)
			}
		}
	}

	return rootDepartments
}

func (s *departmentService) GetUsers(departmentID int, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := xdb.GetDB().Where("department_id = ?", departmentID)
	db.Count(&total)
	db = db.Scopes(models.Page(page, pageSize)).Find(&users)

	return users, total, db.Error
}

func (s *departmentService) SetLeader(id, leaderID int) error {
	var department models.Department
	if err := xdb.GetDB().First(&department, id).Error; err != nil {
		return errors.New("部门不存在")
	}

	var leader models.User
	if err := xdb.GetDB().First(&leader, leaderID).Error; err != nil {
		return errors.New("用户不存在")
	}

	return xdb.GetDB().Model(&department).Update("leader_id", leaderID).Error
}