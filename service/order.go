package service

import (
	"fmt"
	"oms/consts"
	"oms/domain"
	"oms/model"
	"oms/types"
	"time"
)

// OrderService provides business logic for order operations
type orderService struct {
	orderRepository domain.OrderRepository
	storeService    domain.StoreService
	cityService     domain.CityService
}

func NewOrderService(
	orderRepository domain.OrderRepository,
	storeService domain.StoreService,
	cityService domain.CityService) domain.OrderService {
	return &orderService{
		orderRepository: orderRepository,
		storeService:    storeService,
		cityService:     cityService,
	}
}

func (os orderService) CreateOrder(order types.OrderCreateRequest) (types.OrderCreateResponse, error) {
	if _, err := os.storeService.GetStoreByID(order.StoreID); err != nil {
		return types.OrderCreateResponse{}, fmt.Errorf("invalid store_id: store not found")
	}

	// Generate unique consignment ID
	consignmentID := generateConsignmentID()

	// Calculate delivery fee (this would typically involve business logic)
	deliveryFee := os.calculateDeliveryFee(order)

	// Calculate COD fee
	codFee := calculateCodFee(order.OrderAmount)

	// Calculate total fee
	totalFee := deliveryFee + codFee - order.PromoDiscount - order.Discount

	newOrder := model.Order{
		ConsignmentID:      consignmentID,
		UserID:             order.UserId,
		StoreID:            order.StoreID,
		MerchantOrderID:    order.MerchantOrderID,
		RecipientName:      order.RecipientName,
		RecipientPhone:     order.RecipientPhone,
		RecipientAddress:   order.RecipientAddress,
		RecipientCity:      order.RecipientCity,
		RecipientZone:      order.RecipientZone,
		RecipientArea:      order.RecipientArea,
		OrderType:          consts.OrderTypeDelivery,
		DeliveryTypeID:     order.DeliveryType,
		ItemType:           order.ItemType,
		ItemQuantity:       order.ItemQuantity,
		ItemWeight:         order.ItemWeight,
		ItemDescription:    order.ItemDescription,
		SpecialInstruction: order.SpecialInstruction,
		OrderAmount:        order.OrderAmount,
		AmountToCollect:    order.OrderAmount + totalFee,
		DeliveryFee:        deliveryFee,
		CodFee:             codFee,
		PromoDiscount:      order.PromoDiscount,
		Discount:           order.Discount,
		TotalFee:           totalFee,
		OrderStatus:        consts.OrderStatusPending,
	}

	err := os.orderRepository.CreateOrder(newOrder)
	if err != nil {
		return types.OrderCreateResponse{}, err
	}

	response := types.OrderCreateResponse{
		ConsignmentID:   consignmentID,
		MerchantOrderID: order.MerchantOrderID,
		OrderStatus:     newOrder.OrderStatus,
		DeliveryFee:     deliveryFee,
	}

	return response, nil
}

func (os orderService) GetOrderByConsignmentID(consignmentID string, userId int64) (types.OrderResponse, error) {
	existingOrder, err := os.orderRepository.GetOrderByConsignmentID(consignmentID)
	if err != nil {
		return types.OrderResponse{}, err
	}

	if existingOrder.UserID != userId {
		return types.OrderResponse{}, fmt.Errorf("unauthorized")
	}

	return os.mapOrderToResponse(existingOrder), nil
}

func (os orderService) ListAllOrders(listReq types.OrderListRequest) (types.OrderListResponse, error) {
	orders, pagination, err := os.orderRepository.ListAllOrders(listReq)
	if err != nil {
		return types.OrderListResponse{}, err
	}

	var orderResponses []types.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, os.mapOrderToResponse(order))
	}

	response := types.OrderListResponse{
		Orders:     orderResponses,
		Pagination: pagination,
	}

	return response, nil
}

