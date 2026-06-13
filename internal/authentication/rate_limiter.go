package authentication

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

const (
	registrationRateLimitCount  = 3
	registrationRateLimitWindow = time.Minute * 10
	otpRequestRateLimitCount    = 3
	otpRequestRateLimitWindow   = time.Minute * 10
	loginRateLimitCount         = 3
	loginRateLimitWindow        = time.Minute * 5
)

func (s *Service) enforceRateLimit(ctx context.Context, key string, limit int64, window time.Duration, logger log.Entry) error {
	count, err := s.cache.IncrWithExpire(ctx, key, window)
	if err != nil {
		logger.WithError(err).Error("failed to enforce rate limit")
		return errors.New(errors.ErrInternal, "request failed")
	}

	if count > limit {
		logger.WithFields(map[string]interface{}{"key": key, "count": count, "limit": limit}).Warning("rate limit exceeded")
		return errors.New(errors.ErrTooManyRequests, "too many requests, please try again later")
	}

	return nil
}

func (s *Service) resetRateLimit(ctx context.Context, key string, logger log.Entry) {
	if err := s.cache.Delete(ctx, key); err != nil {
		logger.WithError(err).Error("failed to reset rate limit")
	}
}

func registrationRateLimitKey(email string) string {
	return fmt.Sprintf("registration:%s", strings.ToLower(strings.TrimSpace(email)))
}

func requestOTPRateLimitKey(email string) string {
	return fmt.Sprintf("otp_request:%s", strings.ToLower(strings.TrimSpace(email)))
}

func loginRateLimitKey(identifier string) string {
	return fmt.Sprintf("login:%s", strings.TrimSpace(identifier))
}
