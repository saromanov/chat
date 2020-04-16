package server

import (
	"net/http"

	"github.com/go-chi/render"
)

// AddMessage provides adding of the new message
// @Summary Retrieves user based on given ID
// @Produce json
// @Success 201
// @Router /messages [post]
func (s *Server) AddMessage(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("func", "AddMessage").Info("registered request for add new message")
	totalRequests.Inc()
	render.Status(r, http.StatusCreated)
}
