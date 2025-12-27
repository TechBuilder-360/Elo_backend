package pkg

type WelcomeMailRequest struct {
	ToMail   string
	ToName   string
	FullName string
}

type OTPMailRequest struct {
	ToMail string
	ToName string
	OTP    string
}

type IEmailProvider interface {
	SendWelcomeMail(data WelcomeMailRequest) error
	SendOTPMail(data OTPMailRequest) error
}
