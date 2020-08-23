package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/albuquerq/stone-desafio-go/pkg/domain"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/common"
)

// Handler http rest handler.
type Handler struct {
	handler  http.Handler
	registry domain.Registry
}

// New returns a new rest handler.
func New(registry domain.Registry) http.Handler {
	h := Handler{
		registry: registry,
	}

	h.DefineRoutes() // Create de mux handler and define api routes.

	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.handler == nil {
		common.Logger().Fatal("http handler not defined")
	}
	h.handler.ServeHTTP(w, r)
}

// WrapControl returns handler functions from custom API controls.
func (h *Handler) WrapControl(fn func(w http.ResponseWriter, r *http.Request) Response) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1000000) // Limit body size.

		resp := fn(w, r)

		w.Header().Set("content-type", "application/json")

		if resp.Error != nil {
			w.WriteHeader(resp.Error.Code)
		}

		err := json.NewEncoder(w).Encode(resp)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Error: &Error{
					Code:    http.StatusInternalServerError,
					Message: errors.New("error on encode JSON data").Error(),
				},
			})
		}
		r.Body.Close() // connection not persistent.
	}
}
