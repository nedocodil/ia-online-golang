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
	Password string `json:"password" validate:"required,complexpassword"`
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
