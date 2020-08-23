package config

import (
	"os"

	match "github.com/go-ozzo/ozzo-validation/v4/is"

	valid "github.com/go-ozzo/ozzo-validation/v4"
)

// FromEnv returns a new configuration from env vars.
func FromEnv() (Config, error) {
	cfg := Config{
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		JwtSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}

	err := valid.ValidateStruct(
		&cfg,
		valid.Field(&cfg.DBHost, valid.Required, match.Host),
		valid.Field(&cfg.DBPort, match.Digit),
		valid.Field(&cfg.DBName, valid.Required),
		valid.Field(&cfg.DBUser, valid.Required),
		valid.Field(&cfg.Port, match.Digit),
	)

	if cfg.Port == "" {
		cfg.Port = "80" // Default listen port.
	}

	return cfg, err
}
