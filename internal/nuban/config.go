package nuban

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/samber/do/v2"
	"go.temporal.io/sdk/worker"
)

var taskQueue string

var workflows = []interface{}{
	InitiateNubanTransferWorkFlow,
	ProcessNubanTransferChildWorkFlow,
}

func registerWorkflow(i do.Injector, activity *activityWrapper) {
	// If the Temporal worker wasn't provided in DI (e.g. Temporal unavailable),
	// do.MustInvoke will panic. Recover and skip registration in that case.
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	taskQueue = configuration.Instance.TemporalWorker

	w := do.MustInvoke[worker.Worker](i)

	// register workflow
	for _, v := range workflows {
		w.RegisterWorkflow(v)
	}

	// register activities
	activities := []interface{}{
		activity.GetBankCode,
		activity.GetTransferService,
		activity.AccountResolutionActivity,
		activity.BeneficiaryChecklist,
		activity.createTransferTransaction,
		activity.transferToProvider,
		activity.transferStatus,
		activity.updateTransaction,
		activity.reverseTransfer,
	}
	for _, v := range activities {
		w.RegisterActivity(v)
	}

}
