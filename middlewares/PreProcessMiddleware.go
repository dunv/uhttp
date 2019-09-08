package middlewares

import (
	"net/http"

	"github.com/dunv/uhttp/helpers"
	"github.com/dunv/uhttp/models"
)

func PreProcess(handler models.Handler) models.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if handler.PreProcess != nil {
				err := (handler.PreProcess)(r.Context())
				if err != nil {
					helpers.RenderError(w, r, err)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
	}
}
