package cloudinary

import (
	"fmt"
	"strings"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	c "github.com/cloudinary/cloudinary-go/v2"
)

type config struct {
	Apikey    string `env:"CLOUDINARY_API_KEY" required:"true"`
	Secret    string `env:"CLOUDINARY_SECRET" required:"true"`
	CloudName string `env:"CLOUDINARY_CLOUD_NAME" required:"true"`
}

type Cloud struct {
	config config
	cld    *c.Cloudinary
}

// DisplayName implements [provider.Impl].
func (c *Cloud) DisplayName() string {
	return constant.Cloudinary.ToString()
}

// Slug implements [provider.Impl].
func (c *Cloud) Slug() string {
	return strings.ToLower(constant.Cloudinary.ToString())
}

func New() *Cloud {
	config := config{}
	configuration.Load(&config)

	url := fmt.Sprintf("cloudinary://%s:%s@%s", config.Apikey, config.Secret, config.CloudName)
	cld, err := c.NewFromURL(url)
	if err != nil {
		log.WithError(err).Error("Failed to initialize cloudinary client")
	}

	return &Cloud{
		config: config,
		cld:    cld,
	}
}

var _ provider.Impl = (*Cloud)(nil)
