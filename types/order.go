package types

import (
	"fmt"
	"oms/model"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type OrderCreateRequest struct {
	StoreID            int64   `json:"store_id" validate:"required,min=1"`
	MerchantOrderID    string  `json:"merchant_order_id" validate:"omitempty,max=100"`
	RecipientName      string  `json:"recipient_name" validate:"required,min=1,max=255"`
	RecipientPhone     string  `json:"recipient_phone" validate:"required,regexp=^(01)[3-9]{1}[0-9]{8}$"`
	RecipientAddress   string  `json:"recipient_address" validate:"required,min=1"`
	RecipientCity      int64   `json:"recipient_city" validate:"required,min=1"`
	RecipientZone      int64   `json:"recipient_zone" validate:"required,min=1"`
	RecipientArea      string  `json:"recipient_area"`
	DeliveryType       int64   `json:"delivery_type" validate:"required,min=1"`
	ItemType           int64   `json:"item_type" validate:"required,min=1"`
	ItemQuantity       int     `json:"item_quantity" validate:"required,min=1"`
	ItemWeight         float64 `json:"item_weight" validate:"required,gt=0"`
	OrderAmount        float64 `json:"order_amount" validate:"required,gt=0"`
	ItemDescription    string  `json:"item_description"`
	SpecialInstruction string  `json:"special_instruction"`
	PromoDiscount      float64 `json:"promo_discount" validate:"omitempty,gte=0"`
	Discount           float64 `json:"discount" validate:"omitempty,gte=0"`
	UserId             int64   `json:"user_id,omitempty"` // Usually set from JWT token
}

type OrderUpdateRequest struct {
	UserId             int64   `json:"user_id,omitempty"`
	ConsignmentID      string  `json:"consignment_id" validate:"required"`
	MerchantOrderID    string  `json:"merchant_order_id" validate:"omitempty,max=100"`
	RecipientName      string  `json:"recipient_name" validate:"omitempty,min=1,max=255"`
	RecipientPhone     string  `json:"recipient_phone" validate:"omitempty,regexp=^(01)[3-9]{1}[0-9]{8}$"`
	RecipientAddress   string  `json:"recipient_address" validate:"omitempty,min=1"`
	ItemWeight         float64 `json:"item_weight" validate:"omitempty,gt=0"`
	OrderAmount        float64 `json:"order_amount" validate:"omitempty,gt=0"`
	SpecialInstruction string  `json:"special_instruction"`
}

type OrderCreateResponse struct {
	ConsignmentID   string  `json:"consignment_id"`
	MerchantOrderID string  `json:"merchant_order_id"`
	OrderStatus     string  `json:"order_status"`
	DeliveryFee     float64 `json:"delivery_fee"`
}

type OrderResponse struct {
	ConsignmentID    string    `json:"consignment_id"`
	OrderCreatedAt   time.Time `json:"order_created_at"`
	OrderDescription string    `json:"order_description"`
	MerchantOrderID  string    `json:"merchant_order_id"`
	RecipientName    string    `json:"recipient_name"`
	RecipientAddress string    `json:"recipient_address"`
	RecipientPhone   string    `json:"recipient_phone"`
	OrderAmount      float64   `json:"order_amount"`
	TotalFee         float64   `json:"total_fee"`
	Instruction      string    `json:"instruction"`
	OrderType        string    `json:"order_type"`
	CodFee           float64   `json:"cod_fee"`
	PromoDiscount    float64   `json:"promo_discount"`
	Discount         float64   `json:"discount"`
	DeliveryFee      float64   `json:"delivery_fee"`
	OrderStatus      string    `json:"order_status"`
	ItemType         int64     `json:"item_type"`
}

type OrderListRequest struct {
	UserId      int64  `json:"user_id" form:"user_id"`
	OrderStatus string `json:"order_status" form:"order_status"`
	PageNumber  int    `json:"page_number" form:"page" validate:"omitempty,min=1"`
	PageLength  int    `json:"page_length" form:"limit" validate:"omitempty,min=1,max=100"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"data"`
	model.Pagination
}

type OrderStatusUpdateRequest struct {
	ConsignmentID string `json:"consignment_id" form:"consignment_id" validate:"required"`
	UserId        int64  `json:"user_id"`
}

// ValidationErrorResponse represents the error response format
type ValidationErrorResponse struct {
	Message string              `json:"message"`
	Type    string              `json:"type"`
	Code    int                 `json:"code"`
	Errors  map[string][]string `json:"-"`
}

func (r *OrderCreateRequest) Validate() *ValidationErrorResponse {
	validate := validator.New()

	// Register custom validation for Bangladeshi phone number
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		pattern := fl.Param()
		matched, _ := regexp.MatchString(pattern, fl.Field().String())
		return matched
	})

	err := validate.Struct(r)
	if err == nil {
		return nil
	}

	// Create error response
	errorResponse := &ValidationErrorResponse{
		Message: "Please fix the given errors",
		Type:    "error",
		Code:    422,
		Errors:  make(map[string][]string),
	}

	// Process validation errors
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := strings.ToLower(err.Field())
		errorMessage := getErrorMessage(fieldName, err.Tag(), err.Param())

		errorResponse.Errors[fieldName] = append(errorResponse.Errors[fieldName], errorMessage)
	}

	return errorResponse
}

// getErrorMessage returns user-friendly error messages
func getErrorMessage(fieldName, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("The %s field is required.", getFieldDisplayName(fieldName))
	case "min":
		return fmt.Sprintf("The %s must be at least %s.", getFieldDisplayName(fieldName), param)
	case "max":
		return fmt.Sprintf("The %s may not be greater than %s characters.", getFieldDisplayName(fieldName), param)
	case "gt":
		return fmt.Sprintf("The %s must be greater than %s.", getFieldDisplayName(fieldName), param)
	case "gte":
		return fmt.Sprintf("The %s must be greater than or equal to %s.", getFieldDisplayName(fieldName), param)
	case "regexp":
		return "The phone number format is invalid. Must be a valid Bangladeshi phone number (01XXXXXXXXX)."
	default:
		return fmt.Sprintf("The %s field is invalid.", getFieldDisplayName(fieldName))
	}
}

// getFieldDisplayName returns user-friendly field names
func getFieldDisplayName(fieldName string) string {
	displayNames := map[string]string{
		"storeid":            "store",
		"merchantorderid":    "merchant order ID",
		"recipientname":      "recipient name",
		"recipientphone":     "recipient phone",
		"recipientaddress":   "recipient address",
		"recipientcity":      "recipient city",
		"recipientzone":      "recipient zone",
		"recipientarea":      "recipient area",
		"deliverytype":       "delivery type",
		"itemtype":           "item type",
		"itemquantity":       "item quantity",
		"itemweight":         "item weight",
		"orderamount":        "order amount",
		"itemdescription":    "item description",
		"specialinstruction": "special instruction",
		"promodiscount":      "promo discount",
		"discount":           "discount",
		"userid":             "user ID",
	}

	if displayName, exists := displayNames[strings.ToLower(fieldName)]; exists {
		return displayName
	}

	return fieldName
}
