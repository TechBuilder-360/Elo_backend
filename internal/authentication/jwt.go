package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) generateJWT(ctx context.Context, userId string) (*model.LoginResponse, error) {
	logger := log.LoggerInContext(ctx)
	issuedAt := time.Now()
	expireAt := issuedAt.Add(time.Hour * 24)
	jti := util.GenerateCUID()
	sub := userId

	claims := jwt.MapClaims{
		"sub": sub,
		"exp": expireAt.Unix(),
		"iss": configuration.Instance.AppName,
		"iat": issuedAt.Unix(),
		"jti": jti,
	}

	//encoded string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(configuration.Instance.JWTSecret))
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New("token could not be generated %s", err.Error())
	}

	// Store token to enable revoking token 30 Days
	// marshal, err := json.Marshal(jtiEncoded)
	// if err != nil {
	// 	return nil, err
	// }

	err = s.cache.Set(ctx, fmt.Sprintf("auth::%s", userId), jti, time.Duration(time.Hour*time.Duration(configuration.Instance.TOKENLIFESPAN)))
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
		ExpireAt:    expireAt.Unix(),
	}, nil
}

func (s *Service) VerifyJWT(ctx context.Context, tokenString string) (string, bool) {
	logger := log.LoggerInContext(ctx)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		hmacSampleSecret := []byte(configuration.Instance.JWTSecret)
		return hmacSampleSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		logger.WithError(err).Error("failed to validate jwt")
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		jti := claims["jti"].(string)
		userId := claims["sub"].(string)

		v, err := s.cache.Get(ctx, fmt.Sprintf("auth::%s", userId))
		if err != nil {
			logger.WithError(err).Error("failed to validate jwt jti")
			return "", false
		}

		// err = json.Unmarshal([]byte(v), &jtiToken)
		// if err != nil {
		// 	logger.WithError(err).Error("failed to marshal jwt jti")
		// 	return "", false
		// }
		if jti != v {
			logger.WithField("jti", jti).WithField("user_id", userId).Error("failed to validate jwt jti")
			return "", false
		}

		return userId, true
	}

	return "", false
}
