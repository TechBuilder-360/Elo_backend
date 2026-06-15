package middlewares

import (
	"github.com/Toflex/directory_v2/pkg/configuration"
	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]bool{}

func isOriginAllowed(origin string) bool {
	// Allow localhost in development mode
	if !configuration.IsProduction() {
		return true
	}
	return allowedOrigins[origin]
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// origin := c.GetHeader("Origin")

		// if origin != "" && isOriginAllowed(origin) {
		// Set the critical Access-Control-Allow-Origin header
		// c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Business-Id")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// }

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}

		c.Next()
	}
}
