package dto

type ReferralDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	Active      bool   `json:"active"`
}
