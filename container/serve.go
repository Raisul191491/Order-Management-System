package container

import (
	"context"
	"errors"
	"log"
	"net/http"
	"oms/config"
	"oms/routes"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(e *gin.Engine) {
	logger := log.New(os.Stdout, "[OMS] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting OMS server initialization...")

	// Load configuration
	logger.Println("Loading configuration...")
	cfg := config.LoadConfig()
	if cfg == nil {
		logger.Fatal("Failed to load configuration")
	}

	// Initialize routes
	logger.Println("Initializing routes...")
	routes.InitRoutes(e)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: e,
	}

	go func() {
		logger.Printf("Server starting on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	logger.Printf("Server is running on http://localhost:%s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Printf("Received signal: %s", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	logger.Println("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("Server shutdown error: %v", err)
	} else {
		logger.Println("HTTP server stopped gracefully")
	}

	// Wait a moment for any pending operations
	select {
	case <-shutdownCtx.Done():
		if errors.Is(context.DeadlineExceeded, shutdownCtx.Err()) {
			logger.Println("Shutdown completed with timeout")
		}
	default:
	}

	logger.Println("OMS server terminated")
}
