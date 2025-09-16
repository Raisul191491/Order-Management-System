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

type CityHandler struct {
	cityService domain.CityService
}

func NewCityHandler(cityService domain.CityService) *CityHandler {
	return &CityHandler{cityService: cityService}
}

func (handler CityHandler) CreateCity(ctx *gin.Context) {
	var req types.CityCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	if req.BaseDeliveryFee < 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "base_delivery_fee cannot be negative", nil)
		return
	}

	err := handler.cityService.CreateCity(req)
	if err != nil {
		if err.Error() == "city with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "city with name already exists", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create city", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created city", nil)
}

func (handler CityHandler) GetCityByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid city ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.cityService.GetCityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "city with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "city not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch city", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched city", response)
}

func (handler CityHandler) GetAllCities(ctx *gin.Context) {
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

	responses, err := handler.cityService.GetAllCities(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch cities", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched cities", responses)
}

func (handler CityHandler) UpdateCity(ctx *gin.Context) {
	var req types.CityUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	if req.BaseDeliveryFee < 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "base_delivery_fee cannot be negative", nil)
		return
	}

	err := handler.cityService.UpdateCity(req)
	if err != nil {
		if err.Error() == "city with name '"+req.Name+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "city with name already exists", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update city", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated city", nil)
}

func (handler CityHandler) DeleteCity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid city ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.cityService.DeleteCity(id)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete city", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted city", nil)
}

func (handler CityHandler) GetCityByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "City name is required", nil)
		return
	}

	response, err := handler.cityService.GetCityByName(name)
	if err != nil {
		if err.Error() == "city with name '"+name+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "City not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch city", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched city", response)
}
