package nuban

import (
	"time"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type NubanTransferRequest struct {
	OwnerID       string             `json:"owner_id"`
	AccountNumber string             `json:"account_number"`
	BankSlug      string             `json:"bank_slug"`
	Narration     *string            `json:"narration"`
	Amount        uint64             `json:"amount"`
	Reference     string             `json:"reference"`
	Currency      types.CurrencyCode `json:"currency"`
}

type transactionPayload struct {
	OwnerId          string
	AccountName      string             `json:"account_name"`
	AccountNumber    string             `json:"account_number"`
	ProviderBankCode string             `json:"provider_bank_code"`
	Narration        *string            `json:"narration"`
	Amount           uint64             `json:"amount"`
	Reference        string             `json:"reference"`
	Currency         types.CurrencyCode `json:"currency"`
	Bank             *bankResponse      `json:"bank"`
	Service          provider.ServiceLocator
	BeneficiaryID    string
}

type NubanTransferResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type bankResponse struct {
	ID           string
	ProviderCode string
}

func InitiateNubanTransferWorkFlow(ctx workflow.Context, param *NubanTransferRequest) (*NubanTransferResponse, error) {
	var (
		transactionResult *NubanTransferResponse
		service           provider.ServiceLocator
		bank              *bankResponse
		beneficiaryResult *nubanBeneficiaryResult
		logger            = workflow.GetLogger(ctx)
	)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    100 * time.Second,
			MaximumAttempts:    2, // 0 = unlimited retries until timeout
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Get Transfer Service Provider
	err := workflow.ExecuteActivity(ctx, "GetTransferService").Get(ctx, &service)
	if err != nil {
		logger.Error("Failed to get transfer service", err.Error())
		return nil, err
	}

	// Get Bank Code
	bankCodeRequest := &bankPayload{
		Slug:     param.BankSlug,
		Provider: service.ActiveProvider.Slug,
	}
	err = workflow.ExecuteActivity(ctx, "GetBankCode", bankCodeRequest).Get(ctx, &bank)
	if err != nil {
		return nil, err
	}

	// Account number resolution
	resolveAccountRequest := &ResolveAccountNamePayload{
		AccountNumber: param.AccountNumber,
		BankCode:      bank.ProviderCode,
	}

	var resolveAccountResponse *ResolveAccountNameResult

	err = workflow.ExecuteActivity(ctx, "AccountResolutionActivity", resolveAccountRequest).Get(ctx, &resolveAccountResponse)
	if err != nil {
		return nil, err
	}

	// Beneficiary
	beneficiartPayload := &nubanBeneficiaryPayload{
		OwnerID:       param.OwnerID,
		AccountName:   resolveAccountResponse.AccountName,
		AccountNumber: resolveAccountResponse.AccountNumber,
		BankID:        bank.ID,
	}

	err = workflow.ExecuteActivity(ctx, "BeneficiaryChecklist", beneficiartPayload).Get(ctx, &beneficiaryResult)
	if err != nil {
		return nil, err
	}

	tRequest := &transactionPayload{
		OwnerId:          param.OwnerID,
		AccountName:      resolveAccountResponse.AccountName,
		AccountNumber:    resolveAccountResponse.AccountNumber,
		ProviderBankCode: resolveAccountResponse.BankCode,
		Narration:        param.Narration,
		Amount:           param.Amount,
		Reference:        param.Reference,
		Currency:         param.Currency,
		Bank:             bank,
		Service:          service,
		BeneficiaryID:    beneficiaryResult.ID,
	}

	err = workflow.ExecuteActivity(ctx, "createTransferTransaction", tRequest).Get(ctx, &transactionResult)
	if err != nil {
		return nil, err
	}

	// Start child Workflow that make request to provider and wait for status update
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, ProcessNubanTransferChildWorkFlow, &ProcessTransferPayload{
		TransactionID: transactionResult.ID,
		Narration:     param.Narration,
		Amount:        param.Amount,
		Provider:      service.ActiveProvider.Slug,
		Currency:      param.Currency.ToString(),
		Beneficiary:   beneficiaryResult,
		Bank:          bank,
	})
	// Wait for the Child Workflow Execution to spawn
	var childWE workflow.Execution
	if err := childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		logger.Error("failed to start child workflow to prosess nuban transfer %s", err.Error())
		return nil, err
	}

	return transactionResult, nil
}

