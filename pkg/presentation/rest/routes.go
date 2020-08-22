package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Handler) DefineRoutes() {

	mux := chi.NewRouter()

	mux.Route("/api/v1", func(r chi.Router) {

		r.Use(CORS) // Accept CORS for api v1.

		r.Post("/login", h.WrapControl(h.Login))

		r.Post("/accounts", h.WrapControl(h.AccountCreate))
		r.Get("/accounts", h.WrapControl(h.AccountList))
		r.Get("/accounts/{accountID}/balance", h.WrapControl(h.AccountBalance))

		r.Route("/transfers", func(tr chi.Router) {
			tr.Use(AccountAccessCtxMiddleware)
			tr.Get("/", h.WrapControl(h.AccountTransferList))
			tr.Post("/", h.WrapControl(h.TransferCreate))
		})
	})

	h.handler = mux
}

func (h *Handler) getRouteParam(r *http.Request, param string) (val string, exists bool) {
	val = chi.URLParam(r, param)
	if val == "" {
		return "", false
	}
	return val, true
}
