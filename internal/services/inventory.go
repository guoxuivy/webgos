package services

import (
	"errors"
	"webgos/internal/database"
	"webgos/internal/models"
)

// InventoryService 库存服务接口
type InventoryService interface {
	ProductIn(record *models.InventoryRecord) error
	ProductOut(record *models.InventoryRecord) error
}

// inventoryService 实现 InventoryService 接口
type inventoryService struct{}

// NewInventoryService 创建库存服务实例
func NewInventoryService() InventoryService {
	return &inventoryService{}
}

// ProductIn 处理商品入库
func (s *inventoryService) ProductIn(record *models.InventoryRecord) error {
	record.Type = "in"
	if err := database.DB.Create(record).Error; err != nil {
		return err
	}

	// 更新商品库存
	var product models.Product
	if err := database.DB.First(&product, record.ProductID).Error; err != nil {
		return err
	}
	product.Stock += record.Quantity
	return database.DB.Save(&product).Error
}

// ProductOut 处理商品出库
func (s *inventoryService) ProductOut(record *models.InventoryRecord) error {
	record.Type = "out"
	var product models.Product
	if err := database.DB.First(&product, record.ProductID).Error; err != nil {
		return err
	}

	if product.Stock < record.Quantity {
		return errors.New("库存不足")
	}

	if err := database.DB.Create(record).Error; err != nil {
		return err
	}

	// 更新商品库存
	product.Stock -= record.Quantity
	return database.DB.Save(&product).Error
}
