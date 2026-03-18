package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ssnakerss/TypingHero/internal/console"
	"github.com/Ssnakerss/TypingHero/internal/database"
	"github.com/Ssnakerss/TypingHero/internal/web"
)

func main() {
	// Initialize database
	dbPath := "data/typing_sessions.sqlite"

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Warning: Could not initialize database: %v", err)
	}
	defer db.CloseDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		sig := <-exit
		slog.Warn("signal received", "termination", sig)
		slog.Info("stopping server")
		cancel()
	}()

	go console.RunGame(ctx, cancel, db)
	go web.StartWeb(ctx, db)
	<-ctx.Done()
}
