package db

import (
	"gorm.io/gorm/logger"
	"time"
)

// Options defines optsions for mysql database.
type Options struct {
	Host                  string
	Port                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}
