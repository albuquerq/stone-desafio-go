package psql

import (
	"context"
	"fmt"

	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/config"
)

// getContext generator.
func getContext() context.Context {
	return context.Background()
}

// URLFromConfig returns a database connection URL from the application's configuration.
func URLFromConfig(cfg config.Config) string {

	host := cfg.DBHost
	if cfg.DBPort != "" {
		host += ":" + cfg.DBPort
	}

	cred := cfg.DBUser
	if cfg.DBPass != "" {
		cred += ":" + cfg.DBPass
	}

	return fmt.Sprintf("postgres://%s@%s/%s", cred, host, cfg.DBName)
}
