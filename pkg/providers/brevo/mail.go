package brevo

import (
	"context"
	"errors"

	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/apm"
	"github.com/Toflex/directory_v2/pkg/log"
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

type Success struct {
	MessageID string `json:"messageId"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (b *brevo) Send(ctx context.Context, payload email.MailPayload) error {
	logger := log.LoggerInContext(ctx).WithField("provider", b.DisplayName())
	errorResponse := &ErrorResponse{}
	to := []Sender{}

	for _, recipient := range payload.To {
		to = append(to, Sender{
			Email: recipient.Email,
			Name:  recipient.Name,
		})
	}

	request := &MailPayload{
		HTMLContent: payload.HTMLContent,
		To:          to,
		Subject:     payload.Subject,
		Sender: Sender{
			Email: "tech.builder.circle@gmail.com",
			Name:  "Elo",
		},
	}

	resp, err := apm.HTTPClientRequest().
		SetContext(ctx).
		SetHeader("api-key", b.config.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetError(errorResponse).
		Post(b.config.BaseURL + "/v3/smtp/email")

	if err != nil {
		logger.WithError(err).Error("failed to send email via Brevo")
		return err
	}

	if resp.IsError() {
		logger.Error("failed to send email via Brevo: %v", errorResponse)
		return errors.New(errorResponse.Message)
	}
	return nil
}
