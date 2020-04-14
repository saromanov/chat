package main

import (
	"github.com/saromanov/experiments/chat/pkg/config"
	"github.com/saromanov/experiments/chat/pkg/server"
	"github.com/saromanov/experiments/chat/pkg/storage"
	"github.com/sirupsen/logrus"
	"os"
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
			Address: ":3000",
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
