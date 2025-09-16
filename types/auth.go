package types

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	UserResponse
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type PasswordUpdateRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,oldPassword"`
	NewPassword string `json:"newPassword" binding:"required,newPassword"`
}
