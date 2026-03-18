package router

import (
	"github.com/Ssnakerss/TypingHero/internal/web/handlers"
	"github.com/go-chi/chi"
)

func New(hm *handlers.HandlerMaster) *chi.Mux {
	r := chi.NewRouter()
	// r.Get("/", hm.HomeHandler)
	// r.Get("/", http.FileServer(http.Dir("./html")))
	r.Handle("/", hm.HomeHandler)

	r.Get("/api/text", hm.GetTextHandler)
	r.Get("/api/result", hm.CalculateResultHandler)
	return r
}
