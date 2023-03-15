package config

import "github.com/ZakirAvrora/TechHome/pkg/env"

type Config struct {
	Environment string
	Database    *DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

func NewConfig(path string) *Config {
	env.CheckDotEnv(path)
	return &Config{
		Environment: env.MustGet("ENV"),
		Database: &DatabaseConfig{
			Host:     env.MustGet("DATABASE_HOST"),
			Port:     env.MustGet("DATABASE_PORT"),
			User:     env.MustGet("DATABASE_USER"),
			DbName:   env.MustGet("DATABASE_DB"),
			Password: env.MustGet("DATABASE_PASSWORD"),
		},
	}
}
