package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/internal/email"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/redis/go-redis/v9"
)

func (s *Service) Login(ctx context.Context, payload Login, log log.Entry) (*model.LoginResponse, error) {
	cacheOTPString, err := s.cache.Get(ctx, payload.Identifier)
	if err != nil {
		log.WithError(err).Error("unable to fetch stored OTP")
		if err == redis.Nil {
			return nil, errors.New(errors.ErrFailed, "Expired OTP")
		}
		return nil, errors.New(errors.ErrNotFound, "request failed")
	}

	data := CacheOTP{}

	err = json.Unmarshal([]byte(cacheOTPString), &data)
	if err != nil {
		return nil, errors.New(errors.ErrNotFound, "request failed")
	}

	// verify otp
	if !util.DecryptPassword(payload.Otp, data.EncryptedOTP) {
		data.Attempts += 1
		if data.Attempts < 3 {
			b, _ := json.Marshal(data)
			err = s.cache.Set(ctx, payload.Identifier, string(b), time.Duration(time.Minute*5))
			if err != nil {
				log.WithError(err).Error("failed to add OTP token to cache")
				return nil, errors.New(errors.ErrFailed, "request failed")
			}
		} else {
			s.cache.Delete(ctx, payload.Identifier)
		}
	}

	user, err := s.repo.GetUserByID(ctx, data.UserID)
	if err != nil {
		log.WithField("User ID", data.UserID).WithError(err).Error("login failed")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if user.Disabled {
		log.WithFields(map[string]interface{}{"ID": user.ID, "disabled": user.Disabled, "reason": user.DisableReason}).Error("login failed user is disabled")
		return nil, errors.New(errors.ErrFailed, "account is disabled")
	}

	return s.generateJWT(ctx, data.UserID)
}

func (s *Service) RequestOTP(ctx context.Context, payload OTPRequest, log log.Entry) (*string, error) {
	user, err := s.repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
		return nil, errors.New(errors.ErrInternal, "request failed")
	}

	if user == nil {
		log.Error("user with email '%s' not found", payload.Email)
		return nil, errors.New(errors.ErrNotFound, "user not found!")
	}

	otp := "111111"
	if configuration.IsProduction() {
		otp = util.GenerateOTP()
	}

	reference := util.GenerateReference(configuration.Instance.AppName)
	encryptedOTP := util.EncryptPassword(otp)

	data := CacheOTP{
		UserID:       user.ID,
		EncryptedOTP: encryptedOTP,
		Attempts:     0,
	}

	b, _ := json.Marshal(data)
	err = s.cache.Set(ctx, reference, string(b), time.Duration(time.Minute*5))
	if err != nil {
		log.WithError(err).Error("failed to add OTP token to cache")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	// send OTP email
	if configuration.IsProduction() {
		err = email.NewOTPTask(email.OTPMailRequest{
			ToMail: user.EmailAddress,
			ToName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			OTP:    otp,
		})
		if err != nil {
			log.WithError(err).Error("Error occurred when sending activation otp email")
		}
	}

	return &reference, nil
}
