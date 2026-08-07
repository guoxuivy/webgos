package services

import (
	"context"
	"errors"

	"webgos/internal/dto"
	"webgos/internal/models"
	"webgos/internal/xlog"

	"gorm.io/gorm"
)

type DepartmentService interface {
	Create(ctx context.Context, dtoModel dto.AddDepartmentDTO) (*models.Department, error)
	Update(ctx context.Context, dtoModel dto.EditDepartmentDTO) error
	Delete(ctx context.Context, id int) error
	GetTree(ctx context.Context) ([]models.Department, error)
	AddUsers(ctx context.Context, departmentID int, userIDs []int) error
	RemoveUser(ctx context.Context, userID int) error
}

type departmentService struct{}

func NewDepartmentService() DepartmentService {
	return &departmentService{}
}

func (s *departmentService) Create(ctx context.Context, dtoModel dto.AddDepartmentDTO) (*models.Department, error) {
	if dtoModel.ParentID > 0 {
		var parent models.Department
		if err := ctxDB(ctx).First(&parent, dtoModel.ParentID).Error; err != nil {
			return nil, errors.New("父部门不存在")
		}
	}

	if dtoModel.LeaderID != nil && *dtoModel.LeaderID > 0 {
		var leader models.User
		if err := ctxDB(ctx).First(&leader, *dtoModel.LeaderID).Error; err != nil {
			return nil, errors.New("负责人不存在")
		}
	}

	department := dtoModel.ToModel()
	if department.Status == 0 {
		department.Status = 1
	}

	if err := ctxDB(ctx).Create(&department).Error; err != nil {
		return nil, err
	}

	return &department, nil
}

func (s *departmentService) Update(ctx context.Context, dtoModel dto.EditDepartmentDTO) error {
	var department models.Department
	if err := ctxDB(ctx).First(&department, dtoModel.ID).Error; err != nil {
		return errors.New("部门不存在")
	}

	if dtoModel.ParentID != nil && *dtoModel.ParentID > 0 {
		var parent models.Department
		if err := ctxDB(ctx).First(&parent, *dtoModel.ParentID).Error; err != nil {
			return errors.New("父部门不存在")
		}
	}

	if dtoModel.LeaderID != nil && *dtoModel.LeaderID > 0 {
		var leader models.User
		if err := ctxDB(ctx).First(&leader, *dtoModel.LeaderID).Error; err != nil {
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
	if dtoModel.Sort != nil {
		department.Sort = *dtoModel.Sort
	}

	return ctxDB(ctx).Select("*").Updates(&department).Error
}

func (s *departmentService) Delete(ctx context.Context, id int) error {
	var department models.Department
	if err := ctxDB(ctx).First(&department, id).Error; err != nil {
		return errors.New("部门不存在")
	}

	return ctxDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("parent_id = ?", id).Delete(&models.Department{}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("department_id = ?", id).Update("department_id", 0).Error; err != nil {
			return err
		}

		return tx.Delete(&department, id).Error
	})
}

func (s *departmentService) GetTree(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	if err := ctxDB(ctx).Preload("Leader").Order("parent_id ASC, sort ASC").Find(&departments).Error; err != nil {
		return nil, err
	}
	return s.buildDepartmentTree(ctx, departments, 0), nil
}

func (s *departmentService) buildDepartmentTree(ctx context.Context, departments []models.Department, parentID int) []models.Department {
	var tree []models.Department
	for i := range departments {
		if departments[i].ParentID == parentID {
			children := s.buildDepartmentTree(ctx, departments, departments[i].ID)
			departments[i].Children = children

			if err := ctxDB(ctx).Where("department_id = ?", departments[i].ID).Find(&departments[i].Users).Error; err != nil {
				xlog.Error("加载部门成员失败: %v", err)
			}
			tree = append(tree, departments[i])
		}
	}
	return tree
}

func (s *departmentService) AddUsers(ctx context.Context, departmentID int, userIDs []int) error {
	var department models.Department
	if err := ctxDB(ctx).First(&department, departmentID).Error; err != nil {
		return errors.New("部门不存在")
	}

	if len(userIDs) == 0 {
		return nil
	}

	return ctxDB(ctx).Model(&models.User{}).Where("id IN ?", userIDs).Update("department_id", departmentID).Error
}

func (s *departmentService) RemoveUser(ctx context.Context, userID int) error {
	var user models.User
	if err := ctxDB(ctx).First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	db := ctxDB(ctx)

	// 如果该用户是其所属部门的负责人，先将对应部门的负责人置空
	if user.DepartmentID != 0 {
		if err := db.Model(&models.Department{}).
			Where("id = ? AND leader_id = ?", user.DepartmentID, userID).
			Update("leader_id", nil).Error; err != nil {
			return err
		}
	}

	// 再将用户移出部门
	return db.Model(&models.User{}).Where("id = ?", userID).Update("department_id", 0).Error
}
