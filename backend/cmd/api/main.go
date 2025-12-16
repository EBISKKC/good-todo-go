package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"good-todo-go/internal/infrastructure/database"
	"good-todo-go/internal/infrastructure/environment"
	"good-todo-go/internal/presentation/public/router/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	cfg, err := environment.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	entClient, err := database.NewEntClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseEntClient(entClient)

	// Initialize Echo
	e := echo.New()

	// Setup middleware
	middleware.SetupMiddleware(e)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// TODO: Setup routes with dependency injection

	// Start server
	go func() {
		addr := ":" + cfg.Port
		log.Printf("Starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
