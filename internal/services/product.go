package services

import (
	"webgos/internal/database"
	"webgos/internal/models"
)

// ProductService 商品服务接口
type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
}

// productService 实现 ProductService 接口
type productService struct{}

// NewProductService 创建商品服务实例
func NewProductService() ProductService {
	return &productService{}
}

// CreateProduct 创建商品
func (s *productService) CreateProduct(product *models.Product) error {
	return database.DB.Create(product).Error
}

// GetProductByID 根据ID获取商品
func (s *productService) GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
