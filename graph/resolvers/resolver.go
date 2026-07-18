package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/Toflex/directory_v2/graph/generated"
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/internal/business"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/internal/transaction"
	"github.com/Toflex/directory_v2/internal/wallet"
	"github.com/Toflex/directory_v2/pkg/verification"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthenticationService authentication.IService
	VerificationService   verification.IService
	BusinessService       business.IService
	WalletService         wallet.IService
	CurrencyService       currency.IService
	TransactionService    transaction.IService
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
