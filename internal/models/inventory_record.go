package models

// InventoryRecord 进销存记录模型
type InventoryRecord struct {
	BaseModel[InventoryRecord]
	ProductID int
	Quantity  int
	Type      string // "in" 或 "out"
}
