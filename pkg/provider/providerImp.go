package provider

import (
	"strings"
	"sync"
)

type Impl interface {
	DisplayName() string
	Slug() string
}

var mutex = sync.RWMutex{}
var providers = map[string]Impl{}

// RegisterProvider registers a provider
func RegisterProvider(provider ...Impl) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, p := range provider {
		providers[p.Slug()] = p
	}
}

// GetProviderImpl returns a provider implementation
func GetProviderImpl(slug string) (Impl, bool) {
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

// GetImpl checks if a provider conforms to a type
func GetImpl[T any](slug string) (T, bool) {
	var result T

	p, ok := providers[strings.ToLower(slug)]
	if !ok {
		return result, false
	}

	if v, ok := p.(T); ok {
		return v, true
	}

	return result, false
}

// GetProviders returns list of registered providers
func GetProviders() map[string]Impl {
	return providers
}
