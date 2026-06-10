package verification

import "context"

type VerifyResult struct {
	Link         string `json:"url"`
	ProviderLink string `json:"provider_link"`
	ReferenceID  string `json:"reference_id"`
}

type Verifier interface {
	VerifyUser(ctx context.Context, referenceID string) (*VerifyResult, error)
	VerifyBusiness(ctx context.Context, referenceID string) (*VerifyResult, error)
}
