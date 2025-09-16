package repository

import (
	"errors"
	"fmt"
	"math"
	"oms/domain"
	"oms/model"
	"oms/types"

	"gorm.io/gorm"
)

type orderRepository struct {
	masterDb  *gorm.DB
	replicaDb *gorm.DB
}

func NewOrderRepository(masterDB, replicaDB *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		masterDb:  masterDB,
		replicaDb: replicaDB,
	}
}

func (r *orderRepository) CreateOrder(order model.Order) error {
	return r.masterDb.Create(&order).Error
}

func (r *orderRepository) GetOrderByConsignmentID(consignmentID string) (model.Order, error) {
	var order model.Order
	err := r.replicaDb.Where("consignment_id = ?", consignmentID).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Order{}, fmt.Errorf("order with consignment ID '%s' not found", consignmentID)
		}
		return model.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) ListAllOrders(listReq types.OrderListRequest) ([]model.Order, model.Pagination, error) {
	var orders []model.Order
	query := r.replicaDb.Model(&model.Order{}).Where("deleted_at IS NULL")

	if listReq.UserId != 0 {
		query = query.Where("user_id = ?", listReq.UserId)
	}

	if listReq.OrderStatus != "" {
		query = query.Where("order_status = ?", listReq.OrderStatus)
	}

	pageNumber := listReq.PageNumber
	pageLength := listReq.PageLength
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageLength <= 0 {
		pageLength = 10
	}

	// Count total
	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return nil, model.Pagination{}, err
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(pageLength)))

	// Apply pagination
	offset := (pageNumber - 1) * pageLength
	query = query.Offset(offset).Limit(pageLength)

	// Execute query
	if err := query.Find(&orders).Error; err != nil {
		return nil, model.Pagination{}, err
	}

	pagination := model.Pagination{
		Total:       totalPages,
		CurrentPage: pageNumber,
		TotalInPage: totalRows,
		PerPage:     pageLength,
		LastPage:    totalPages,
	}

	return orders, pagination, nil
}

func (r *orderRepository) UpdateOrder(order model.Order) error {
	result := r.masterDb.Save(&order)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", order.ID)
	}

	return nil
}

func (r *orderRepository) UpdateOrderStatus(id int64, status string) error {
	result := r.masterDb.Model(&model.Order{}).Where("id = ?", id).Update("order_status", status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}

func (r *orderRepository) DeleteOrder(id int64) error {
	result := r.masterDb.Delete(&model.Order{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}
