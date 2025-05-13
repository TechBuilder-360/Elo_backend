package runtime

import "github.com/gin-gonic/gin"

func server() *gin.Engine {
	engine := gin.Default()

	// gin
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	return engine
}
