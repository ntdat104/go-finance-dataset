package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
)

type SystemController struct {
	router        *gin.Engine
	baseService   *service.BaseService
	systemService *service.SystemService
}

func NewSystemController(router *gin.Engine, baseService *service.BaseService, systemService *service.SystemService) *SystemController {
	systemController := SystemController{
		router:        router,
		baseService:   baseService,
		systemService: systemService,
	}

	apiGroup := router.Group("/api/v1/system")
	{
		apiGroup.GET("/time", systemController.GetTime)
	}

	return &systemController
}

func (systemController SystemController) GetTime(ctx *gin.Context) {
	systemController.baseService.NewResponse(ctx, systemController.systemService.GetTime())
}
