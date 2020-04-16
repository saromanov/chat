package models

// User defines User model
type User struct {
	Model
	FirstName string
	LastName string
	Email string
	Password string
}