package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID           int64
	PhoneNumber  string
	Email        string
	Name         string
	Telegram     string
	City         string
	PasswordHash string
	ReferralCode string
	CreatedAt    time.Time
	Roles        pq.StringArray
	IsActive     bool
}
