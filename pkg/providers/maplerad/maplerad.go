package maplerad

import (
	"strings"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type IMaplerad interface {
	RegisterRoutes(engine *gin.Engine)
	GetConfig() config
}

type config struct {
	BaseURL   string `env:"MAPLERAD_BASE_URL" required:"true"`
	SecretKey string `env:"MAPLERAD_SECRET_KEY" required:"true"`
}

type maplerad struct {
	config config
	engine *gin.Engine
}

func New(i do.Injector) *maplerad {
	config := config{}
	configuration.Load(&config)

	// http engine
	engine := do.MustInvoke[*gin.Engine](i)

	return &maplerad{
		config: config,
		engine: engine,
	}
}

// DisplayName implements [provider.Impl].
func (m *maplerad) DisplayName() string {
	return constant.Maplerad.ToString()
}

// Slug implements [provider.Impl].
func (m *maplerad) Slug() string {
	return strings.ToLower(constant.Maplerad.ToString())
}

var _ provider.Impl = (*maplerad)(nil)
