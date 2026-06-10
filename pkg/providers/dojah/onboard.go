package dojah

import (
	"context"
	"fmt"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/verification"
)

// VerifyUser implements [verification.Verifier].
func (d *dojah) VerifyUser(ctx context.Context, referenceID string) (*verification.VerifyResult, error) {
	url := fmt.Sprintf("%s?widget_id=%s&reference_id=%s", d.config.IdentityBaseURL, d.config.UserVerificationWidgetID, referenceID)

	return &verification.VerifyResult{
		Link:         url,
		ReferenceID:  referenceID,
		ProviderLink: fmt.Sprintf("%s/%s", configuration.GetBaseURL(), referenceID),
	}, nil
}

// VerifyBusiness implements [verification.Verifier].
func (d *dojah) VerifyBusiness(ctx context.Context, referenceID string) (*verification.VerifyResult, error) {
	url := fmt.Sprintf("%s?widget_id=%s&reference_id=%s", d.config.IdentityBaseURL, d.config.BusinessVerificationWidgetID, referenceID)

	return &verification.VerifyResult{
		Link:         url,
		ReferenceID:  referenceID,
		ProviderLink: fmt.Sprintf("%s/%s", configuration.GetBaseURL(), referenceID),
	}, nil
}
