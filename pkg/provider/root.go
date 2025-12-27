package provider

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/providers/mock"
)

func Register() {
	if configuration.IsSandbox() || configuration.IsDevelopment() {
		RegisterProvider(
			mock.New())
	} else {
		RegisterProvider()
	}
}
