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

type StoreHandler struct {
	storeService domain.StoreService
}

func NewStoreHandler(storeService domain.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

func (handler StoreHandler) CreateStore(ctx *gin.Context) {
	var req types.StoreCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.storeService.CreateStore(req)
	if err != nil {
		if err.Error() == "store with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "store with name already exists", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create store", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created store", nil)
}

func (handler StoreHandler) GetStoreByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid store ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.storeService.GetStoreByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "store with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "store not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch store", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched store", response)
}

func (handler StoreHandler) GetAllStores(ctx *gin.Context) {
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

	responses, err := handler.storeService.GetAllStores(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch stores", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched stores", responses)
}

func (handler StoreHandler) UpdateStore(ctx *gin.Context) {
	var req types.StoreUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.storeService.UpdateStore(req)
	if err != nil {
		if err.Error() == "store with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "store with name already exists", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update store", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated store", nil)
}

func (handler StoreHandler) DeleteStore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid store ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.storeService.DeleteStore(id)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete store", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted store", nil)
}
