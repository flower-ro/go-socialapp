package config

import (
	"go-socialapp/internal/socialserver/options"
)

// Config is the running configuration structure of the IAM pump service.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM pump command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {

	config := &Config{Options: opts}

	return config, nil
}
