package runtime

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/providers/brevo"
	"github.com/Toflex/directory_v2/pkg/providers/dojah"
	"github.com/Toflex/directory_v2/pkg/providers/mock"
	"github.com/samber/do/v2"
)

func Register(i do.Injector) {
	if !configuration.IsProduction() {
		provider.RegisterProvider(mock.New())
	} else {
		provider.RegisterProvider(brevo.New(),
			dojah.New(i))
	}
}
