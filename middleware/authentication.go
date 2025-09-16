package middleware

import (
	"net/http"
	"oms/consts"
	"oms/domain"
	"oms/utility"

	"github.com/gin-gonic/gin"
)

func Auth(userSessionSvc domain.UserSessionService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utility.SendErrorResponse(ctx, http.StatusUnauthorized, "Authorization header required", nil)
			ctx.Abort()
			return
		}

		userToken := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			userToken = authHeader[7:]
		}

		if userToken == "" {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		_, err := userSessionSvc.ValidateSession(userToken)
		if err != nil {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", []any{err.Error()})
			ctx.Abort()
			return
		}

		claims, err := utility.VerifyJWT(userToken)
		if err != nil {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", []any{err.Error()})
			ctx.Abort()
			return
		}

		if claims == nil || claims.UserID == 0 {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", nil)
			ctx.Abort()
			return
		}
		ctx.Set(consts.UserIdKey, claims.UserID)
		ctx.Set(consts.AccessTokenKey, userToken)
		ctx.Next()
	}
}
