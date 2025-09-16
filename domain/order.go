package domain

import (
	"oms/model"
	"oms/types"
)

type OrderRepository interface {
	CreateOrder(order model.Order) error
	GetOrderByConsignmentID(consignmentID string) (model.Order, error)
	ListAllOrders(listReq types.OrderListRequest) ([]model.Order, model.Pagination, error)
	UpdateOrder(order model.Order) error
	UpdateOrderStatus(id int64, status string) error
	DeleteOrder(id int64) error
}

type OrderService interface {
	CreateOrder(order types.OrderCreateRequest) (types.OrderCreateResponse, error)
	GetOrderByConsignmentID(consignmentID string, userId int64) (types.OrderResponse, error)
	ListAllOrders(listReq types.OrderListRequest) (types.OrderListResponse, error)
	UpdateOrder(order types.OrderUpdateRequest) error
	UpdateOrderStatus(updateReq types.OrderStatusUpdateRequest, status string) error
	DeleteOrder(consignmentID string, userId int64) error
}
