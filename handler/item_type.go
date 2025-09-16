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

type ItemTypeHandler struct {
	itemTypeService domain.ItemTypeService
}

func NewItemTypeHandler(itemTypeService domain.ItemTypeService) *ItemTypeHandler {
	return &ItemTypeHandler{itemTypeService: itemTypeService}
}

func (handler ItemTypeHandler) CreateItemType(ctx *gin.Context) {
	var req types.ItemTypeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.itemTypeService.CreateItemType(req)
	if err != nil {
		if err.Error() == "item type with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "item type with name already exists", []any{err.Error()})
			return
		}
		if err.Error() == "item type name cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "item type name cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create item type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created item type", nil)
}

func (handler ItemTypeHandler) GetItemTypeByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid item type ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.itemTypeService.GetItemTypeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "item type with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "item type not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch item type", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched item type", response)
}

func (handler ItemTypeHandler) GetAllItemTypes(ctx *gin.Context) {
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

	responses, err := handler.itemTypeService.GetAllItemTypes(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch item types", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched item types", responses)
}

func (handler ItemTypeHandler) UpdateItemType(ctx *gin.Context) {
	var req types.ItemTypeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.itemTypeService.UpdateItemType(req)
	if err != nil {
		if err.Error() == "item type with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "item type with name already exists", []any{err.Error()})
			return
		}
		if err.Error() == "item type name cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "item type name cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update item type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated item type", nil)
}

func (handler ItemTypeHandler) DeleteItemType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid item type ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.itemTypeService.DeleteItemType(id)
	if err != nil {
		if err.Error() == "item type does not exist" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "item type not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete item type", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted item type", nil)
}
