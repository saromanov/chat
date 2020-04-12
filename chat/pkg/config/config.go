package config

import (
	"github.com/sirupsen/logrus"
)

// Project defines root configuration
type Project struct {
	Log *logrus.Logger
	Server *Server
}

// Server provides definition for server config
type Server struct {
	Address string
}
