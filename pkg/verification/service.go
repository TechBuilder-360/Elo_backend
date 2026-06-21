package verification

import (
	"context"

	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/providers/cloudinary"
)

type IService interface {
	RequestVerificationLink(ctx context.Context, payload *VerificationRequest) (_ *VerificationResponse, err error)
	ProcessVerification(ctx context.Context, payload *VerificationResult) error
	GetProviderLink(ctx context.Context, reference string) (string, error)
}

type Service struct {
	db             *ent.Client
	serviceLocator provider.IService
	cloudinary     cloudinary.Cloud
}

func NewService() IService {
	return &Service{
		db:             database.DBInstance(),
		serviceLocator: provider.NewService(),
		cloudinary:     *cloudinary.New(),
	}
}
