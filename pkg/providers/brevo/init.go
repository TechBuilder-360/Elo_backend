package brevo

import (
	"strings"

	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
)

type brevo struct {
	config config
}

type config struct {
	baseURL string `env:"BREVO_BASE_URL" required:"true"`
	apiKey  string `env:"BREVO_API_KEY" required:"true"`
}

func (*brevo) Slug() string {
	return strings.ToLower(constant.Brevo.ToString())
}

func (*brevo) DisplayName() string {
	return constant.Brevo.ToString()
}

func New() *brevo {
	config := config{}
	configuration.Load(&config)

	return &brevo{
		config: config,
	}
}

var _ provider.Impl = (*brevo)(nil)
var _ email.IEmailProvider = (*brevo)(nil)
