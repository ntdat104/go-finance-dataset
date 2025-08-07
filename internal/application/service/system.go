package service

import (
	"github.com/ntdat104/go-finance-dataset/internal/application/dto"
	"github.com/ntdat104/go-finance-dataset/pkg/datetime"
)

type SystemSvc interface {
	GetTime() *dto.SystemTime
}

type systemSvc struct{}

func NewSystemSvc() SystemSvc {
	return &systemSvc{}
}

func (s *systemSvc) GetTime() *dto.SystemTime {
	return &dto.SystemTime{
		ServerTime: datetime.GetCurrentMiliseconds(),
	}
}
