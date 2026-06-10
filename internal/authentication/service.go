package authentication

import (
	"context"
	"github.com/Toflex/directory_v2/database/redis"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/samber/do/v2"
)

type IService interface {
	RegisterUser(ctx context.Context, body Registration, log log.Entry) (*string, error)
	RequestOTP(ctx context.Context, payload OTPRequest, log log.Entry) (*string, error)
	Login(ctx context.Context, payload Login, log log.Entry) (*model.LoginResponse, error)
}

type Service struct {
	repo  IRepository
	cache *redis.Client
}

func NewService(i do.Injector) IService {
	db := do.MustInvoke[*ent.Client](i)
	return &Service{
		repo:  Newrepository(db),
		cache: do.MustInvoke[*redis.Client](i),
	}
}

// 	if len(uid) == 0 {
// 		return errors.New("token has expired")
// 	}

// 	user, err := a.userRepo.GetUserByID(uid)
// 	if err != nil || user == nil {
// 		return errors.New("account not found")
// 	}

// 	if user.EmailVerified {
// 		return nil
// 	}

// 	user.EmailVerified = true
// 	user.Active = true
// 	user.EmailVerifiedAt = time.Now()

// 	if err = d.userRepo.Update(user); err != nil {
// 		logger.Error("An Error occurred while Activating your account, Please try again. %s", err.Error())
// 		return errors.New("account activation failed")
// 	}

// 	err = d.redis.Delete(token)
// 	if err != nil {
// 		logger.Error(err.Error())
// 	}

// 	return nil
// }

// // Login
// // Handles authentication logic
// func (a *authenticationService) Login(body *types.AuthRequest) (*types.LoginResponse, error) {
// 	response := new(types.LoginResponse)

// 	user, err := d.userRepo.GetByEmail(util.ToLower(body.EmailAddress))
// 	if err != nil {
// 		log.Error("An error occurred when fetching user profile. %s", err.Error())
// 		return nil, errors.New(constant.InternalServerError)
// 	}
// 	if user == nil {
// 		return nil, errors.New("account not found")
// 	}

// 	if !user.Active {
// 		return nil, errors.New("account is inactive")
// 	}

// 	// Validate OTP token
// 	token, err := d.repo.GetToken(user.ID)
// 	if err != nil {
// 		log.Error("An Error occurred when validating login token. %s", err.Error())
// 		return nil, errors.New("token validation failed")
// 	}

// 	if token == nil {
// 		return nil, errors.New("invalid OTP")
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(util.AddToStr(token)), []byte(body.Otp))
// 	if err != nil {
// 		log.Error("An Error occurred when comparing login token. %s", err.Error())
// 		return nil, errors.New("invalid OTP")
// 	}

// 	// Generate JWT for user
// 	tk, err := d.generateJWT(user.ID)
// 	if err != nil {
// 		log.Error("An error occurred when generating jwt token. %s", err.Error())
// 		return nil, errors.New("request failed")
// 	}

// 	err = d.repo.DeleteToken(user.ID)
// 	if err != nil {
// 		log.Error("an error occurred when removing jwt token. %s", err.Error())
// 	}

// 	profile := types.UserProfile{
// 		ID:            user.ID,
// 		FirstName:     user.FirstName,
// 		LastName:      user.LastName,
// 		DisplayName:   user.DisplayName,
// 		EmailAddress:  user.EmailAddress,
// 		PhoneNumber:   user.PhoneNumber,
// 		EmailVerified: user.EmailVerified,
// 		LastLogin:     user.LastLogin,
// 	}

// 	response.Profile = profile
// 	response.Authentication = *tk

// 	defer func() {
// 		user.LastLogin = time.Now()
// 		d.userRepo.Update(user)
// 	}()

// 	return response, nil
// }

// func (a *authenticationService) RequestToken(body *types.EmailRequest, logger log.Entry) error {
// 	if !util.ValidateEmail(body.EmailAddress) {
// 		return errors.New("account not found")
// 	}
// 	email := strings.ToLower(body.EmailAddress)

// 	// Check if email address exist
// 	user, err := d.userRepo.GetByEmail(email)
// 	if err != nil {
// 		logger.Error(err.Error())
// 		return errors.New("request failed")
// 	}

// 	if user == nil {
// 		logger.Error(err.Error())
// 		return errors.New("user not found")
// 	}

// 	duration := uint(5)
// 	otp := "123456"

// 	if configs.IsProduction() {
// 		otp = util.GenerateNumericToken(6)

// 		mailTemplate := &sendgrid.OTPMailRequest{
// 			Code:     otp,
// 			ToMail:   user.EmailAddress,
// 			ToName:   user.LastName + " " + user.FirstName,
// 			Name:     user.DisplayName,
// 			Duration: duration,
// 		}
// 		err = sendgrid.SendOTPMail(mailTemplate)
// 		if err != nil {
// 			logger.Error("Error occurred when sending otp email. %s", err.Error())
// 		}
// 	}

// 	hashedOTP, _ := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
// 	err = d.repo.StoreToken(user.ID, string(hashedOTP), duration)
// 	if err != nil {
// 		logger.Error("Error occurred when sending token %s", err)
// 		return errors.New("request failed please try again")
// 	}

// 	return nil
// }

// type AuthToken struct {
// 	Token        string
// 	RefreshToken string
// }

// func (a *authenticationService) ValidateToken(encodedToken string) (*jwt.RegisteredClaims, error) {
// 	claims := &jwt.RegisteredClaims{}
// 	tkn, err := jwt.ParseWithClaims(encodedToken, claims, func(token *jwt.Token) (any, error) {
// 		authToken, err := d.repo.GetToken(fmt.Sprintf("auth::%s", claims.ID))
// 		if err != nil {
// 			return nil, err
// 		}

// 		if err != nil {
// 			return nil, err
// 		}

// 		tk := AuthToken{}

// 		err = json.Unmarshal([]byte(util.AddToStr(authToken)), &tk)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if tk.Token != encodedToken {
// 			return nil, errors.New("jwt token not found")
// 		}

// 		key := fmt.Sprintf("%s-%s", configs.Instance.Secret, tk.RefreshToken)
// 		return []byte(key), nil
// 	})

// 	if err != nil || !tkn.Valid {
// 		return nil, err
// 	}

// 	return claims, nil
// }

// func (a *authenticationService) RefreshUserToken(body *types.RefreshTokenRequest, logger log.Entry) (*types.Authentication, error) {
// 	claims, err := d.ValidateToken(body.Token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response *types.Authentication
// 	response, err = d.generateJWT(claims.ID)
// 	if err != nil {
// 		return nil, errors.New("token could not be generated")
// 	}

// 	return response, nil
// }

// func (a *authenticationService) Logout(Token string) error {
// 	claims, err := d.ValidateToken(Token)
// 	if err != nil {
// 		return err
// 	}

// 	// Invalidate Refresh token
// 	return d.repo.DeleteToken(fmt.Sprintf("auth::%s", claims.ID))
// }
