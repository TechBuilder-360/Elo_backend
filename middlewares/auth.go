package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Toflex/directory_v2/database/database"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "user"
)

func authJWT(a authentication.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.LoggerInContext(c)

		authHeader := c.GetHeader("Authorization")
		token, err := ExtractBearerToken(authHeader)
		if err != nil {
			logger.WithError(err).Error("failed to extract token from Authorization header")
		}

		// validate token
		userID, valid := a.VerifyJWT(c, token)
		if !valid || userID == "" {
			logger.Error("unable to validate jwt token")
		}

		// fetch user by ID
		usr, err := database.DBInstance().User.Query().Where(user.IDEQ(userID)).First(c)
		if err != nil {
			logger.WithError(err).WithField("user_id", userID).Error("failed to fetch user")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// add authenticated user to request context
		ctx := context.WithValue(c.Request.Context(), UserContextKey, usr)

		c.Request = c.Request.WithContext(ctx)
		// Pre-handler phase
		c.Next()
		// Post-handler phase

	}
}

func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization format")
	}

	return parts[1], nil
}

func UserFromContext(ctx context.Context) (*ent.User, error) {
	usr := ctx.Value(UserContextKey)

	user, ok := usr.(*ent.User)
	if !ok {
		return nil, errors.New(errors.ErrFailed, "user not found")
	}

	return user, nil
}
