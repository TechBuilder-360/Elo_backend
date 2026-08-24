package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/Toflex/directory_v2/internal/nuban"
	"github.com/Toflex/directory_v2/pkg/util"
)

// GenerateMockNUBANAccountNumber generates a mock NUBAN account number for testing purposes.
func (m *mock) generateMockNUBANAccountNumber(length int) string {
	if length <= 0 {
		return ""
	}

	digits := "0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = digits[rand.Intn(10)]
	}

	return string(result)
}

func (m *mock) GenerateStaticNUBANAccountNumber(ctx context.Context, payload *nuban.StaticNubanRequest) (*nuban.StaticNubanResponse, error) {
	return &nuban.StaticNubanResponse{
		AccountNumber:    m.generateMockNUBANAccountNumber(10), // Generate a random 10-digit account number
		AccountName:      payload.AccountName,
		BankCode:         "000",
		BankName:         "ELO Bank",
		PartnerReference: util.GenerateRandomString(15), // Generate a random partner reference
		Type:             nuban.StaticNUBANAccountType,
	}, nil
}

// GenerateDynamicNUBANAccountNumber generates a mock dynamic NUBAN account number for testing purposes.
func (m *mock) GenerateDynamicNUBANAccountNumber(ctx context.Context, payload *nuban.DynamicNubanRequest) (*nuban.DynamicNubanResponse, error) {
	return &nuban.DynamicNubanResponse{
		AccountNumber:    m.generateMockNUBANAccountNumber(10), // Generate a random 10-digit account number
		AccountName:      payload.AccountName,
		Amount:           payload.Amount,
		BankCode:         "000",
		BankName:         "ELO Bank",
		PartnerReference: util.GenerateRandomString(15), // Generate a random partner reference
		Type:             nuban.DynamicNUBANAccountType,
		ExpiryDate:       time.Now().Add(30 * time.Minute), // Set expiry date to 30 minutes from now
	}, nil
}

var _ nuban.NUBANProvider = (*mock)(nil)
