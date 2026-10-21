package nuban

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
)

func (s *service) GetBankBySlug(ctx context.Context, slug string) (*ent.NubanBank, error) {
	return s.repo.GetBankBySlug(ctx, slug)
}
