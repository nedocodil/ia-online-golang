package dto

type UserDTO struct {
	ID             *int64   `json:"id" validate:"omitempty"`
	Roles          []string `json:"roles" validate:"omitempty"`
	ReferralCode   string   `json:"referral_code" validate:"omitempty"`
	Email          string   `json:"email" validate:"omitempty"`
	Name           string   `json:"name" validate:"omitempty"`
	PhoneNumber    string   `json:"phone_number" validate:"omitempty"`
	Telegram       string   `json:"telegram" validate:"omitempty"`
	City           string   `json:"city" validate:"omitempty"`
	RewardInternet float64  `json:"reward_internet" validate:"omitempty"`
	RewardCleaning float64  `json:"reward_cleaning" validate:"omitempty"`
	RewardShipping float64  `json:"reward_shipping" validate:"omitempty"`
	RewardReferral float64  `json:"reward_referral" validate:"omitempty"`
}
