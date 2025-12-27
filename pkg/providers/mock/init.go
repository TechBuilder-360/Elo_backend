package mock

import "github.com/Toflex/directory_v2/pkg/constant"

type mock struct{}

func New() mock {
	return mock{}
}

func (mock) Slug() string {
	return constant.SendGrid.ToString()
}

func (mock) DisplayName() string {
	return constant.SendGrid.ToString()
}
