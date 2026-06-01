package models

type InventoryRecord struct {
	BaseFields
	ProductID int
	Quantity  int
	Type      string
}