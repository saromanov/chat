package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/ory/graceful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saromanov/experiments/chat/pkg/config"
	"github.com/saromanov/experiments/chat/pkg/models"
	"github.com/saromanov/experiments/chat/pkg/storage"
	"github.com/sirupsen/logrus"
)

var tokenAuth *jwtauth.JWTAuth

type Server struct {
	db  *storage.Storage
	log *logrus.Logger
}

// AddUser provides adding of the new user
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
		render.Render(w, r, ErrInvalidRequest(err, 400))
		return
	}
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
		user, err := s.db.GetUserByID(id)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err, 404))
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Make provides making of server
func Make(st *storage.Storage, p *config.Project) {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	s := &Server{
		db: st,
	}
	initPrometheus()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		totalRequests.Inc()
		w.Write([]byte("welcome"))
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", s.AddUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/users/{id}", s.GetUser)
	})
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(s.UsersCtx)
		r.Get("/users/{id}", s.GetUser)
	})

	r.Handle("/metrics", promhttp.Handler())
	server := graceful.WithDefaults(&http.Server{
		Addr:    p.Server.Address,
		Handler: r,
	})
	p.Log.WithField("package", "server").Info("main: Starting the server")
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		p.Log.WithField("package", "server").Fatalln("main: Failed to gracefully shutdown")
	}
	p.Log.WithField("package", "server").Info("main: Server was shutdown gracefully")
}
