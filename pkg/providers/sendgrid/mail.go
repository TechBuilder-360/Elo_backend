package sendgrid

import (
	"context"

	"github.com/Toflex/directory_v2/internal/email"
)

type MailPayload struct {
	HTMLContent string   `json:"htmlContent"`
	Sender      Sender   `json:"sender"`
	Subject     string   `json:"subject"`
	To          []Sender `json:"to"`
}

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *sendgrid) Send(ctx context.Context, payload email.MailPayload) error {

	return nil
}
