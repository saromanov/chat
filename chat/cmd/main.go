package main

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/saromanov/experiments/chat/pkg/server"
	"github.com/saromanov/experiments/chat/pkg/storage"
	"github.com/saromanov/experiments/chat/pkg/config"
)

func makeLogger()*logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
	return l  
}
func main() {
	logger := makeLogger()
	p := &config.Project{
		Log: logger,
		Server: &config.Server {
			Address: ":3000",
		},
		DatabaseHost: "postgres",
		DatabaseUser: "chatapp",
		DatabasePassword: os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName: "chatapp",
	}
	st, err := storage.New(p)
	if err != nil {
		logger.WithError(err).Fatalf("unable to load storage")
	}
	server.Make(st, p)
}