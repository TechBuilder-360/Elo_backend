package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Toflex/directory_v2/graph/model"
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/golang-jwt/jwt/v4"
)

func (s *Service) generateJWT(ctx context.Context, userId string) (*model.LoginResponse, error) {
	logger := log.LoggerInContext(ctx)
	issuedAt := time.Now()
	expireAt := issuedAt.Add(time.Hour * 24)
	claims := jwt.RegisteredClaims{
		Issuer:    configuration.Instance.AppName,
		ExpiresAt: jwt.NewNumericDate(expireAt),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ID:        util.GenerateUUID(),
	}

	//encoded string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(configuration.Instance.JWTSecret))
	if err != nil {
		logger.Error(err.Error())
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
