package models

import "time"

type PasswordCode struct {
	ID        int
	UserID    int64
	Code      string
	ExpiresAt time.Time
}
