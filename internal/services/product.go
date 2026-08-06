package services

import (
	"context"
	"errors"
	"strconv"

	"webgos/internal/models"
	"webgos/internal/xdb"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
}

type productService struct{}

func NewProductService() ProductService {
	return &productService{}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return errors.New("产品名称不能为空")
	}
	return xdb.GetDB().WithContext(ctx).Create(product).Error
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = xdb.GetDB().WithContext(ctx).First(&product, productID).Error
	return &product, err
}
