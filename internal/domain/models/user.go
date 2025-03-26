package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64
	PhoneNumber    string
	Email          string
	Name           string
	Telegram       string
	City           string
	PasswordHash   string
	ReferralCode   string
	CreatedAt      time.Time
	Role           string
	IsActive       bool
	RewardInternet sql.NullFloat64
	RewardCleaning sql.NullFloat64
	RewardShipping sql.NullFloat64
}
