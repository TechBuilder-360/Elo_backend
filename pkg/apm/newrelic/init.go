package newrelic

import (
	"os"
	"time"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type config struct {
	NewRelicAppName string `env:"NEW_RELIC_APP_NAME"`
	NewRelicLicense string `env:"NEW_RELIC_LICENSE"`
}

func InitialiseNewRelic(l *logrus.Logger) {
	conf := &config{}
	conf = configuration.Load(conf).(*config)

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(conf.NewRelicAppName),
		newrelic.ConfigLicense(conf.NewRelicLicense),
		newrelic.ConfigAppLogForwardingEnabled(false),
		//nrlogrus.ConfigLogger(l),
		newrelic.ConfigAppLogEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigAppLogDecoratingEnabled(true),
		newrelic.ConfigInfoLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		log.Error(err.Error())
	}

	err = app.WaitForConnection(10 * time.Second)
	if err != nil {
		log.Error(err.Error())
	}

}
