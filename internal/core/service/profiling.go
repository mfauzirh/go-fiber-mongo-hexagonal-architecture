package service

import (
	"context"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/port"
)

type ProfilingService struct {
	profilingRepository port.ProfilingRepository
}

func NewProfilingService(profilingRepository port.ProfilingRepository) port.ProfilingService {
	return &ProfilingService{
		profilingRepository: profilingRepository,
	}
}

func (s *ProfilingService) InsertProfilingData(ctx context.Context, data *domain.Profiling) (*domain.Profiling, error) {
	return s.profilingRepository.InsertProfilingData(ctx, data)
}
