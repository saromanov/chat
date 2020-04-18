package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/saromanov/chat/pkg/models"
)

// AddUser provides adding of the new user
// @Summary Retrieves user based on given ID
// @Produce json
// @Success 201
// @Router /users [post]
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("func", "AddUser").Info("registered request for add new user")
	totalRequests.Inc()
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		s.log.WithField("func", "AddUser").WithError(err).Errorf("unable to unmarshal request")
		render.Render(w, r, ErrInvalidRequest(err, 400))
		return
	}
	if err := s.db.AddUser(&models.User{
		Email:     data.Email,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}); err != nil {
		s.log.WithField("func", "AddUser").WithError(err).Errorf("unable to add user")
		render.Render(w, r, ErrInvalidRequest(err, 400))
		return
	}
	render.Status(r, http.StatusCreated)
}

// GetUser godoc
// @Summary Retrieves user based on given ID
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("func", "AddUser").Info("get user by id")
	user := r.Context().Value("user").(*models.User)
	if err := render.Render(w, r, &UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}); err != nil {
		writeStatusCode(http.StatusInternalServerError, "GetUser")
		render.Render(w, r, ErrInvalidRequest(err, 500))
		return
	}
	writeStatusCode(http.StatusOK, "GetUser")
	render.Status(r, http.StatusOK)
}

// UsersCtx provides handling of context for users endpoints
func (s *Server) UsersCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		if userID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("userid is not found on request"), 400))
			return
		}
		id, err := strconv.ParseInt(userID, 10, 32)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err, 400))
			return
		}
		user, err := s.db.GetUser(id)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err, 404))
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
