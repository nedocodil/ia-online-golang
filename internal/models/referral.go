package models

import "time"

type Referral struct {
	ID           int64
	UserID       int64
	ReferralCode string
	Cost         float64
	CreatedAt    time.Time
	Active       bool
}

type ReferralAndUser struct {
	Referral Referral
	User     User
}
