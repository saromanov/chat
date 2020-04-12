package models

import "time"

// Message defines messaging to user
type Message struct {
	ID int64
	Text string
	CreatedAt time.Time
	UpdateAt time.Time
}
