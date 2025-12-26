package services

import (
	"errors"
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

// ProductIn 处理商品入库（使用 BaseModel）
func (s *inventoryService) ProductIn(record *models.InventoryRecord) error {
	record.Type = "in"
	
	// 创建库存记录（使用 BaseModel）
	var inventoryModel models.InventoryRecord
	if err := inventoryModel.Create(record); err != nil {
		return err
	}

	// 更新商品库存（使用 BaseModel）
	var productModel models.Product
	product, err := productModel.Read(record.ProductID)
	if err != nil {
		return err
	}
	
	product.Stock += record.Quantity
	return product.Update(product)
}

// ProductOut 处理商品出库（使用 BaseModel）
func (s *inventoryService) ProductOut(record *models.InventoryRecord) error {
	record.Type = "out"
	
	// 检查商品库存（使用 BaseModel）
	var productModel models.Product
	product, err := productModel.Read(record.ProductID)
	if err != nil {
		return err
	}

	if product.Stock < record.Quantity {
		return errors.New("库存不足")
	}

	// 创建库存记录（使用 BaseModel）
	var inventoryModel models.InventoryRecord
	if err := inventoryModel.Create(record); err != nil {
		return err
	}

	// 更新商品库存（使用 BaseModel）
	product.Stock -= record.Quantity
	return product.Update(product)
}
