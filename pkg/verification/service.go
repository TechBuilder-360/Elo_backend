package verification

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/provider"
)

type IService interface {
	RequestVerificationLink(ctx context.Context, payload *VerificationRequest) (_ *VerificationResponse, err error)
	ProcessVerification(ctx context.Context, payload *VerificationResult) error
	GetProviderLink(ctx context.Context, reference string) (string, error)
}

type Service struct {
	db             *ent.Client
	serviceLocator provider.IService
}

func NewService() IService {
	return &Service{
		db:             database.DBInstance(),
		serviceLocator: provider.NewService(),
	}
}
