package business

import (
	"context"

	"github.com/Toflex/directory_v2/pkg/log"
)

func (s *service) GetBusinessByID(ctx context.Context, id string, logger log.Entry) (*businessResult, error) {
	return s.repo.GetBusinessByID(ctx, id)
}
