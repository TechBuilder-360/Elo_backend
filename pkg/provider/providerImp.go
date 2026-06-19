package provider

import (
	"strings"
	"sync"

	"github.com/Toflex/directory_v2/pkg/constant"
)

type Impl interface {
	DisplayName() string
	Slug() string
}

var mutex = sync.RWMutex{}
var providers = map[string]Impl{}

var ServiceProviders = []string{}

func init() {
	ServiceProviders = append(ServiceProviders,
		constant.Brevo.ToString(), constant.Dojah.ToString())
}

// RegisterProvider registers a provider
func RegisterProvider(provider ...Impl) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, p := range provider {
		providers[p.Slug()] = p
	}
}

// GetImpl returns a provider implementation
func GetImpl(slug string) (Impl, bool) {
	p, ok := providers[strings.ToLower(slug)]
	return p, ok
}

// ConformsTo checks if a provider conforms to a type
func ConformsTo[T any](p Impl) (T, bool) {
	var result T

	if v, ok := p.(T); ok {
		return v, true
	}

	return result, false
}

// GetProviders returns list of registered providers
func GetProviders() map[string]Impl {
	return providers
}
