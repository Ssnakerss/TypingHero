package router

import (
	"github.com/go-chi/chi"
)

func New(hm *handlers.HandlerMaster) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/text", web.getTextHandler)
	r.Get("/api/result", web.calculateResultHandler)
	return r
}
