package dto

type UserDTO struct {
	ID             int64   `json:"id"`
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	PhoneNumber    string  `json:"phone_number"`
	City           string  `json:"city"`
	RewardInternet float64 `json:"reward_internet"`
	RewardCleaning float64 `json:"reward_clean"`
	RewardShipping float64 `json:"reward_shipping"`
}
