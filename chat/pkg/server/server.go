package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ory/graceful"
	"github.com/saromanov/experiments/chat/pkg/config"
	"github.com/saromanov/experiments/chat/pkg/storage"
)

type Server struct {
	db *storage.Storage
}

// AddUser provides adding of the new user
func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {

}

// Make provides making of server
func Make(p *config.Server) {
	s := &Server {
		db: p.db,
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) 
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Route("/users", func(r chi.Route){
		r.Post("/register", s.AddUser)
	})
	server := graceful.WithDefaults(&http.Server{
        Addr: p.Address,
        Handler: r,
	})
	p.Log.WithField("package", "server").Info("main: Starting the server")
    if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
        p.Log.WithField("package", "server").Fatalln("main: Failed to gracefully shutdown")
    }
	p.Log.WithField("package", "server").Info("main: Server was shutdown gracefully")
}
