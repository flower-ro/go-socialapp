package options

import (
	"go-socialapp/pkg/db"
	"time"

	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

// MySQLOptions defines options for mysql database.
type DBOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Port                  string        `json:"port,omitempty"                     mapstructure:"port"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

// NewMySQLOptions create a `zero` value instance.
func NewDBOptions() *DBOptions {
	return &DBOptions{
		Host:                  "127.0.0.1",
		Port:                  "3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *DBOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *DBOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "db.host", o.Host, ""+
		"db service host address. If left blank, the following related mysql options will be ignored.")
	fs.StringVar(&o.Host, "db.port", o.Port, ""+
		"db service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Username, "db.username", o.Username, ""+
		"Username for access to db service.")

	fs.StringVar(&o.Password, "db.password", o.Password, ""+
		"Password for access to db, should be used pair with password.")

	fs.StringVar(&o.Database, "db.database", o.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&o.MaxIdleConnections, "db.max-idle-connections", o.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to db.")

	fs.IntVar(&o.MaxOpenConnections, "db.max-open-connections", o.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to db.")

	fs.DurationVar(&o.MaxConnectionLifeTime, "db.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to db.")

	fs.IntVar(&o.LogLevel, "db.log-mode", o.LogLevel, ""+
		"Specify gorm log level.")
}

// NewPGClient create mysql store with the given config.
func (o *DBOptions) NewPGClient() (*gorm.DB, error) {
	opts := &db.Options{
		Host:                  o.Host,
		Username:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdleConnections:    o.MaxIdleConnections,
		MaxOpenConnections:    o.MaxOpenConnections,
		MaxConnectionLifeTime: o.MaxConnectionLifeTime,
		LogLevel:              o.LogLevel,
		//Logger:  未初始化
	}

	return db.NewPG(opts)
}

func (o *DBOptions) NewMysqlClient() (*gorm.DB, error) {
	opts := &db.Options{
		Host:                  o.Host,
		Username:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdleConnections:    o.MaxIdleConnections,
		MaxOpenConnections:    o.MaxOpenConnections,
		MaxConnectionLifeTime: o.MaxConnectionLifeTime,
		LogLevel:              o.LogLevel,
		//Logger:  未初始化
	}

	return db.New(opts)
}
