package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/constants"
	"github.com/ntdat104/go-finance-dataset/internal/application/response"
	"github.com/ntdat104/go-finance-dataset/internal/application/service"
)

type SystemHandler interface {
	GetTime(ctx *gin.Context)
}

type systemHandler struct {
	router    *gin.Engine
	systemSvc service.SystemSvc
}

func NewSystemHandler(router *gin.Engine, systemSvc service.SystemSvc) SystemHandler {
	h := &systemHandler{
		router:    router,
		systemSvc: systemSvc,
	}
	h.initRoutes()
	return h
}

func (h *systemHandler) initRoutes() {
	h.router.GET(constants.ApiSystemTime, h.GetTime)
}

func (h *systemHandler) GetTime(ctx *gin.Context) {
	response.Success(ctx, h.systemSvc.GetTime())
}
