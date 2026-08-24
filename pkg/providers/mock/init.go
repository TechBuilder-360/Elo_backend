package mock

import (
	"strings"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
)

type mock struct{}

func New() *mock {
	if configuration.IsProduction() {
		return nil
	}

	return &mock{}
}

func (*mock) Slug() string {
	return strings.ToLower(constant.Mock.ToString())
}

func (*mock) DisplayName() string {
	return constant.Mock.ToString()
}

var _ provider.Impl = (*mock)(nil)