func (os orderService) UpdateOrder(order types.OrderUpdateRequest) error {
	existingOrder, err := os.orderRepository.GetOrderByConsignmentID(order.ConsignmentID)
	if err != nil {
		return err
	}

	if existingOrder.UserID != order.UserId {
		return fmt.Errorf("unauthorized")
	}

	// Update fields if provided
	if order.RecipientName != "" {
		existingOrder.RecipientName = order.RecipientName
	}
	if order.RecipientPhone != "" {
		existingOrder.RecipientPhone = order.RecipientPhone
	}
	if order.RecipientAddress != "" {
		existingOrder.RecipientAddress = order.RecipientAddress
	}
	if order.ItemWeight != 0 {
		existingOrder.ItemWeight = order.ItemWeight
	}
	if order.OrderAmount != 0 {
		existingOrder.OrderAmount = order.OrderAmount
		// Recalculate fees
		existingOrder.DeliveryFee = os.calculateDeliveryFeeFromOrder(existingOrder)
		existingOrder.CodFee = calculateCodFee(existingOrder.OrderAmount)
		existingOrder.TotalFee = existingOrder.DeliveryFee + existingOrder.CodFee - existingOrder.PromoDiscount - existingOrder.Discount
		existingOrder.AmountToCollect = existingOrder.AmountToCollect + existingOrder.TotalFee
	}
	if order.SpecialInstruction != "" {
		existingOrder.SpecialInstruction = order.SpecialInstruction
	}

	return os.orderRepository.UpdateOrder(existingOrder)
}

func (os orderService) UpdateOrderStatus(updateReq types.OrderStatusUpdateRequest, status string) error {
	existingOrder, err := os.orderRepository.GetOrderByConsignmentID(updateReq.ConsignmentID)
	if err != nil {
		return err
	}

	if existingOrder.UserID != updateReq.UserId {
		return fmt.Errorf("unauthorized")
	}

	return os.orderRepository.UpdateOrderStatus(existingOrder.ID, status)
}

func (os orderService) DeleteOrder(consignmentID string, userId int64) error {
	existingOrder, err := os.orderRepository.GetOrderByConsignmentID(consignmentID)
	if err != nil {
		return err
	}

	if existingOrder.UserID != userId {
		return fmt.Errorf("unauthorized")
	}

	return os.orderRepository.DeleteOrder(existingOrder.ID)
}

func generateConsignmentID() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("CON%d", timestamp)
}

func (os orderService) calculateDeliveryFee(order types.OrderCreateRequest) float64 {
	baseFee := 60.0

	city, err := os.cityService.GetCityByID(order.RecipientCity)
	if err != nil {
		return baseFee
	}

	baseFee = city.BaseDeliveryFee

	// Add weight-based fee
	if order.ItemWeight > 1.0 {
		baseFee += (order.ItemWeight - 1.0) * 10.0
	}

	return baseFee
}

func (os orderService) calculateDeliveryFeeFromOrder(order model.Order) float64 {
	baseFee := 60.0

	city, err := os.cityService.GetCityByID(order.RecipientCity)
	if err != nil {
		return baseFee
	}

	baseFee = city.BaseDeliveryFee

	// Add weight-based fee
	if order.ItemWeight > 1.0 {
		baseFee += (order.ItemWeight - 1.0) * 10.0
	}

	return baseFee
}

func calculateCodFee(orderAmount float64) float64 {
	codFeePercentage := 0.01 // 1%
	return orderAmount * codFeePercentage
}

func (os orderService) mapOrderToResponse(order model.Order) types.OrderResponse {
	return types.OrderResponse{
		ConsignmentID:    order.ConsignmentID,
		OrderCreatedAt:   order.CreatedAt,
		OrderDescription: order.ItemDescription,
		MerchantOrderID:  order.MerchantOrderID,
		RecipientName:    order.RecipientName,
		RecipientAddress: order.RecipientAddress,
		RecipientPhone:   order.RecipientPhone,
		OrderAmount:      order.OrderAmount,
		TotalFee:         order.TotalFee,
		Instruction:      order.SpecialInstruction,
		CodFee:           order.CodFee,
		PromoDiscount:    order.PromoDiscount,
		Discount:         order.Discount,
		DeliveryFee:      order.DeliveryFee,
		OrderStatus:      order.OrderStatus,
		OrderType:        order.OrderType,
		ItemType:         order.ItemType,
	}
}
