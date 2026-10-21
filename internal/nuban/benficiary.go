package nuban

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
)

const blacklistToken = "nuban::blacklist"

func (s *service) GetBeneficiaryByAccountNumber(ctx context.Context, ownerId, accountNumber, bankID string) (*beneficiary, error) {
	logger := log.LoggerInContext(ctx)

	beneficiary, err := s.repo.GetBeneficiaryByAccountNumber(ctx, ownerId, accountNumber, bankID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch beneficiary")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if beneficiary == nil {
		return nil, nil
	}

	if beneficiary.Disabled {
		return nil, errors.New(errors.ErrFailed, "beneficiary is disabled")
	}

	return beneficiary, nil
}

func (s *service) CreateBeneficiary(ctx context.Context, payload *createBeneficiary) (*beneficiary, error) {
	logger := log.LoggerInContext(ctx)

	beneficiary, err := s.repo.CreateBeneficiary(ctx, payload)
	if err != nil {
		logger.WithError(err).Error("failed to save beneficiary")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	return beneficiary, nil
}

type nubanBeneficiaryResult struct {
	OwnerID       string
	ID            string
	AccountName   string
	AccountNumber string
	BankName      string
	BankID        string
	Disabled      bool
}

type nubanBeneficiaryPayload struct {
	OwnerID       string
	AccountName   string
	AccountNumber string
	BankID        string
}

func (s *service) IsBlacklisted(ctx context.Context, payload *nubanBeneficiaryPayload) (bool, error) {
	logger := log.LoggerInContext(ctx)

	key := fmt.Sprintf("%s::%s::%s", blacklistToken, payload.BankID, payload.AccountNumber)
	result, err := s.cache.Get(ctx, key)
	if err == nil {
		data := beneficiaryBlacklist{}
		err = json.Unmarshal([]byte(result), &data)
		if err == nil {
			return data.Disabled, nil
		}
	}

	response, err := s.repo.GetNubanBlacklist(ctx, payload.AccountNumber, payload.BankID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch nuban blackist")
		return true, errors.New(errors.ErrFailed, "failed to fetch blacklist")
	}

	if response == nil {
		return false, nil
	}

	responseByte, err := json.Marshal(response)
	if err == nil {
		// cache response
		s.cache.Set(ctx, key, string(responseByte), 0)
	}

	return response.Disabled, nil
}

func (s *service) beneficiaryCheck(ctx context.Context, payload *nubanBeneficiaryPayload) (*nubanBeneficiaryResult, error) {
	logger := log.LoggerInContext(ctx)

	// check that beneficiary isn't on blacklist
	blacklisted, err := s.IsBlacklisted(ctx, payload)
	if err != nil || blacklisted {
		logger.WithError(err).WithField("blacklisted", blacklisted).Error("%s", err.Error())
		return nil, err
	}

	beneficiary, err := s.GetBeneficiaryByAccountNumber(ctx, payload.OwnerID, payload.AccountNumber, payload.BankID)
	if err != nil {
		return nil, err
	}

	// if beneficiary is nil, save beneficiary
	if beneficiary == nil {
		bank, err := s.repo.GetBankByID(ctx, payload.BankID)
		if err != nil {
			logger.WithError(err).WithField("bank_id", payload.BankID).Error("failed to fetch bank with ID")
			return nil, errors.New(errors.ErrFailed, "bank not found")
		}

		beneficiary, err = s.repo.CreateBeneficiary(ctx, &createBeneficiary{
			AccountName:   payload.AccountName,
			AccountNumber: payload.AccountNumber,
			BankName:      bank.Name,
			BankID:        bank.ID,
			OwnerID:       payload.OwnerID,
			Nickname:      payload.AccountName,
		})
		if err != nil {
			logger.WithError(err).Error("failed to create beneficiary")
			return nil, errors.New(errors.ErrFailed, "request failed")
		}
	}

	if beneficiary.Disabled {
		logger.WithField("beneficiary", beneficiary).Error("beneficiary is disabled")
		return nil, errors.New(errors.ErrFailed, "beneficiary is disabled")
	}

	return &nubanBeneficiaryResult{
		OwnerID:       payload.OwnerID,
		ID:            beneficiary.ID,
		AccountName:   beneficiary.AccountName,
		AccountNumber: beneficiary.AccountNumber,
		Disabled:      beneficiary.Disabled,
		BankID:        beneficiary.BankID,
		BankName:      beneficiary.BankName,
	}, nil
}
