package cloudinary

// FileDestination represents the location where a file should be uploaded
type FileDestination struct {
	Type string // e.g., "verification", "document"
	ID   string // business_id or user_id
}

// Upload type constants
const (
	VerificationType = "verification"
	DocumentType     = "document"
)

// Entity type constants
const (
	BusinessEntity = "businesses"
	UserEntity     = "users"
)
