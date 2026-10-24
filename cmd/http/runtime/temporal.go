package runtime

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func initializeTemporal(i do.Injector) {
	c, err := client.Dial(client.Options{
		HostPort: configuration.Instance.TemporalURL,
	})
	if err != nil {
		log.WithError(err).Panic("Unable to create temporal client; skipping worker") // require Temporal worker to be setup
	}
	w := worker.New(c, configuration.Instance.TemporalWorker, worker.Options{})

	// Register Temporal client on DI
	do.Provide(i, func(injector do.Injector) (client.Client, error) {
		return c, nil
	})

	// Register Worker on DI
	do.Provide(i, func(injector do.Injector) (worker.Worker, error) {
		return w, nil
	})

	// Start the worker in a goroutine so it doesn't block initialization.
	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.WithError(err).Warning("Unable to start temporal worker")
		}
	}()
}
