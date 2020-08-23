package rest

import (
	"crypto/rand"
	"io"

	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

var tokenSecret []byte

// SetTokenSecret sets the secret of the token used to encrypt jwt.
func SetTokenSecret(secret string) {
	tokenSecret = []byte(secret)
}

func getTokenSecret() []byte {
	if len(tokenSecret) == 0 {
		tokenSecret = make([]byte, 50)
		_, err := io.ReadFull(rand.Reader, tokenSecret)
		if err != nil {
			common.Logger().Fatal("token secret undefined")
		}
	}
	return tokenSecret
}
