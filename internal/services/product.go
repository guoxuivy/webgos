package services

import (
	"strconv"

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

// CreateProduct 创建商品（使用 BaseModel）
func (s *productService) CreateProduct(product *models.Product) error {
	var productModel models.Product
	return productModel.Create(product)
}

// GetProductByID 根据ID获取商品（使用 BaseModel）
func (s *productService) GetProductByID(id string) (*models.Product, error) {
	var productModel models.Product
	
	// 转换ID为整数
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	
	product, err := productModel.Read(productID)
	if err != nil {
		return nil, err
	}
	return product, nil
}
