package dto

type StockInDTO struct {
	WarehouseID int   `json:"warehouse_id" validate:"required" label:"仓库ID"`
	InstanceIDs []int `json:"instance_ids" validate:"required,min=1" label:"商品实例ID列表"`
	OperatorID  int   `json:"operator_id" validate:"required" label:"操作人ID"`
}
