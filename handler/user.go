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

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler UserHandler) CreateUser(ctx *gin.Context) {
	var req types.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "user with email '"+req.Email+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "user with email already exists", []any{err.Error()})
			return
		}
		if err.Error() == "email cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "email cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to create user", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusCreated, "Successfully created user", nil)
}

func (handler UserHandler) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	response, err := handler.userService.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "user with ID "+idStr+" not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "user not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch user", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched user", response)
}

func (handler UserHandler) GetAllUsers(ctx *gin.Context) {
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

	responses, err := handler.userService.GetAllUsers(limit, offset)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch users", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched users", responses)
}

func (handler UserHandler) UpdateUserEmail(ctx *gin.Context) {
	var req types.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	err := handler.userService.UpdateUserEmail(req)
	if err != nil {
		if err.Error() == "user with email '"+req.Email+"' already exists" {
			utility.SendErrorResponse(ctx, http.StatusConflict, "user with email already exists", []any{err.Error()})
			return
		}
		if err.Error() == "email cannot be empty" {
			utility.SendErrorResponse(ctx, http.StatusBadRequest, "email cannot be empty", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to update user", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully updated user", nil)
}

func (handler UserHandler) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format", []any{err.Error()})
		return
	}

	if id <= 0 {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Id should be positive", nil)
		return
	}

	err = handler.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user does not exist" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "user not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to delete user", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully deleted user", nil)
}

func (handler UserHandler) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	if email == "" {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Email parameter is required", nil)
		return
	}

	response, err := handler.userService.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user with email '"+email+"' not found" {
			utility.SendErrorResponse(ctx, http.StatusNotFound, "user not found", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to fetch user", []any{err.Error()})
		return
	}
	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully fetched user", response)
}
