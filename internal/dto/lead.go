package dto

import "time"

type LeadDTO struct {
	Name           string  `json:"name" validate:"required"`
	PhoneNumber    string  `json:"phone_number" validate:"required"`
	Address        string  `json:"address" validate:"required"`
	Comment        string  `json:"comment" validate:"omitempty"`
	RewardInternet float64 `json:"reward_internet" validate:"omitempty"`
	RewardCleaning float64 `json:"reward_cleaning" validate:"omitempty"`
	RewardShipping float64 `json:"reward_shipping" validate:"omitempty"`
	IsInternet     bool    `json:"is_internet"`
	IsShipping     bool    `json:"is_shipping"`
	IsCleaning     bool    `json:"is_cleaning" validate:"atLeastOneService"`
}

type LeadFilterDTO struct {
	StatusID   *int64     `json:"status_id"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	UserID     *int64     `json:"user_id"`
	Limit      int64      `json:"limit"`
	Offset     int64      `json:"offset"`
	IsInternet *bool      `json:"is_internet"`
	IsShipping *bool      `json:"is_shipping"`
	IsCleaning *bool      `json:"is_cleaning"`
	Search     *string    `json:"search"`
}

type UserStatistic struct {
	Internet  float64 `json:"internet"`
	Cleaning  float64 `json:"cleaning"`
	Shipping  float64 `json:"shipping"`
	Referrals float64 `json:"referrals"`
	Total     float64 `json:"total"`
}
