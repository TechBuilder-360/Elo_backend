package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/hibiken/asynq"
)

func GetEmailProvider(ctx context.Context, logger log.Entry) (IEmailProvider, error) {
	p := provider.NewService()

	serviceLocator, err := p.GetServiceByIdentifier(ctx, constant.EmailServiceIdentifier)
	if err != nil {
		logger.WithError(err).Error("failed to fetch email service provider")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	activeProvider, ok := provider.GetImpl(serviceLocator.ActiveProvider.Slug)
	if !ok {
		logger.Error("provider not supported for Email service")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	providerImpl, ok := provider.ConformsTo[IEmailProvider](activeProvider)
	if !ok {
		logger.WithField("active provider", activeProvider.DisplayName()).Error("provider does not conform to Email provider interface")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return providerImpl, nil
}

func extractEmailTemplate(data interface{}, logger log.Entry, templateName string) (bytes.Buffer, error) {
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("static/email/%s.html", templateName)))
	var htmlContent bytes.Buffer

	err := tmpl.Execute(&htmlContent, data)
	if err != nil {
		logger.WithError(err).Error("failed to parse email template")
		return bytes.Buffer{}, errors.New(errors.ErrFailed, "request failed")
	}
	return htmlContent, nil
}

func HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	logger := log.LoggerInContext(ctx)

	var data WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Email delivery code ...
	impl, err := GetEmailProvider(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get active email provider")
		return errors.New(errors.ErrFailed, "request failed")
	}

	htmlContent, err := extractEmailTemplate(data, logger, "welcome")
	if err != nil {
		return err
	}

	err = impl.Send(ctx, MailPayload{
		HTMLContent: htmlContent.String(),
		Subject:     "Welcome to ELO",
		To:          []Sender{{Email: data.ToMail, Name: data.ToName}},
	})

	return err
}

func HandleOTPEmailTask(ctx context.Context, t *asynq.Task) error {
	logger := log.LoggerInContext(ctx)

	var data OTPMailRequest
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Email delivery code ...
	impl, err := GetEmailProvider(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get active email provider")
		return errors.New(errors.ErrFailed, "request failed")
	}

	htmlContent, err := extractEmailTemplate(data, logger, "otp")
	if err != nil {
		return err
	}

	err = impl.Send(ctx, MailPayload{
		HTMLContent: htmlContent.String(),
		Subject:     "ELO OTP",
		To:          []Sender{{Email: data.ToMail, Name: data.ToName}},
	})

	return err
}

func HandleVericationEmailTask(ctx context.Context, t *asynq.Task) error {
	logger := log.LoggerInContext(ctx)

	var data VerificationMailPayload
	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Email delivery code ...
	impl, err := GetEmailProvider(ctx, logger)
	if err != nil {
		logger.WithError(err).Error("failed to get active email provider")
		return errors.New(errors.ErrFailed, "request failed")
	}

	htmlContent, err := extractEmailTemplate(data, logger, "user_verification")
	if err != nil {
		return err
	}

	err = impl.Send(ctx, MailPayload{
		HTMLContent: htmlContent.String(),
		Subject:     "Your Account Verification Status Has Been Updated",
		To:          []Sender{{Email: data.ToMail, Name: data.FullName}},
	})

	return err
}
