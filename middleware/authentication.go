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
		userToken := ctx.Request.Header.Get("Authorization")
		if userToken == "" {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", nil)
			return
		}

		_, err := userSessionSvc.ValidateSession(userToken)
		if err != nil {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", []any{err.Error()})
		}

		claims, err := utility.VerifyJWT(userToken)
		if err != nil {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", []any{err.Error()})
			return
		}

		if claims == nil || claims.UserID == 0 {
			utility.SendErrorResponse(ctx, http.StatusForbidden, "Unauthorized", nil)
		}
		ctx.Set(consts.UserIdKey, claims.UserID)
		ctx.Next()
	}
}
