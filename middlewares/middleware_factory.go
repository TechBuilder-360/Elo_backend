package middlewares

import (
	"github.com/Toflex/directory_v2/internal/authentication"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

var registeredMiddleWare middleware

type middleware struct {
	Jwt gin.HandlerFunc
}

func InitializeMiddleWares(i do.Injector) {
	var (
		authenticationService = authentication.NewService(i)
	)

	registeredMiddleWare = middleware{
		Jwt: authJWT(authenticationService),
	}
}

func MiddleWare() *middleware {
	return &registeredMiddleWare
}
