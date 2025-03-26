package models

import "time"

type Referral struct {
	ID           int64
	UserID       int64
	ReferralCode string
	CreatedAt    time.Time
}
