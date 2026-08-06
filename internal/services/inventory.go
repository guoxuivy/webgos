package services

import (
	"context"
	"errors"

	"webgos/internal/models"
	"webgos/internal/xdb"
)

type InventoryService interface {
	ProductIn(ctx context.Context, record *models.InventoryRecord) error
	ProductOut(ctx context.Context, record *models.InventoryRecord) error
}

type inventoryService struct{}

func NewInventoryService() InventoryService {
	return &inventoryService{}
}

func (s *inventoryService) ProductIn(ctx context.Context, record *models.InventoryRecord) error {
	if record.ProductID == 0 || record.Quantity <= 0 {
		return errors.New("产品ID和数量必须大于0")
	}
	record.Type = "in"

	db := xdb.GetDB().WithContext(ctx)

	if err := db.Create(record).Error; err != nil {
		return err
	}

	var product models.Product
	if err := db.First(&product, record.ProductID).Error; err != nil {
		return err
	}

	product.Stock += record.Quantity
	return db.Updates(&product).Error
}

func (s *inventoryService) ProductOut(ctx context.Context, record *models.InventoryRecord) error {
	if record.ProductID == 0 || record.Quantity <= 0 {
		return errors.New("产品ID和数量必须大于0")
	}
	record.Type = "out"

	db := xdb.GetDB().WithContext(ctx)

	var product models.Product
	if err := db.First(&product, record.ProductID).Error; err != nil {
		return err
	}

	if product.Stock < record.Quantity {
		return errors.New("库存不足")
	}

	if err := db.Create(record).Error; err != nil {
		return err
	}

	product.Stock -= record.Quantity
	return db.Updates(&product).Error
}
