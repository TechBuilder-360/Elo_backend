package hook

import (
	"context"
	"fmt"
	"sort"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/feerule"
	"github.com/Toflex/directory_v2/ent/feetier"
	"github.com/Toflex/directory_v2/ent/hook"
)

func validateTierRangesHook(next ent.Mutator) ent.Mutator {
	return hook.FeeTierFunc(func(ctx context.Context, m *ent.FeeTierMutation) (ent.Value, error) {
		ruleID, ok := m.FeeRuleID()
		if !ok {
			return next.Mutate(ctx, m)
		}

		existing, err := m.Client().FeeTier.Query().
			Where(feetier.HasFeeRuleWith(feerule.ID(ruleID))).
			All(ctx)
		if err != nil {
			return nil, err
		}

		if id, ok := m.ID(); ok {
			filtered := existing[:0]
			for _, t := range existing {
				if t.ID != id {
					filtered = append(filtered, t)
				}
			}
			existing = filtered
		}

		pending, err := buildPendingTier(ctx, m)
		if err != nil {
			return nil, err
		}

		if err := validateTierRanges(append(existing, pending)); err != nil {
			return nil, err
		}

		return next.Mutate(ctx, m)
	})
}

func buildPendingTier(ctx context.Context, m *ent.FeeTierMutation) (*ent.FeeTier, error) {
	t := &ent.FeeTier{}

	if id, ok := m.ID(); ok {
		t.ID = id
	}

	minAmount, ok := m.MinAmount()
	if !ok {
		old, err := m.OldMinAmount(ctx)
		if err != nil {
			return nil, fmt.Errorf("min_amount required: %w", err)
		}
		minAmount = old
	}
	t.MinAmount = minAmount

	if maxAmount, ok := m.MaxAmount(); ok {
		t.MaxAmount = &maxAmount
	} else if !m.MaxAmountCleared() {
		if old, err := m.OldMaxAmount(ctx); err == nil {
			t.MaxAmount = old
		}
	}

	return t, nil
}

func validateTierRanges(tiers []*ent.FeeTier) error {
	sorted := make([]*ent.FeeTier, len(tiers))
	copy(sorted, tiers)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MinAmount < sorted[j].MinAmount
	})

	openEndedCount := 0
	for i, t := range sorted {
		if t.MaxAmount != nil && *t.MaxAmount < t.MinAmount {
			return fmt.Errorf("tier %d: max_amount (%d) is less than min_amount (%d)", i, *t.MaxAmount, t.MinAmount)
		}
		if t.MaxAmount == nil {
			openEndedCount++
		}
		if i == 0 {
			continue
		}
		prev := sorted[i-1]
		if prev.MaxAmount == nil {
			return fmt.Errorf("tier %d overlaps: tier %d is open-ended but is not the last tier", i, i-1)
		}
		if t.MinAmount <= *prev.MaxAmount {
			return fmt.Errorf("tier %d (min=%d) overlaps tier %d (max=%d)", i, t.MinAmount, i-1, *prev.MaxAmount)
		}
		if t.MinAmount != *prev.MaxAmount+1 {
			return fmt.Errorf("gap between tier %d (max=%d) and tier %d (min=%d)", i-1, *prev.MaxAmount, i, t.MinAmount)
		}
	}

	if openEndedCount != 1 {
		return fmt.Errorf("expected exactly one open-ended (unbounded) tier, got %d", openEndedCount)
	}
	if sorted[len(sorted)-1].MaxAmount != nil {
		return fmt.Errorf("the open-ended tier must be the last tier (highest min_amount)")
	}
	if sorted[0].MinAmount != 0 {
		return fmt.Errorf("first tier must start at 0, got %d", sorted[0].MinAmount)
	}
	return nil
}
