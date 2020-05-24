package main

import (
	"os"

	"github.com/saromanov/chat/pkg/config"
	"github.com/saromanov/chat/pkg/server"
	"github.com/saromanov/chat/pkg/storage"
	logrusloki "github.com/schoentoon/logrus-loki"
	"github.com/sirupsen/logrus"
)

func makeLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
	loki, err := logrusloki.NewLoki("http://loki:3100/api/prom/push", 10, 1)
	if err != nil {
		panic("unable to init loki")
	}
	l.AddHook(loki)
	return l
}

func main() {
	logger := makeLogger()
	p, err := config.LoadConfig("./configs/app/app.yaml")
	if err != nil {
		logger.WithError(err).Fatalf("unable to load config")
	}
	if p == nil {
		logger.Fatalf("unable to load config")
	}
	p.DatabasePassword = os.Getenv("POSTGRES_PASSWORD")
	p.DatabaseHost = os.Getenv("POSTGRES_HOST")
	p.Log = logger
	st, err := storage.New(p)
	if err != nil {
		logger.WithError(err).Fatalf("unable to load storage")
	}
	if err := st.Prepare(); err != nil {
		logger.WithError(err).Fatalf("unable to prepare storage")
	}
	server.Make(st, p)
}
