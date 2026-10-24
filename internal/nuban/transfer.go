package nuban

import (
	"context"
	stderrors "errors"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/transaction"
	"github.com/Toflex/directory_v2/internal/currency"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

func (s *service) InitiateTransfer(ctx context.Context, payload *NubanTransferRequest) (*NubanTransferResponse, error) {
	logger := log.LoggerInContext(ctx)

	if payload.Currency != currency.NGN {
		return nil, errors.New(errors.ErrInvalidInput, "currency not supported")
	}

	if payload.Reference == "" {
		payload.Reference = util.GenerateReference("NUBAN")
	}

	existing, err := s.db.Transaction.Query().
		Where(transaction.ReferenceEQ(payload.Reference)).
		Only(ctx)
	if err == nil {
		return &NubanTransferResponse{
			ID:     existing.ID,
			Status: existing.Status.String(),
		}, nil
	}
	if !ent.IsNotFound(err) {
		logger.WithError(err).WithField("reference", payload.Reference).Error("failed to check transfer reference")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	tc := s.temporalClient

	// check that temporal client is not nil
	if tc == nil {
		logger.Error("temporal client is not initialized")
		return nil, errors.New(errors.ErrFailed, "")
	}

	options := client.StartWorkflowOptions{
		ID:                                       payload.Reference,
		TaskQueue:                                taskQueue,
		WorkflowIDConflictPolicy:                 enumspb.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowIDReusePolicy:                    enumspb.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
		WorkflowExecutionErrorWhenAlreadyStarted: false,
	}

	workflowRun, err := tc.ExecuteWorkflow(ctx, options, InitiateNubanTransferWorkFlow, payload)
	if err != nil {
		logger.WithField("payload", payload).WithError(err).Error("Unable to execute workflow")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	var result *NubanTransferResponse
	err = workflowRun.Get(ctx, &result)
	if err != nil {
		logger.WithError(err).Error("workflow returned error")
		return nil, unwrapTemporalError(err)
	}
	if result == nil {
		logger.Error("workflow returned an empty result")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return result, nil

}

func unwrapTemporalError(err error) error {
	var applicationErr *temporal.ApplicationError
	if stderrors.As(err, &applicationErr) {
		msg := errors.SanitizeMessage(applicationErr.Message())
		return errors.New(errors.ErrFailed, msg)
	}

	for {
		cause := stderrors.Unwrap(err)
		if cause == nil {
			if errors.IsDBError(err) {
				return errors.New(errors.ErrFailed, "request failed")
			}
			return err
		}

		err = cause
	}
}
