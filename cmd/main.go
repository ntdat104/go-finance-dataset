package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
	"github.com/ntdat104/go-finance-dataset/internal/interfaces"
	"github.com/ntdat104/go-finance-dataset/pkg/config"
	"github.com/ntdat104/go-finance-dataset/pkg/logger"
	"github.com/ntdat104/go-finance-dataset/pkg/middleware"
)

type Meta struct {
	MessageID string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
	Datetime  string `json:"datetime"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data,omitempty"`
}

func main() {
	// Initialize cfg
	cfg := config.NewConfig("./config/dev.yml")

	logg := logger.InitLogger("./log", cfg.App.Name, cfg.App.Version)
	defer logg.Sync()

	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production
	router := gin.New()          // Clean router (no default logger)
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ZapLoggerWithBody(logg))

	baseService := service.NewBaseService()
	systemService := service.NewSystemService()
	interfaces.NewSystemController(router, baseService, systemService)

	// Run server
	port := strconv.Itoa(cfg.HTTP.Port)
	log.Printf("%v started on http://%v:%v", cfg.App.Name, cfg.HTTP.Host, port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf(cfg.App.Name+" failed to start: %v", err)
	}
}
