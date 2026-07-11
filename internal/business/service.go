package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/internal/manager"
	rbac "github.com/Toflex/directory_v2/pkg/RBAC"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/providers/cloudinary"
	"github.com/samber/do/v2"
)

type IService interface {
	CreateBusiness(ctx context.Context, payload CreateBusinessRequest, logger log.Entry) error
	AddDocument(ctx context.Context, payload UploadDocumentRequest, logger log.Entry) error
	DeleteDocument(ctx context.Context, business *ent.Business, documentID string, logger log.Entry) error
	BusinessDetail(ctx context.Context, user *ent.User, business *ent.Business, payload *BusinessDetailRequest, logger log.Entry) error
	GetBusiness(ctx context.Context, user *ent.User, id string, logger log.Entry) (*BusinessResult, error)
	Businesses(ctx context.Context, user *ent.User, logger log.Entry) ([]*MyBusinessResult, error)
	GetDocuments(ctx context.Context) ([]DocumentResult, error)
	GetKYBDocuments(ctx context.Context, b *ent.Business) ([]KYBDocument, error)
	GetCategory(ctx context.Context) []string
}

type service struct {
	db                *ent.Client
	repo              IRepository
	managerRepository manager.IRepository
	cloudinary        *cloudinary.Cloud
	rbacRepository    rbac.IRepository
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &service{
		db:                db,
		repo:              Newrepository(db),
		managerRepository: manager.NewRepository(db),
		cloudinary:        cloudinary.New(),
		rbacRepository:    rbac.NewRepository(db),
	}
}
