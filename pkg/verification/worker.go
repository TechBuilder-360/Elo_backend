package verification

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/hibiken/asynq"
)

func ProcessVerificationTask(ctx context.Context, t *asynq.Task) error {
	logger := log.LoggerInContext(ctx)

	var payload string
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		logger.WithError(err).Error("Failed to unmarshal task payload")
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	decodedRequest, err := base64.RawStdEncoding.DecodeString(payload)
	if err != nil {
		logger.WithError(err).Error("Failed to decode request string")
		return fmt.Errorf("base64 decode failed: %v: %w", err, asynq.SkipRetry)
	}

	var data VerificationResult
	if err := json.Unmarshal(decodedRequest, &data); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return NewService().ProcessVerification(ctx, &data)
}
