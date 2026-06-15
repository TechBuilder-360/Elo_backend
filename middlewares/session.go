package middlewares

import (
	"net/http"

	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {

		appName := configuration.Instance.AppName
		sessionID := util.GenerateReference(appName)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     appName,
			Value:    sessionID,
			MaxAge:   1800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		c.Next()
	}
}
