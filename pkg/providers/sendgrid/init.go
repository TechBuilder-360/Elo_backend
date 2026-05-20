package sendgrid

import (
	"strings"

	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/provider"
)

type sendgrid struct {
}

func (*sendgrid) Slug() string {
	return strings.ToLower(constant.SendGrid.ToString())
}

func (*sendgrid) DisplayName() string {
	return constant.SendGrid.ToString()
}

func New() *sendgrid {
	return &sendgrid{}
}

var _ provider.Impl = (*sendgrid)(nil)
var _ email.IEmailProvider = (*sendgrid)(nil)
