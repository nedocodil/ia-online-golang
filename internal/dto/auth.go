package dto

type RegisterUserDTO struct {
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,complexpassword"`
	RepeatPassword string `json:"repeat_password" validate:"required,eqfield=Password"`
	Telegram       string `json:"telegram" validate:"omitempty"`
	PhoneNumber    string `json:"phone_number" validate:"e164"`
	Name           string `json:"name" validate:"required"`
	City           string `json:"city" validate:"required"`
	ReferralCode   string `json:"referral_code" validate:"omitempty"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type NewPasswordDTO struct {
	OldPassword       string `json:"old_password" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,complexpassword"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type RecoverPasswordDTO struct {
	Code              string `json:"code" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,complexpassword"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type SendNewPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type SendPasswordCodeDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
