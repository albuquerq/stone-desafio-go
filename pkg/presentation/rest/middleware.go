package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/presentation/rest/contextkey"
)

// current methods allowed in rest api.
var allowedMethods = strings.Join([]string{
	http.MethodOptions,
	http.MethodGet,
	http.MethodPost,
}, ", ")

var allowedHeaders = strings.Join([]string{
	"Content-Type",
	"Authorization",
}, ", ")

// CORS middleware for REST API handler.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Add("Access-Control-Allow-Methods", allowedMethods)

		if r.Method == http.MethodOptions {
			if rh := r.Header.Get("Access-Control-Request-Headers"); rh != "" {
				w.Header().Add("Access-Control-Allow-Headers", rh) // In options accept all headers.
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AccountAccessCtxMiddleware obtain encrypted data from the jwt token and add to the request context.
func AccountAccessCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if ah == "" {
			http.Error(w, "missing authorization header with access token", http.StatusUnauthorized)
			return
		}

		tokenparts := strings.Fields(ah)

		if len(tokenparts) != 2 || strings.ToLower(tokenparts[0]) != "bearer" {
			http.Error(w, "missing bearer access token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenparts[1], func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]), http.StatusBadRequest)
			}

			return getTokenSecret(), nil
		})
		if err != nil {
			http.Error(w, "error on parse jwt token", http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			http.Error(w, "token not valid", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			descr := access.Description{}

			descr.AccountID = claims["account_id"].(string)
			descr.CPF = claims["account_cpf"].(string)
			descr.Name = claims["account_name"].(string)

			r = r.WithContext(context.WithValue(r.Context(), contextkey.AccountDescription, descr))

		} else {
			http.Error(w, "missing jwt claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
