package models

import "time"

type Lead struct {
	ID          int64    `json:"id"`
	UserID      int64    `json:"user_id"`
	FIO         string   `json:"fio"`
	Address     string   `json:"address"`
	StatusID    int64    `json:"status_id"`
	PhoneNumber string   `json:"phone_number"`
	Internet    bool     `json:"is_internet"`
	Cleaning    bool     `json:"is_cleaning"`
	Shipping    bool     `json:"is_shipping"`
	Comments    []string `json:"comments"`

	RewardInternet float64 `json:"reward_internet"`
	RewardCleaning float64 `json:"reward_cleaning"`
	RewardShipping float64 `json:"reward_shipping"`

	CreatedAt   *time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	PaymentAt   *time.Time `json:"payment_at"`
}