func ProcessNubanTransferChildWorkFlow(ctx workflow.Context, params *ProcessTransferPayload) (*NubanTransferResponse, error) {
	var (
		result *NubanTransferResult
		logger = workflow.GetLogger(ctx)
	)
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    3 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1, // 0 = unlimited retries until timeout
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	err := workflow.ExecuteActivity(ctx, "transferToProvider", params).Get(ctx, &result)
	if err != nil {
		logger.Error("transferToProvider failed", "error", err.Error())

		// Update transaction status to FAILED
		updateErr := workflow.ExecuteActivity(ctx, "updateTransaction", &transactionUpdatePayload{
			ID:     params.TransactionID,
			Status: constant.TransactionFailedStatus,
		}).Get(ctx, nil)
		if updateErr != nil {
			logger.Error("failed to update transaction status to FAILED", "error", updateErr.Error())
		}

		// Reverse funds for failed transfer
		revActivityOptions := workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    3 * time.Second,
				BackoffCoefficient: 2.0,
				MaximumInterval:    time.Minute,
				MaximumAttempts:    5,
			},
		}
		revCtx := workflow.WithActivityOptions(ctx, revActivityOptions)
		revErr := workflow.ExecuteActivity(revCtx, "reverseTransfer", params.TransactionID).Get(ctx, nil)
		if revErr != nil {
			logger.Error("failed to reverse transfer funds", "error", revErr.Error())
		}

		return nil, err
	}

	// If Status = 'PENDING', poll status
	if result.Status == constant.TransactionPendingStatus {
		statusActivityOptions := workflow.ActivityOptions{
			StartToCloseTimeout: 10 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    3 * time.Second,
				BackoffCoefficient: 2.0,
				MaximumInterval:    time.Minute,
				MaximumAttempts:    5,
			},
		}
		statusCtx := workflow.WithActivityOptions(ctx, statusActivityOptions)
		for i := 0; i < 5; i++ {
			_ = workflow.Sleep(ctx, 3*time.Second)
			err = workflow.ExecuteActivity(statusCtx, "transferStatus", params.TransactionID, params.Provider).Get(ctx, &result)
			if err != nil || result.Status != constant.TransactionPendingStatus {
				break
			}
		}
	}

	// Update transaction Status
	updateActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    3 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5, // 0 = unlimited retries until timeout
		},
	}
	updateCtx := workflow.WithActivityOptions(ctx, updateActivityOptions)
	err = workflow.ExecuteActivity(updateCtx, "updateTransaction", &transactionUpdatePayload{
		ID:               params.TransactionID,
		Status:           result.Status,
		SessionID:        result.SessionID,
		PartnerReference: &result.PartnerID,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("failed to update transaction", "error", err.Error())
		return nil, err
	}

	// Perform post transfer operation
	// Reverse funds for failed transfer
	if result.Status == constant.TransactionFailedStatus {
		revActivityOptions := workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    3 * time.Second,
				BackoffCoefficient: 2.0,
				MaximumInterval:    time.Minute,
				MaximumAttempts:    5,
			},
		}
		revCtx := workflow.WithActivityOptions(ctx, revActivityOptions)
		err = workflow.ExecuteActivity(revCtx, "reverseTransfer", params.TransactionID).Get(ctx, nil)
		if err != nil {
			logger.Error("failed to reverse funds for failed transfer", "error", err.Error())
			return nil, err
		}
	}

	return &NubanTransferResponse{
		ID:     params.TransactionID,
		Status: string(result.Status),
	}, nil
}
