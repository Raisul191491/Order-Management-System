package domain

import "oms/types"

type AuthService interface {
	Login(loginRequest types.UserLoginRequest) (types.UserLoginResponse, error)
	Logout(accessToken string) error
}
