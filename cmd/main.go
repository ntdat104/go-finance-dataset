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
	config.InitConfig("./config/dev.yml")
	logger.InitProduction("./log")
	defer logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ZapLoggerWithBody())

	systemSvc := service.NewSystemSvc()
	interfaces.NewSystemHandler(router, systemSvc)

	binanceSvc := service.NewBinanceSvc()
	interfaces.NewBinanceHandler(router, binanceSvc)

	cfg := config.GetGlobalConfig()
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("%v started on http://%v:%v", cfg.App.Name, cfg.HTTP.Host, strconv.Itoa(cfg.HTTP.Port))
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
