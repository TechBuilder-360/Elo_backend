package verification

import (
	"net/http"

	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/gin-gonic/gin"
)

type VerificationController struct {
	service *Service
}

func NewVerificationController(service *Service) *VerificationController {
	return &VerificationController{
		service: service,
	}
}

func (vc *VerificationController) RegisterRoutes(engine *gin.Engine) {
	verification := engine.Group("verification")
	v1 := verification.Group("v1")
	v1.GET("/verification/{expiration_time}", vc.initializeVerification)
}

// initializeVerification handles the verification request processing
func (vc *VerificationController) initializeVerification(ctx *gin.Context) {
	logger := log.LoggerInContext(ctx)
	logger.Info("Initialize verification request")

	var payload VerificationResult
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.WithError(err).Error("failed to bind verification payload")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := vc.service.ProcessVerification(ctx.Request.Context(), &payload); err != nil {
		logger.WithError(err).Error("failed to process verification request")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "verified"})
}
