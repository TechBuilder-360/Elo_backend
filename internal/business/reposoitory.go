package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

type IRepository interface {
	GetBusinessByName(ctx context.Context, name string) (*businessResult, error)
	GetBusinessByID(ctx context.Context, id string) (*businessResult, error)
	Create(ctx context.Context, payload createbusiness) (*ent.Business, error)
	Update(ctx context.Context, b *ent.Business, payload UpdateBusiness) error
	CreateBusinssLocation(ctx context.Context, payload createBusinssLocation) error
	AddDocument(ctx context.Context, business *ent.Business, payload AddBusinessDocument) error
	GetKYBDocuments(ctx context.Context, id string) ([]*ent.BusinessDocument, error)
	DeleteDocument(ctx context.Context, id string) error
	AddBusinessManager(ctx context.Context, businessID, userID string, isOwner bool) error
	GetManager(ctx context.Context, user *ent.User, business *ent.Business) (*ent.Manager, error)
	GetBusinessDocument(ctx context.Context, b *ent.Business, documentID string) (*ent.BusinessDocument, error)
	UpdateDcument(ctx context.Context, bixID, documentID, url string) error
	KYBDocumentByID(ctx context.Context, id string) (*ent.KYBDocument, error)
	KYBDOcuments(ctx context.Context) ([]*ent.KYBDocument, error)
	WithTransaction(tx *ent.Tx) IRepository
}

type repository struct {
	db *ent.Client
}

func Newrepository(db *ent.Client) IRepository {
	return &repository{
		db: db,
	}
}
