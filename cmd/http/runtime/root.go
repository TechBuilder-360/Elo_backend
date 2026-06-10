package runtime

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/providers/brevo"
	"github.com/Toflex/directory_v2/pkg/providers/mock"
)

func Register() {
	if !configuration.IsProduction() {
		provider.RegisterProvider(mock.New())
	} else {
		provider.RegisterProvider(brevo.New())
	}
}
