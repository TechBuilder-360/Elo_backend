package constant

import "github.com/Toflex/directory_v2/pkg/types"

const (
	JWTId types.JWTKey = "JWT"
)

// Queue task types
const (
	TaskTypeWelcomeEmail string = "email:welcome"
	TaskTypeOTPEmail     string = "email:otp"
)
