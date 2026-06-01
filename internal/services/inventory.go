package services

import (
	"errors"

	"webgos/internal/xdb"
	"webgos/internal/models"
)

type InventoryService interface {
	ProductIn(record *models.InventoryRecord) error
	ProductOut(record *models.InventoryRecord) error
}

type inventoryService struct{}

func NewInventoryService() InventoryService {
	return &inventoryService{}
}

func (s *inventoryService) ProductIn(record *models.InventoryRecord) error {
	if record.ProductID == 0 || record.Quantity <= 0 {
		return errors.New("产品ID和数量必须大于0")
	}
	record.Type = "in"

	if err := xdb.GetDB().Create(record).Error; err != nil {
		return err
	}

	var product models.Product
	if err := xdb.GetDB().First(&product, record.ProductID).Error; err != nil {
		return err
	}

	product.Stock += record.Quantity
	return xdb.GetDB().Updates(&product).Error
}

func (s *inventoryService) ProductOut(record *models.InventoryRecord) error {
	if record.ProductID == 0 || record.Quantity <= 0 {
		return errors.New("产品ID和数量必须大于0")
	}
	record.Type = "out"

	var product models.Product
	if err := xdb.GetDB().First(&product, record.ProductID).Error; err != nil {
		return err
	}

	if product.Stock < record.Quantity {
		return errors.New("库存不足")
	}

	if err := xdb.GetDB().Create(record).Error; err != nil {
		return err
	}

	product.Stock -= record.Quantity
	return xdb.GetDB().Updates(&product).Error
}
