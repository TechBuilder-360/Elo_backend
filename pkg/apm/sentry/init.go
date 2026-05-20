package sentry

import (
	"time"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	sentrylogrus "github.com/getsentry/sentry-go/logrus"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sirupsen/logrus"
)

type config struct {
	SentryURL string `env:"SENTRY_URL"`
}

var client sentry.ClientOptions

// var sentryFiber fiber.Handler

func initClientOpt() {
	conf := &config{}
	conf = configuration.Load(conf).(*config)

	client = sentry.ClientOptions{
		Dsn:              conf.SentryURL,
		TracesSampleRate: 0.5,
		Debug:            !configuration.IsProduction(),
		AttachStacktrace: true,
		EnableTracing:    configuration.IsProduction(),
	}
}

func ClientOpt() sentry.ClientOptions {
	return client
}

func InitializeSentry(i do.Injector, l *logrus.Logger) (*sentrylogrus.Hook, error) {
	initClientOpt()

	if err := sentry.Init(ClientOpt()); err != nil {
		log.Info("Sentry initialization failed: %v\n", err)
		return nil, err
	}

	// Send only ERROR and higher level logs to Sentry
	sentryLevels := []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}

	// Initialize Sentry
	sentryHook, err := sentrylogrus.New(sentryLevels, client)
	if err != nil {
		log.Info("Sentry initialization failed: %v\n", err)
		return nil, err
	}
	log.AddHook(sentryHook)

	logrus.RegisterExitHandler(func() { sentryHook.Flush(5 * time.Second) })

	sentryGin := sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
	})

	engine := do.MustInvoke[*gin.Engine](i)
	engine.Use(sentryGin)

	return sentryHook, nil
}
