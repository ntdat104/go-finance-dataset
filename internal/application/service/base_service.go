package service

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/response"
	"github.com/ntdat104/go-finance-dataset/pkg/datetime"
	"github.com/ntdat104/go-finance-dataset/pkg/http"
	"github.com/ntdat104/go-finance-dataset/pkg/uuid"
)

type BaseService struct{}

func NewBaseService() *BaseService {
	return &BaseService{}
}

func (bs *BaseService) NewResponse(c *gin.Context, data any) {
	now := datetime.GetCurrentMiliseconds()
	c.JSON(http.StatusOK, response.Response{
		Meta: response.Meta{
			MessageID: uuid.NewShortUUID(),
			Timestamp: now,
			Datetime:  datetime.ConvertMillisecondsToString(now, datetime.YYYY_MM_DD_HH_MM_SS),
			Code:      http.StatusOK,
			Message:   http.StatusText(http.StatusOK),
		},
		Data: data,
	})
}

func (bs *BaseService) NewErrorResponse(c *gin.Context, code int) {
	now := datetime.GetCurrentMiliseconds()
	c.JSON(code, response.Response{
		Meta: response.Meta{
			MessageID: uuid.NewShortUUID(),
			Timestamp: now,
			Datetime:  datetime.ConvertMillisecondsToString(now, datetime.YYYY_MM_DD_HH_MM_SS),
			Code:      code,
			Message:   http.StatusText(code),
		},
	})
}
