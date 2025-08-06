package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/constants"
	"github.com/ntdat104/go-finance-dataset/internal/application/response"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
)

type SystemController struct {
	router        *gin.Engine
	systemService *service.SystemService
}

func NewSystemController(router *gin.Engine, systemService *service.SystemService) *SystemController {
	systemController := SystemController{
		router:        router,
		systemService: systemService,
	}

	router.GET(constants.ApiSystemTime, systemController.GetTime)

	return &systemController
}

func (c SystemController) GetTime(ctx *gin.Context) {
	response.Success(ctx, c.systemService.GetTime())
}
