package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ruscigno/ticker-signals/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StartHttpServer starts the HTTP server using gin package and the routes defined in the handlers package
func StartHttpServer(ctx context.Context, cfg *config.AppConfig, start time.Time) {
	router := gin.New()
	// Register logger middleware.
	router.Use(Logger(), Recovery())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register HTTP route handlers.
	registerRoutes(router, cfg)

	// Create new HTTP server instance.
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HttpHost(), cfg.HttpPort()),
		Handler: router,
	}

	zap.L().Debug(fmt.Sprintf("server: successfully initialized [%s]", time.Since(start)))

	// Start HTTP server.
	go func() {
		zap.L().Info(fmt.Sprintf("server: listening at %s", server.Addr))

		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				zap.L().Info("server: shutdown complete")
			} else {
				zap.L().Error("server error", zap.Error(err))
			}
		}
	}()

	// Graceful HTTP server shutdown.
	<-ctx.Done()
	zap.L().Info("server: shutting down")
	err := server.Close()
	if err != nil {
		zap.L().Error("server: shutdown failed", zap.Error(err))
	}
}
