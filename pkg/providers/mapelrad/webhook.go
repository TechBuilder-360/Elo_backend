package mapelrad

import "github.com/gin-gonic/gin"

func (m *maplerad) RegisterRoutes(engine *gin.Engine) {
	_ = engine.Group("/webhook")

	// route.POST("/dojah", d.handleDojahWebhook)
}
