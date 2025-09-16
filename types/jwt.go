package types

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}
