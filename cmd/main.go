package main

import (
	"log"
	"os"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/saromanov/chat/pkg/config"
	"github.com/saromanov/chat/pkg/server"
	"github.com/saromanov/chat/pkg/storage"
	"github.com/sirupsen/logrus"
)

func makePromtail(sourceName, jobName string) {
	labels := "{source=\"" + sourceName + "\",job=\"" + jobName + "\"}"
	conf := promtail.ClientConfig{
		PushURL:            "http://localhost:3100/api/prom/push",
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	loki, err := promtail.NewClientJson(conf)
	if err != nil {
		log.Printf("promtail.NewClient: %s\n", err)
		os.Exit(1)
	}

	loki.Debugf("source = %s time = %s\n", sourceName, time.Now().Format(time.RFC3339))
}
func makeLogger() *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
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
	makePromtail("test", "data")
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
