package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/samber/do/v2"
	"go.temporal.io/sdk/client"
)

type IService interface {
	GenerateBusinessStaticNuban(ctx context.Context, ownerID string) (*StaticNubanResponse, error)
	GenerateBusinessDynamicNuban(ctx context.Context, ownerID string) (*DynamicNubanResponse, error)
	InitiateTransfer(ctx context.Context, payload *NubanTransferRequest) (*NubanTransferResponse, error)
	GetBankBySlug(ctx context.Context, slug string) (*ent.NubanBank, error)
	getServiceProvider(ctx context.Context, identifier string) (provider.ServiceLocator, error)
	ResolveAccountName(ctx context.Context, payload ResolveAccountNamePayload) (*ResolveAccountNameResult, error)
	beneficiaryCheck(ctx context.Context, payload *nubanBeneficiaryPayload) (*nubanBeneficiaryResult, error)
	createTransferTransaction(ctx context.Context, payload *transactionPayload) (*NubanTransferResponse, error)
	updateTransferTransactionStatus(ctx context.Context, param *transactionUpdatePayload) error
	reverseTransfer(ctx context.Context, transactionID string) error
}

type service struct {
	db             *ent.Client
	cache          *redis.Client
	repo           IRepository
	walletService  wallet.IService
	temporalClient client.Client
}

func NewService(i do.Injector) IService {
	temporalClient, err := do.Invoke[client.Client](i)
	if err != nil {
		log.Panic("failed to invoke temporal client %s", err.Error())
	}

	db := do.MustInvoke[*ent.Client](i)
	cache := do.MustInvoke[*redis.Client](i)
	walletService := do.MustInvoke[wallet.IService](i)
	svc := &service{
		db:             db,
		repo:           NewRepository(db),
		cache:          cache,
		walletService:  walletService,
		temporalClient: temporalClient,
	}

	if temporalClient != nil {
		// Register Temporal handlers with a wrapper bound to this service.
		registerWorkflow(i, &activityWrapper{svc: svc})
	}

	return svc
}
