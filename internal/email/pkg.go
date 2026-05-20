package email

import "context"

type WelcomeEmailPayload struct {
	ToMail   string
	ToName   string
	FullName string
}

type OTPMailRequest struct {
	ToMail string
	ToName string
	OTP    string
}

type MailPayload struct {
	HTMLContent string   `json:"htmlContent"`
	Subject     string   `json:"subject"`
	To          []Sender `json:"to"`
}

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type IEmailProvider interface {
	// SendWelcomeMail(data WelcomeMailRequest) error
	// SendOTPMail(data OTPMailRequest) error

	Send(context.Context, MailPayload) error
}
