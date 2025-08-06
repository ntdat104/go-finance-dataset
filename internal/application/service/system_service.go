package service

import (
	"github.com/ntdat104/go-finance-dataset/internal/application/dto"
	"github.com/ntdat104/go-finance-dataset/pkg/datetime"
)

type SystemService struct {
}

func NewSystemService() *SystemService {
	return &SystemService{}
}

func (s *SystemService) GetTime() dto.SystemTime {
	return dto.SystemTime{
		ServerTime: datetime.GetCurrentMiliseconds(),
	}
}
