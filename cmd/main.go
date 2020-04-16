package main

import (
	"os"

	"github.com/saromanov/chat/pkg/config"
	"github.com/saromanov/chat/pkg/server"
	"github.com/saromanov/chat/pkg/storage"
	"github.com/sirupsen/logrus"
)

func makeLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
	return l
}
func main() {
	logger := makeLogger()
	p := &config.Project{
		Log: logger,
		Server: &config.Server{
			Address: ":3005",
		},
		DatabaseHost:     os.Getenv("POSTGRES_HOST"),
		DatabaseUser:     os.Getenv("POSTGRES_USER"),
		DatabasePassword: os.Getenv("POSTGRES_PASSWORD"),
		DatabaseName:     os.Getenv("POSTGRES_DB"),
	}
	st, err := storage.New(p)
	if err != nil {
		logger.WithError(err).Fatalf("unable to load storage")
	}
	server.Make(st, p)
}
