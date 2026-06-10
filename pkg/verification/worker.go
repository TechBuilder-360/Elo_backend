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

	decodedRequest, err := base64.RawStdEncoding.DecodeString(string(t.Payload()))
	if err != nil {
		logger.WithError(err).Error("Failed to decode request sting")
		return err
	}

	var data Verification
	if err := json.Unmarshal(decodedRequest, &data); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}
