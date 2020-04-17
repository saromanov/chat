package config

import (
	"github.com/saromanov/cowrow"
	"github.com/sirupsen/logrus"
)

// Project defines root configuration
type Project struct {
	Log              *logrus.Logger
	Server           *Server `yaml:"server"`
	DatabaseName     string  `yaml:"databaseName"`
	DatabaseUser     string  `yaml:"databaseUser"`
	DatabasePassword string  `yaml:"databasePassword"`
	DatabaseHost     string  `yaml:"databaseHost"`
	Address          string  `yaml:"address"`
}

// LoadConfig provides loading of the config from path
func LoadConfig(path string) (*Project, error) {
	cfg := &Project{}
	err := cowrow.LoadByPath(path, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Server provides definition for server config
type Server struct {
	Address string `yaml:"address"`
}
