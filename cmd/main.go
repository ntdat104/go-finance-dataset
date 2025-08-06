package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
	"github.com/ntdat104/go-finance-dataset/internal/interfaces"
	"github.com/ntdat104/go-finance-dataset/pkg/config"
	"github.com/ntdat104/go-finance-dataset/pkg/logger"
	"github.com/ntdat104/go-finance-dataset/pkg/middleware"
)

func main() {
	// Initialize cfg
	cfg := config.NewConfig("./config/dev.yml")

	logg := logger.InitLogger("./log", cfg.App.Name, cfg.App.Version)
	defer logg.Sync()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ZapLoggerWithBody(logg))

	systemService := service.NewSystemService()
	interfaces.NewSystemController(router, systemService)

	port := strconv.Itoa(cfg.HTTP.Port)
	addr := ":" + port

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("%v started on http://%v:%v", cfg.App.Name, cfg.HTTP.Host, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(cfg.App.Name+" failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
