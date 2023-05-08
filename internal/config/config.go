package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/partyscript/bledger/internal/model"
)

// GlobalConfig is a generic config used to fetch env attributes
type GlobalConfig struct {
	Cache       *CacheConfig
	Environment *EnvironmentConfig
	DB          *DBConfig
}

// EnvironmentConfig is a config to get the environment
type EnvironmentConfig struct {
	Env  model.ApplicationEnvironment `envconfig:"ENVIRONMENT"`
	Port string                       `envconfig:"PORT"`
}

// DBConfig is a config to interact with a generic database
type DBConfig struct {
	DSN string `envconfig:"DB_DSN"`
}

// CacheConfig is a config for cache integrations
type CacheConfig struct {
	URI      string `envconfig:"CACHE_URI"`
	Password string `envconfig:"CACHE_PASSWORD"`
}

// NewGlobalConfig generates a new instance of GlobalConfig
func NewGlobalConfig() (*GlobalConfig, error) {
	var db DBConfig
	var env EnvironmentConfig
	var cache CacheConfig

	err := envconfig.Process("DB", &db)
	if err != nil {
		return nil, errors.New("database config cannot be null")
	}

	err = envconfig.Process("ENV", &env)
	if err != nil {
		return nil, errors.New("environment config cannot be nil")
	}

	err = envconfig.Process("CACHE", &cache)
	if err != nil {
		return nil, errors.New("redis config cannot be nil")
	}

	return &GlobalConfig{
		Environment: &env,
		DB:          &db,
		Cache:       &cache,
	}, nil
}
