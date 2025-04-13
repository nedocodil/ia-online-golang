package models

import "time"

type ActivationLink struct {
	ID           int
	UserID       int64
	ActivationID string
	ExpiresAt    time.Time
}
