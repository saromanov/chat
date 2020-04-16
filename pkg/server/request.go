package server

import (
	"net/http"
	"strings"
)

type UserRequest struct {
	Email string
	Password string
	FirstName string
	LastName string
}


func (a *UserRequest) Bind(r *http.Request) error {
	a.Email = strings.ToLower(a.Email)
	return nil
}