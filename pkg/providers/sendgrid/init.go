package sendgrid

import "github.com/Toflex/directory_v2/pkg/constant"

type sendgrid struct {
}

func (sendgrid) Slug() string {
	return constant.SendGrid.ToString()
}

func New() sendgrid {
	return sendgrid{}
}
