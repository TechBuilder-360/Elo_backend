package runtime

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

func server(i do.Injector) (*gin.Engine, error) {
	engine := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	// gin
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	return engine, nil
}
