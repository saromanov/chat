package config

import (
	"github.com/sirupsen/logrus"
)

// Project defines root configuration
type Project struct {
	Log *logrus.Logger
	Server *Server
	DatabaseName string
	DatabaseUser string
	DatabasePassword string
	DatabaseHost string
	
}

// Server provides definition for server config
type Server struct {
	Address string
}
