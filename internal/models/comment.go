package models

import "time"

type Comment struct {
	ID        int64
	LeadID    int64
	UserID    int64
	Text      string
	CreatedAt *time.Time
}
