package constant

import "github.com/Toflex/directory_v2/pkg/types"

const (
	JWTId types.JWTKey = "JWT"
)

// Queue task types
const (
	TaskTypeWelcomeEmail         string = "email:welcome"
	TaskTypeOTPEmail             string = "email:otp"
	TaskTypeIdentityVerification string = "identity:verification"
	TaskUserVerification         string = "user:verification"
	TaskKYBDocument              string = "kyb:document"
	TaskDafaultWallet            string = "wallet:default"
)

// Verification Types
const (
	UserEntityType     types.EntityType = "USER_VERIFICATION"
	BusinessEntityType types.EntityType = "BUSINESS_VERIFICATION"
)

const (
	Success types.VerificationStatus = "SUCCESS"
	Failed  types.VerificationStatus = "FAILED"
	Pending types.VerificationStatus = "PENDING"
)

const (
	PDF  types.DocumentType = "pdf"
	JPEG types.DocumentType = "jpeg"
	JPG  types.DocumentType = "jpg"
)

const (
	Business types.LedgerOwner = "BUSINESS"
	User     types.LedgerOwner = "USER"
)
