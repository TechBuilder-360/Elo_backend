package provider

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/pkg/log"
)

type IService interface {
	GetServiceByIdentifier(ctx context.Context, identifier string) (ServiceLocator, error)
}

type service struct {
	db   *ent.Client
	repo IRepository
}

func NewService() IService {
	db := database.DBInstance()
	return &service{
		db:   db,
		repo: Newrepository(db),
	}
}

func (s *service) GetServiceByIdentifier(ctx context.Context, identifier string) (ServiceLocator, error) {
	logger := log.LoggerInContext(ctx)
	serviceLocator := ServiceLocator{}

	// get service by identifier
	serv, err := s.repo.FindByIdentifier(ctx, identifier)
	if err != nil {
		logger.WithField("identifier", identifier).WithError(err).Error("failed to find service by identifier")
		return serviceLocator, err
	}

	// Get active provider
	provider, err := s.repo.GetProvider(ctx, serv.Provider)
	if err != nil {
		logger.WithField("Provider", serv.Provider).WithError(err).Error("failed to get provider")
		return serviceLocator, err
	}

	fees := make([]Fee, 0)

	if len(serv.Edges.FeeRules) > 0 {
		for _, feeRule := range serv.Edges.FeeRules {
			fees = append(fees, Fee{
				ComputationType: FeeComputationType(feeRule.CalculationMethod),
				FeeType:         string(feeRule.FeeType),
				FlatAmount:      feeRule.FlatAmount,
				PercentageRate:  feeRule.PercentageRate,
				Min:             feeRule.MinValue,
				Max:             feeRule.MaxValue,
				MinFee:          feeRule.MinFee,
				MaxFee:          feeRule.MaxFee,
				Tiered: func() []Tier {
					tiers := make([]Tier, 0)
					if len(feeRule.Edges.Tiers) > 0 {
						for _, t := range feeRule.Edges.Tiers {
							tiers = append(tiers, Tier{
								From:            (t.MinAmount),
								To:              (t.MaxAmount),
								ComputationType: FeeComputationType(t.ChargeType),
								FlatAmount:      t.FlatAmount,
								PercentageRate:  t.PercentageRate,
							})
						}
					}
					return tiers
				}(),
			})
		}
	}

	serviceLocator = ServiceLocator{
		Name:       serv.Name,
		Identifier: serv.Identifier,
		ActiveProvider: ActiveProvider{
			Name:   provider.Name,
			Slug:   provider.Slug,
			Active: provider.Active,
		},
		Fees:     fees,
		MinValue: int64(serv.MinValue),
		MaxValue: int64(serv.MaxValue),
	}

	return serviceLocator, nil
}

func CalculateFee(ctx context.Context, service ServiceLocator, amount int64) ([]ServiceFee, error) {
	serviceFee := make([]ServiceFee, 0)
	fees := service.Fees
	for _, fee := range fees {
		if !fee.AppliesTo(amount) {
			continue
		}

		switch fee.ComputationType /* calculation_method */ {
		case FlatFeeType:
			computedFee := fee.ClampFee(fee.FlatAmount)
			serviceFee = append(serviceFee, ServiceFee{
				FeeType: fee.FeeType,
				Fee:     computedFee,
			})

		case PercentageFeeType:
			rawFee := int64(math.Round(float64(amount) * (fee.PercentageRate / 100)))
			computedFee := fee.ClampFee(rawFee)
			serviceFee = append(serviceFee, ServiceFee{
				FeeType: (fee.FeeType),
				Fee:     computedFee,
				Rate:    &fee.PercentageRate,
			})

		case TieredFeeType:
			tfee, err := computeTieredFee(fee.Tiered, amount)
			if err != nil {
				return nil, err
			}
			if tfee == nil {
				return nil, errors.New("no matching tier found for amount")
			}

			tfee.Fee = fee.ClampFee(tfee.Fee)
			tfee.FeeType = fee.FeeType

			serviceFee = append(serviceFee, *tfee)
		}
	}

	return serviceFee, nil
}

// Matches reports whether amount falls within this tier's [From, To] range.
func (t Tier) matches(amount int64) bool {
	if amount < t.From {
		return false
	}
	if t.To != nil && amount > *t.To {
		return false
	}
	return true
}

// ComputeFee returns the fee charged for the given amount under this tier.
func (t Tier) computeFee(amount int64) (*ServiceFee, error) {
	switch t.ComputationType {
	case FlatFeeType:
		return &ServiceFee{
			FeeType: string(t.ComputationType),
			Fee:     t.FlatAmount,
		}, nil

	case PercentageFeeType:
		return &ServiceFee{
			FeeType: string(t.ComputationType),
			Fee:     int64(math.Round(float64(amount) * (t.PercentageRate / 100))),
			Rate:    &t.PercentageRate,
		}, nil

	default:
		return nil, fmt.Errorf("unknown fee type: %q", t.ComputationType)
	}
}

// FindTier returns the tier matching amount, or nil if none matches.
func findTier(tiers []Tier, amount int64) (*Tier, error) {
	for i := range tiers {
		if tiers[i].matches(amount) {
			return &tiers[i], nil
		}
	}
	return nil, fmt.Errorf("no tier matches amount %d", amount)
}

// ComputeTieredFee finds the matching tier for amount and computes the fee.
func computeTieredFee(tiers []Tier, amount int64) (*ServiceFee, error) {
	tier, err := findTier(tiers, amount)
	if err != nil {
		return nil, err
	}
	fee, err := tier.computeFee(amount)
	if err != nil {
		return nil, err
	}
	return fee, nil
}

func SumServiceFees(fees []ServiceFee) int64 {
	total := int64(0)
	for _, fee := range fees {
		total += fee.Fee
	}
	return total
}
