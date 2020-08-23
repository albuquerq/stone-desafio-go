package psql

import (
	"testing"

	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/config"
	"github.com/stretchr/testify/assert"
)

func TestURLFromConfig(t *testing.T) {
	cases := []struct {
		In  config.Config
		Out string
	}{
		{
			In: config.Config{
				DBHost: "localhost",
				DBPort: "1035",
				DBPass: "password",
				DBUser: "user",
				DBName: "database-test",
			},
			Out: "postgres://user:password@localhost:1035/database-test",
		},
		{

			In: config.Config{
				DBHost: "localhost",
				DBPort: "",
				DBPass: "password",
				DBUser: "user",
				DBName: "database-test",
			},
			Out: "postgres://user:password@localhost/database-test",
		},
		{

			In: config.Config{
				DBHost: "localhost",
				DBPort: "",
				DBPass: "",
				DBUser: "",
				DBName: "database-test",
			},
			Out: "postgres://@localhost/database-test",
		},
	}

	for _, c := range cases {
		url := URLFromConfig(c.In)
		assert.Equal(t, c.Out, url)
	}
}
