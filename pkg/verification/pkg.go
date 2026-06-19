package verification

import "context"

type VerifyResult struct {
	Link        string `json:"url"`
	ReferenceID string `json:"reference_id"`
}

type Verifier interface {
	VerifyUser(ctx context.Context, referenceID string) (*VerifyResult, error)
	VerifyBusiness(ctx context.Context, referenceID string) (*VerifyResult, error)
}
