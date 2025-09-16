package handler

import (
	"net/http"
	"oms/consts"
	"oms/domain"
	"oms/types"
	"oms/utility"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler AuthHandler) Login(ctx *gin.Context) {
	var req types.UserLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utility.SendErrorResponse(ctx, http.StatusBadRequest, "Unable to bind request", []any{err.Error()})
		return
	}

	response, err := handler.authService.Login(req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "The user credentials were incorrect.", []any{err.Error()})
			return
		}
		if err.Error() == "failed to create session" {
			utility.SendErrorResponse(ctx, http.StatusInternalServerError, "failed to create session", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to login", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully logged in", response)
}

func (handler AuthHandler) Logout(ctx *gin.Context) {
	accessToken := ctx.GetString(consts.AccessTokenKey)
	if accessToken == "" {
		utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Access token required", nil)
		return
	}

	err := handler.authService.Logout(accessToken)
	if err != nil {
		if err.Error() == "invalid access token" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "invalid access token", []any{err.Error()})
			return
		}
		utility.SendErrorResponse(ctx, http.StatusInternalServerError, "Unable to logout", []any{err.Error()})
		return
	}

	utility.SendSuccessResponse(ctx, http.StatusOK, "Successfully logged out", nil)
}
