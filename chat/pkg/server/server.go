package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ory/graceful"
	"github.com/saromanov/experiments/chat/pkg/config"
)

// Make provides making of server
func Make(p *config.Project) {
	
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	server := graceful.WithDefaults(&http.Server{
        Addr: ":3000",
        Handler: r,
	})
	fmt.Println("STARTD")
	p.Log.WithField("level", "server").Info("main: Starting the server")
    if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
        p.Log.WithField("level", "server").Fatalln("main: Failed to gracefully shutdown")
    }
	p.Log.WithField("level", "server").Info("main: Server was shutdown gracefully")
}
