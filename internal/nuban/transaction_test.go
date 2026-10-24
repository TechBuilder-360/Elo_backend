package nuban

import (
	"testing"

	"github.com/Toflex/directory_v2/pkg/provider"
)

func TestSumServiceFees(t *testing.T) {
	fees := []provider.ServiceFee{
		{FeeType: "VAT", Fee: 50},
		{FeeType: "SERVICE_FEE", Fee: 25},
		{FeeType: "PROCESSING_FEE", Fee: 15},
	}

	got := provider.SumServiceFees(fees)
	want := int64(90)
	if got != want {
		t.Fatalf("sumServiceFees() = %d, want %d", got, want)
	}
}
