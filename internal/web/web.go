package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ssnakerss/TypingHero/internal/models"
	"github.com/Ssnakerss/TypingHero/internal/web/handlers"
	"github.com/Ssnakerss/TypingHero/internal/web/router"
)

func StartWeb(ctx context.Context, db models.Storage) {
	hm := handlers.NewHandlerMaster(db)
	r := router.New(hm)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Typing Trainer server running at http://localhost:8080")
	go server.ListenAndServe()
	<-ctx.Done()
	server.Shutdown(ctx)

}
