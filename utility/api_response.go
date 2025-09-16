package utility

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Message string        `json:"message"`
	Type    string        `json:"type"`
	Code    int           `json:"code"`
	Data    interface{}   `json:"data,omitempty"`
	Errors  []interface{} `json:"errors,omitempty"`
}

func SendSuccessResponse(ctx *gin.Context, code int, message string, data interface{}) {
	response := APIResponse{
		Message: message,
		Type:    "success",
		Code:    code,
		Data:    data,
	}

	ctx.JSON(code, response)
}

func SendErrorResponse(ctx *gin.Context, code int, message string, errs []interface{}) {
	response := APIResponse{
		Message: message,
		Type:    "error",
		Code:    code,
		Errors:  errs,
	}

	ctx.JSON(code, response)
}
