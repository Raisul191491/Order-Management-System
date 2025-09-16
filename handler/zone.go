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

type ZoneHandler struct {
	zoneService domain.ZoneService
}

func NewZoneHandler(zoneService domain.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (handler ZoneHandler) CreateZone(ctx *gin.Context) {
	var req types.ZoneCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.zoneService.CreateZone(req)
	if err != nil {
		if err.Error() == "zone with name '"+req.Name+"' already exists in this city" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "zone with name already exists in this city", []any{err.Error()})
			return
		}
		if err.Error() == "city with ID "+strconv.FormatInt(req.CityID, 10)+" does not exist" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "invalid city ID", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create zone", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created zone", nil)
}

func (handler ZoneHandler) GetZoneByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid zone ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.zoneService.GetZoneByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "zone with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "zone not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch zone", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched zone", response)
}

func (handler ZoneHandler) GetAllZones(ctx *gin.Context) {
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

	responses, err := handler.zoneService.GetAllZones(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch zones", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched zones", responses)
}

func (handler ZoneHandler) GetZonesByCity(ctx *gin.Context) {
	cityIDStr := ctx.Param("cityId")
	cityID, err := strconv.ParseInt(cityIDStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid city ID format", []any{err.Error()})
		return
	}

	if cityID <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "City ID should be positive", nil)
		return
	}

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

	responses, err := handler.zoneService.GetZonesByCityID(cityID, limit, offset)
	if err != nil {
		if err.Error() == "city with ID "+cityIDStr+" does not exist" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "city not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch zones", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched zones", responses)
}

func (handler ZoneHandler) UpdateZone(ctx *gin.Context) {
	var req types.ZoneUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.zoneService.UpdateZone(req)
	if err != nil {
		if err.Error() == "zone with name '"+req.Name+"' already exists in this city" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "zone with name already exists in this city", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update zone", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated zone", nil)
}

func (handler ZoneHandler) DeleteZone(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid zone ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.zoneService.DeleteZone(id)
	if err != nil {
		if err.Error() == "zone does not exist" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "zone not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete zone", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted zone", nil)
}
