package handler

import (
	"errors"
	"net/http"
	"oms/domain"
	"oms/types"
	"oms/utility"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeliveryTypeHandler struct {
	deliveryTypeService domain.DeliveryTypeService
}

func NewDeliveryTypeHandler(deliveryTypeService domain.DeliveryTypeService) *DeliveryTypeHandler {
	return &DeliveryTypeHandler{deliveryTypeService: deliveryTypeService}
}

func (handler DeliveryTypeHandler) CreateDeliveryType(ctx *gin.Context) {
	var req types.DeliveryTypeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.deliveryTypeService.CreateDeliveryType(req)
	if err != nil {
		if err.Error() == "delivery type with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "delivery type with name already exists", []any{err.Error()})
			return
		}
		if err.Error() == "delivery type name cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "delivery type name cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create delivery type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created delivery type", nil)
}

func (handler DeliveryTypeHandler) GetDeliveryTypeByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid delivery type ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.deliveryTypeService.GetDeliveryTypeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "delivery type with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "delivery type not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch delivery type", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched delivery type", response)
}

func (handler DeliveryTypeHandler) GetAllDeliveryTypes(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "invalid limit parameter", nil)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "invalid offset parameter", nil)
		return
	}

	responses, err := handler.deliveryTypeService.GetAllDeliveryTypes(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch delivery types", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched delivery types", responses)
}

func (handler DeliveryTypeHandler) UpdateDeliveryType(ctx *gin.Context) {
	var req types.DeliveryTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.deliveryTypeService.UpdateDeliveryType(req)
	if err != nil {
		if err.Error() == "delivery type with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "delivery type with name already exists", []any{err.Error()})
			return
		}
		if err.Error() == "delivery type name cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "delivery type name cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update delivery type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated delivery type", nil)
}

func (handler DeliveryTypeHandler) DeleteDeliveryType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid delivery type ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.deliveryTypeService.DeleteDeliveryType(id)
	if err != nil {
		if err.Error() == "delivery type does not exist" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "delivery type not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete delivery type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted delivery type", nil)
}
