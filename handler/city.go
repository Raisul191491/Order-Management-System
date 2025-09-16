package handler

import (
	"errors"
	"net/http"
	"oms/domain"
	"oms/types"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BaseDeliveryFee < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "base_delivery_fee cannot be negative"})
		return
	}

	err := handler.cityService.CreateCity(req)
	if err != nil {
		if err.Error() == "city with name '"+req.Name+"' already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, "successfully created city")
}

func (handler CityHandler) GetCityByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid city ID format"})
		return
	}

	if id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id should be positive"})
		return
	}

	response, err := handler.cityService.GetCityByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "city with ID "+idStr+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (handler CityHandler) GetAllCities(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	responses, err := handler.cityService.GetAllCities(limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   responses,
		"total":  len(responses),
		"limit":  limit,
		"offset": offset,
	})
}

func (handler CityHandler) UpdateCity(ctx *gin.Context) {
	var req types.CityUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BaseDeliveryFee < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "base_delivery_fee cannot be negative"})
		return
	}

	err := handler.cityService.UpdateCity(req)
	if err != nil {
		if err.Error() == "city with name '"+req.Name+"' already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "successfully updated city")
}

func (handler CityHandler) DeleteCity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id should be positive"})
		return
	}

	err = handler.cityService.DeleteCity(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (handler CityHandler) GetCityByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "city name is required"})
		return
	}

	response, err := handler.cityService.GetCityByName(name)
	if err != nil {
		if err.Error() == "city with name '"+name+"' not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
