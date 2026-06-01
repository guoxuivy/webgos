package services

import (
	"errors"
	"strconv"

	"webgos/internal/xdb"
	"webgos/internal/models"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
}

type productService struct{}

func NewProductService() ProductService {
	return &productService{}
}

func (s *productService) CreateProduct(product *models.Product) error {
	if product.Name == "" {
		return errors.New("产品名称不能为空")
	}
	return xdb.GetDB().Create(product).Error
}

func (s *productService) GetProductByID(id string) (*models.Product, error) {
	productID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = xdb.GetDB().First(&product, productID).Error
	return &product, err
}