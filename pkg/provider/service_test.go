package provider

import (
	"context"
	"testing"
)

func int64Ptr(v int64) *int64 {
	return &v
}

func TestCalculateFee_AppliesTo(t *testing.T) {
	ctx := context.Background()

	// Service with two rules:
	// Rule 1: Flat NGN 50 EMTL levy, applies only for amount >= 10,000
	// Rule 2: 1.5% processing fee, applies for all amounts, min fee 10, max fee 1000
	service := ServiceLocator{
		Fees: []Fee{
			{
				ComputationType: FlatFeeType,
				FeeType:         "EMTL",
				FlatAmount:      50,
				Min:             int64Ptr(10000),
				Max:             nil,
			},
			{
				ComputationType: PercentageFeeType,
				FeeType:         "PROCESSING_FEE",
				PercentageRate:  1.5,
				MinFee:          int64Ptr(10),
				MaxFee:          int64Ptr(1000),
			},
		},
	}

	// Case 1: Transaction of 5,000 (below EMTL threshold of 10,000)
	// EMTL should NOT apply.
	// Processing fee: 1.5% of 5,000 = 75.
	fees, err := CalculateFee(ctx, service, 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fees) != 1 {
		t.Fatalf("expected 1 fee, got %d", len(fees))
	}
	if fees[0].FeeType != "PROCESSING_FEE" || fees[0].Fee != 75 {
		t.Errorf("expected PROCESSING_FEE 75, got %s %d", fees[0].FeeType, fees[0].Fee)
	}

	// Case 2: Transaction of 20,000 (meets EMTL threshold)
	// EMTL should apply (50).
	// Processing fee: 1.5% of 20,000 = 300.
	fees, err = CalculateFee(ctx, service, 20000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fees) != 2 {
		t.Fatalf("expected 2 fees, got %d", len(fees))
	}
	if fees[0].FeeType != "EMTL" || fees[0].Fee != 50 {
		t.Errorf("expected EMTL 50, got %s %d", fees[0].FeeType, fees[0].Fee)
	}
	if fees[1].FeeType != "PROCESSING_FEE" || fees[1].Fee != 300 {
		t.Errorf("expected PROCESSING_FEE 300, got %s %d", fees[1].FeeType, fees[1].Fee)
	}

	// Case 3: Small transaction triggering MinFee cap
	// Transaction of 100: 1.5% of 100 = 2 -> clamped to MinFee of 10.
	fees, err = CalculateFee(ctx, service, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fees) != 1 {
		t.Fatalf("expected 1 fee, got %d", len(fees))
	}
	if fees[0].Fee != 10 {
		t.Errorf("expected clamped fee 10, got %d", fees[0].Fee)
	}

	// Case 4: Large transaction triggering MaxFee cap
	// Transaction of 200,000: 1.5% of 200,000 = 3,000 -> clamped to MaxFee of 1000.
	fees, err = CalculateFee(ctx, service, 200000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fees) != 2 {
		t.Fatalf("expected 2 fees, got %d", len(fees))
	}
	if fees[1].Fee != 1000 {
		t.Errorf("expected clamped fee 1000, got %d", fees[1].Fee)
	}
}
