package main

import (
	"log"
	"net/http"

	"good-todo-go/internal/infrastructure/database"
	"good-todo-go/internal/presentation/public/router"
)

func main() {
	// NewRouterでEcho、Config、EntClientを取得
	e, cfg, entClient, err := router.NewRouter()
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}
	defer database.CloseEntClient(entClient)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
