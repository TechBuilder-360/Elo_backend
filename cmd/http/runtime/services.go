package runtime

import (
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/internal/business"
	"github.com/Toflex/directory_v2/internal/crypto"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/internal/nuban"
	"github.com/Toflex/directory_v2/internal/transaction"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/Toflex/directory_v2/pkg/verification"
	"github.com/samber/do/v2"
)

func initializeService(i do.Injector) {
	do.Provide(i, func(i do.Injector) (authentication.IService, error) {
		return authentication.NewService(i), nil
	})
	do.Provide(i, func(i do.Injector) (business.IService, error) {
		return business.NewService(i), nil
	})
	do.Provide(i, func(i do.Injector) (wallet.IService, error) {
		return wallet.Newservice(i), nil
	})
	do.Provide(i, func(i do.Injector) (transaction.IService, error) {
		return transaction.Newservice(i), nil
	})
	do.Provide(i, func(i do.Injector) (currency.IService, error) {
		return currency.Newservice(i), nil
	})
	do.Provide(i, func(i do.Injector) (nuban.IService, error) {
		return nuban.NewService(i), nil
	})
	do.Provide(i, func(i do.Injector) (crypto.IService, error) {
		return crypto.NewService(i), nil
	})
}

type RegisteredService struct {
	AuthenticationService authentication.IService
	VerificationService   verification.IService
	BusinessService       business.IService
	WalletService         wallet.IService
	CurrencyService       currency.IService
	TransactionService    transaction.IService
	NubanService          nuban.IService
	CryptoService         crypto.IService
}

// NewService initializes and returns a new instance of the Service struct.
func NewService(i do.Injector) *RegisteredService {
	s := &RegisteredService{
		AuthenticationService: do.MustInvoke[authentication.IService](i),
		VerificationService:   verification.NewService(),
		BusinessService:       do.MustInvoke[business.IService](i),
		WalletService:         do.MustInvoke[wallet.IService](i),
		CurrencyService:       do.MustInvoke[currency.IService](i),
		TransactionService:    do.MustInvoke[transaction.IService](i),
		NubanService:          do.MustInvoke[nuban.IService](i),
		CryptoService:         do.MustInvoke[crypto.IService](i),
	}
	return s
}
