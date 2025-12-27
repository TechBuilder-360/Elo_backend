package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Toflex/directory_v2/domain/pkg"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/provider"
	"github.com/Toflex/directory_v2/pkg/saferoutine"
	"github.com/Toflex/directory_v2/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

func (s *Service) RegisterUser(ctx context.Context, body Registration, log log.Entry) (*string, error) {
	err := body.Validate()
	if err != nil {
		return nil, err
	}

	// Check if email address exist
	existingUser, err := s.repo.GetUserByEmail(ctx, body.EmailAddress)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New(errors.ErrFailed, err.Error())
	}

	if existingUser != nil {
		log.Info("Email address already registered. '%s'", body.EmailAddress)
		return nil, errors.New(errors.ErrFailed, "account already exist")
	}

	onboarding := Onboarding{
		DisplayName:  body.DisplayName,
		EmailAddress: body.EmailAddress,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		PhoneNumber:  body.PhoneNumber,
		Avatar:       body.Avatar,
	}

	uid, err := s.repo.Create(ctx, onboarding)
	if err != nil {
		log.Error("error: occurred when saving new user. %s", err.Error())
		return nil, &errors.CustomError{
			Message: "registration was not successful",
		}
	}

	// send welcome email
	saferoutine.Run(func() {
		m, ok := provider.GetImpl(constant.SendGrid.ToString())
		if !ok {
			log.Error("provider not found")
			return
		}
		mail, ok := provider.ConformsTo[pkg.IEmailProvider](m)
		if !ok {
			log.Error("provider not supported for Email service")
			return
		}

		// Send Wellcome email
		mailData := &pkg.WelcomeMailRequest{
			ToMail:   body.EmailAddress,
			ToName:   fmt.Sprintf("%s %s", body.FirstName, body.LastName),
			FullName: fmt.Sprintf("%s %s", body.FirstName, body.LastName),
		}

		err = mail.SendWelcomeMail(*mailData)
		if err != nil {
			log.Error("Error occurred when sending activation email. %s", err.Error())
		}
	})

	return uid, nil
}

func (s *Service) RequestOTP(ctx context.Context, email string, log log.Entry) (*string, error) {
	email = strings.ToLower(email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to fetch user")
		return nil, errors.New(errors.ErrInternal, "request failed")
	}

	if user == nil {
		log.Error("user with email '%s' not found", email)
		return nil, errors.New(errors.ErrNotFound, "user not found!")
	}

	otp := "111111"
	if configuration.IsProduction() {
		otp = utils.GenerateOTP()
	}

	reference := utils.GenerateReference(configuration.Instance.AppName)
	encryptedOTP := utils.Encrypt(otp)

	data := CacheOTP{
		UserID:       user.ID,
		EncryptedOTP: encryptedOTP,
		Attempts:     0,
	}

	b, _ := json.Marshal(data)
	err = s.cache.Set(ctx, reference, string(b), time.Duration(time.Minute*5))
	if err != nil {
		log.Error("failed to add OTP token to cache %s", err.Error())
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	// send OTP email
	if configuration.IsProduction() {
		m, ok := provider.GetImpl(constant.SendGrid.ToString())
		if !ok {
			log.Error("email provider not found!")
			return nil, errors.New(errors.ErrFailed, "request failed")
		}
		mail, ok := provider.ConformsTo[pkg.IEmailProvider](m)
		if !ok {
			log.Error("provider not supported for Email service")
			return nil, errors.New(errors.ErrFailed, "request failed")
		}

		// Send Wellcome email
		mailData := &pkg.OTPMailRequest{
			ToMail: user.EmailAddress,
			ToName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			OTP:    otp,
		}

		err = mail.SendOTPMail(*mailData)
		if err != nil {
			log.Error("Error occurred when sending activation email. %s", err.Error())
		}
	}

	return &reference, nil
}

func (s *Service) Login(ctx context.Context, payload Login, log log.Entry) (*model.LoginResponse, error) {
	cacheOTPString, err := s.cache.Get(ctx, payload.Identifier)
	if err != nil {
		log.Error("unable to fetch stored OTP %s", err.Error())
		if err == redis.Nil {
			return nil, errors.New(errors.ErrFailed, "Expired OTP")
		}
		return nil, errors.New(errors.ErrNotFound, "request failed")
	}

	data := CacheOTP{}

	json.Unmarshal([]byte(cacheOTPString), &data)

	if !utils.Decrypt(payload.Otp, data.EncryptedOTP) {
		data.Attempts += 1
		if data.Attempts < 3 {
			b, _ := json.Marshal(data)
			err = s.cache.Set(ctx, payload.Identifier, string(b), time.Duration(time.Minute*5))
			if err != nil {
				log.Error("failed to add OTP token to cache %s", err.Error())
				return nil, errors.New(errors.ErrFailed, "request failed")
			}
		} else {
			s.cache.Delete(ctx, payload.Identifier)
		}
	}

	user, err := s.repo.GetUserByID(ctx, data.UserID)
	if err != nil {
		log.Error("login failed to fetch user by ID %s", err.Error())
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	if user.Disabled {
		log.Error("login failed user is disabled")
		return nil, errors.New(errors.ErrFailed, "account is disabled")
	}

	return s.generateJWT(ctx, data.UserID)
}

func (s *Service) generateJWT(ctx context.Context, userId string) (*model.LoginResponse, error) {
	issuedAt := time.Now()
	expireAt := issuedAt.Add(time.Hour * 24)
	claims := jwt.RegisteredClaims{
		Issuer:    configuration.Instance.AppName,
		ExpiresAt: jwt.NewNumericDate(expireAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        utils.GenerateUUID(),
	}

	//encoded string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(configuration.Instance.JWTSecret))
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("token could not be generated %s", err.Error())
	}

	jti := JWTToken{
		UserID: userId,
	}

	// Store token to enable revoking token 30 Days
	marshal, err := json.Marshal(jti)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(ctx, fmt.Sprintf("auth::%s", claims.ID), string(marshal), time.Duration(time.Hour*24*30))
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
		ExpireAt:    expireAt.Unix(),
	}, nil
}
