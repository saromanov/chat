package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/ory/graceful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saromanov/chat/pkg/config"
	"github.com/saromanov/chat/pkg/storage"
	"github.com/sirupsen/logrus"
)

var tokenAuth *jwtauth.JWTAuth

type Server struct {
	db  *storage.Storage
	log *logrus.Logger
}

// Make provides making of server
func Make(st *storage.Storage, p *config.Project) {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	s := &Server{
		db: st,
	}
	initPrometheus()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		totalRequests.Inc()
		w.Write([]byte("welcome"))
	})
	r.Route("/users", func(r chi.Router) {
		totalRequests.Inc()
		r.Post("/register", s.AddUser)
	})

	r.Route("/messages", func(r chi.Router) {
		r.Post("/", s.AddMessage)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/users/{id}", s.GetUser)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Use(s.UsersCtx)
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
