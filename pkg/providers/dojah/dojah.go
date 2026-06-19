package dojah

import (
	"strings"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/verification"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type IDojah interface {
	RegisterRoutes(engine *gin.Engine)
}

type config struct {
	BaseURL                      string `env:"DOJAH_BASE_URL" required:"true"`
	IdentityBaseURL              string `env:"DOJAH_IDENTITY_BASE_URL" required:"true"`
	ApiKey                       string `env:"DOJAH_API_KEY" required:"true"`
	SecretKey                    string `env:"DOJAH_SECRET_KEY" required:"true"`
	UserVerificationWidgetID     string `env:"DOJAH_USER_VERIFICATION_WIDGET_ID" required:"true"`
	BusinessVerificationWidgetID string `env:"DOJAH_BUSINESS_VERIFICATION_WIDGET_ID" required:"true"`
}

type dojah struct {
	config config
	engine *gin.Engine
}

func New(i do.Injector) *dojah {
	config := config{}
	configuration.Load(&config)

	// http engine
	engine := do.MustInvoke[*gin.Engine](i)

	return &dojah{
		config: config,
		engine: engine,
	}
}

// DisplayName implements [provider.Impl].
func (d *dojah) DisplayName() string {
	return constant.Dojah.ToString()
}

// Slug implements [provider.Impl].
func (d *dojah) Slug() string {
	return strings.ToLower(constant.Dojah.ToString())
}

var _ provider.Impl = (*dojah)(nil)
var _ verification.Verifier = (*dojah)(nil)
