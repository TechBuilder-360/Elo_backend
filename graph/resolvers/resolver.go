package resolver

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/Toflex/directory_v2/cmd/http/runtime"
	"github.com/Toflex/directory_v2/graph/generated"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Services *runtime.RegisteredService
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
