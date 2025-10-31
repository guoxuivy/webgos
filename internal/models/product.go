package models

// Product 商品模型
type Product struct {
	BaseModel[Product]
	Name  string
	Stock int
}
