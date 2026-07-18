package runtime

import (
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/internal/business"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/internal/transaction"
	"github.com/Toflex/directory_v2/internal/wallet"
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
}
