package models

import "time"

// Model defines basic model for the app
type Model struct {
	ID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}