package config

import (
	"github.com/sirupsen/logrus"
)

// Project defines root configuration
type Project struct {
	Log *logrus.Logger
}
