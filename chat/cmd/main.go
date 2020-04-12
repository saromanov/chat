package main

import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/saromanov/experiments/chat/pkg/server"
	"github.com/saromanov/experiments/chat/pkg/config"
)

func makeLogger()*logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.SetFormatter(&logrus.JSONFormatter{})
	return l  
}
func main() {
	p := &config.Project{
		Log: makeLogger(),
	}
	server.Make(p)
}