package runtime

import (
	"github.com/Toflex/directory_v2/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func server(i do.Injector) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// gin
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(middlewares.Session())
	engine.Use(middlewares.CORS())

	return engine, nil
}
