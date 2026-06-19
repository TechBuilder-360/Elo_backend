package verification

import (
	"net/http"

	"github.com/Toflex/directory_v2/middlewares"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type IVerificationController interface {
	RegisterRoutes(engine *gin.Engine)
}

type VerificationController struct {
	service *Service
}

func NewVerificationController(service *Service) IVerificationController {
	return &VerificationController{
		service: service,
	}
}

func (vc *VerificationController) RegisterRoutes(engine *gin.Engine) {
	verification := engine.Group("v1/verification")
	verification.Use(middlewares.MiddleWare().Jwt)

	// v1.GET("/{expiration_time}", vc.initialiseVerification)
	verification.GET("/user/request/link", vc.requestUserVerification)
	verification.GET("/:reference_id", vc.redirectVerification)
}

// initialiseVerification handles the verification request processing
func (vc *VerificationController) requestUserVerification(ctx *gin.Context) {
	logger := log.LoggerInContext(ctx)
	logger.Info("request user verification request")

	user, err := middlewares.UserFromContext(ctx.Request.Context())
	if err != nil {
		logger.WithError(err).Error("failed to fetch user in context")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			response.Response{
				Status:  false,
				Message: err.Error(),
			})
		return
	}

	resp, err := vc.service.RequestVerificationLink(ctx, &VerificationRequest{
		Entity: constant.UserEntityType,
		ID:     user.ID,
	})
	if err != nil {
		logger.WithError(err).Error("failed to process verification request")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.Response{
				Status:  false,
				Message: err.Error(),
			})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK,
		response.Response{
			Status:  true,
			Message: "Success",
			Data:    resp,
		})
}

// redirectVerification handles the verification request processing
func (vc *VerificationController) redirectVerification(ctx *gin.Context) {
	logger := log.LoggerInContext(ctx)
	logger.Info("request redirect verification request")

	reference := ctx.Param("reference_id")
	link, err := vc.service.GetProviderLink(ctx, reference)
	if err != nil {
		logger.WithError(err).Error("failed to process verification link")
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			response.Response{
				Status:  false,
				Message: err.Error(),
			})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
		Status:  true,
		Data:    link,
		Message: "Success",
	})

	return
}
