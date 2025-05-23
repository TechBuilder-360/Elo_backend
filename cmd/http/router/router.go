package router

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(engine *gin.Engine) {
	initializeGeneralRoutes(engine)
}
