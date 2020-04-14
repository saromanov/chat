package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/middleware"
	"github.com/ory/graceful"
	"github.com/sirupsen/logrus"
	"github.com/go-chi/jwtauth"
	"github.com/saromanov/experiments/chat/pkg/config"
	"github.com/saromanov/experiments/chat/pkg/models"
	"github.com/saromanov/experiments/chat/pkg/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var tokenAuth *jwtauth.JWTAuth

type Server struct {
	db *storage.Storage
	log *logrus.Logger
}

// AddUser provides adding of the new user
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {
	s.log.WithField("func", "AddUser").Info("registered request for add new user")
	totalRequests.Inc()
	data := &UserRequest{}
	if err := render.Bind(r, data); err != nil {
		s.log.WithField("func", "AddUser").WithError(err).Errorf("unable to unmarshal request")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if err := s.db.AddUser(&models.User{
		Email: data.Email,
		FirstName: data.FirstName,
		LastName: data.LastName,
	}); err != nil {
		s.log.WithField("func", "AddUser").WithError(err).Errorf("unable to add user")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
}

// Make provides making of server
func Make(st *storage.Storage, p *config.Project) {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	s := &Server {
		db: st,

	}
	fmt.Println("ADDERSSERVER: ", p.Server.Address)
	initPrometheus()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		totalRequests.Inc()
		w.Write([]byte("welcome"))
	})
	r.Route("/users", func(r chi.Router){
		r.Post("/register", s.AddUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
	})
	r.Handle("/metrics", promhttp.Handler())
	server := graceful.WithDefaults(&http.Server{
        Addr: p.Server.Address,
        Handler: r,
	})
	p.Log.WithField("package", "server").Info("main: Starting the server")
    if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
        p.Log.WithField("package", "server").Fatalln("main: Failed to gracefully shutdown")
    }
	p.Log.WithField("package", "server").Info("main: Server was shutdown gracefully")
}
