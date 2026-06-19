package authentication

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/log"
)

func (s *Service) Logout(ctx context.Context, user *ent.User) {
	logger := log.LoggerInContext(ctx)

	err := s.cache.Delete(ctx, fmt.Sprintf("auth::%s", user.ID))
	if err != nil {
		logger.WithError(err).Error("failed to delete jwt token")
	}
}
