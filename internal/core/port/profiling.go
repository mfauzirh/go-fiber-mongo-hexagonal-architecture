package port

import (
	"context"

	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/core/domain"
)

type ProfilingRepository interface {
	InsertProfilingData(ctx context.Context, data *domain.Profiling) (*domain.Profiling, error)
}

type ProfilingService interface {
	InsertProfilingData(ctx context.Context, data *domain.Profiling) (*domain.Profiling, error)
}
