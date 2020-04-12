package main

import (
	"github.com/sirupsen/logrus"
	"github.com/saromanov/experiments/chat/pkg/server"
	"github.com/saromanov/experiments/chat/pkg/config"
)

func makeLogger()*logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(l.Writer())
	return l  
}
func main() {
	p := &config.Project{
		Log: makeLogger(),
	}
	server.Make(p)
}