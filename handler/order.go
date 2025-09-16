package handler

import (
	"errors"
	"net/http"
	"oms/consts"
	"oms/domain"
	"oms/types"
	"oms/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	orderService domain.OrderService
}

func NewOrderHandler(orderService domain.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (handler OrderHandler) CreateOrder(ctx *gin.Context) {
	var req types.OrderCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusUnprocessableEntity, "Unable the bind request", nil)
		return
	}

	if validationErrors := req.Validate(); validationErrors != nil {
		utility.SendErrorResponse(ctx, http.StatusUnprocessableEntity, "Please fix the given errors", validationErrors.Errors)
		return
	}

	// Extract user ID from JWT token context (assuming middleware sets this)
	userID, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userID.(int64); ok {
			req.UserId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	response, err := handler.orderService.CreateOrder(req)
	if err != nil {
		// Handle specific validation errors
		if err.Error() == "invalid store_id: store not found" {
			utility.SendErrorResponse(ctx, http.StatusUnprocessableEntity, "Please fix the given errors", map[string][]string{
				"store_id": {"The store field is required", "Wrong Store selected"},
			})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create order", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Order Created Successfully", response)
}

func (handler OrderHandler) GetOrderByConsignmentID(ctx *gin.Context) {
	consignmentID := ctx.Param("consignment_id")
	if consignmentID == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Consignment ID is required", nil)
		return
	}

	var userId int64

	// Extract user ID from JWT token context (assuming middleware sets this)
	userIDAny, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userIDAny.(int64); ok {
			userId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	response, err := handler.orderService.GetOrderByConsignmentID(consignmentID, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order with consignment ID '"+consignmentID+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "Order not found", []any{err.Error()})
			return
		}
		if err.Error() == "unauthorized" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch order", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Order successfully fetched", response)
}

func (handler OrderHandler) ListAllOrders(ctx *gin.Context) {
	var req types.OrderListRequest

	// Parse query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters", []any{err.Error()})
		return
	}

	// Set defaults if not provided
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageLength <= 0 {
		req.PageLength = 10
	}

	// Extract user ID from JWT token context (assuming middleware sets this)
	userID, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userID.(int64); ok {
			req.UserId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	response, err := handler.orderService.ListAllOrders(req)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch orders", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched orders", response)
}

func (handler OrderHandler) UpdateOrder(ctx *gin.Context) {
	var req types.OrderUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	if req.ConsignmentID == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Consignment ID is required", nil)
		return
	}

	// Extract user ID from JWT token context (assuming middleware sets this)
	userID, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userID.(int64); ok {
			req.UserId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	err := handler.orderService.UpdateOrder(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order with consignment ID '"+req.ConsignmentID+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "Order not found", []any{err.Error()})
			return
		}
		if err.Error() == "unauthorized" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update order", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated order", nil)
}

func (handler OrderHandler) CancelOrder(ctx *gin.Context) {
	var req types.OrderStatusUpdateRequest

	consignmentID := ctx.Param("consignment_id")
	if consignmentID == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Consignment ID is required", nil)
		return
	}

	req.ConsignmentID = consignmentID

	// Extract user ID from JWT token context (assuming middleware sets this)
	userID, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userID.(int64); ok {
			req.UserId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	err := handler.orderService.UpdateOrderStatus(req, consts.OrderStatusCancelled)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order with consignment ID '"+req.ConsignmentID+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "Order not found", []any{err.Error()})
			return
		}
		if err.Error() == "unauthorized" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update order status", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Order Cancelled Successfully", nil)
}

func (handler OrderHandler) DeleteOrder(ctx *gin.Context) {
	consignmentID := ctx.Param("consignment_id")
	if consignmentID == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Consignment ID is required", nil)
		return
	}

	var userId int64

	// Extract user ID from JWT token context (assuming middleware sets this)
	userIDAny, exists := ctx.Get(consts.UserIdKey)
	if exists {
		if id, ok := userIDAny.(int64); ok {
			userId = id
		}
	}
	if !exists {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Please login first...", nil)
	}

	err := handler.orderService.DeleteOrder(consignmentID, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			err.Error() == "order with consignment ID '"+consignmentID+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "Order not found", []any{err.Error()})
			return
		}
		if err.Error() == "unauthorized" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete order", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted order", nil)
}
