package dojah

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/types"
	"github.com/Toflex/directory_v2/pkg/verification"
	"github.com/gin-gonic/gin"
)

func (d *dojah) RegisterRoutes(engine *gin.Engine) {
	route := engine.Group("/webhook")

	route.POST("/dojah", d.handleDojahWebhook)
}

func (d *dojah) handleDojahWebhook(ctx *gin.Context) {
	logger := log.LoggerInContext(ctx)
	logger = logger.WithField("provider", "dojah")
	logger.Info("Process dojah webhook")

	h := sha256.Sum256([]byte(d.config.SecretKey))
	hash := hex.EncodeToString(h[:])

	if hash != ctx.GetHeader("x-dojah-signature-v2") {
		logger.Error("Failed to verify dojah webhook signature")
		ctx.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	payload := new(WebhookPayload)
	err := ctx.ShouldBindJSON(payload)
	if err != nil {
		logger.WithError(err).Error("Failed to bind dojah webhook body")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = d.processDojahWebhook(ctx, payload, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to process dojah webhook")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.AbortWithStatus(200)
}

func (d *dojah) processDojahWebhook(ctx context.Context, payload *WebhookPayload, logger log.Entry) error {
	var status types.VerificationStatus = constant.Pending

	switch payload.VerificationStatus {
	case "Completed":
		status = constant.Success
	case "Failed":
		status = constant.Failed
	default:
		return nil
	}

	var (
		identificationDetails = verification.VerificationResult{
			Provider:    constant.Dojah,
			ReferenceID: payload.ReferenceID,
			Metadata:    payload.Metadata,
			Status:      status,
			Message:     payload.Message,
			// NationalID:     verification.NationalID{},
			// Passport:       verification.Passport{},
			// BVN:            verification.BVN{},
			// DriversLicense: verification.DriversLicense{},
			// VoterID:        verification.VoterID{},
		}
		documentType = payload.Data.ID.Data.IDData.DocumentType
		idData       = payload.Data.ID.Data.IDData
		bvnDetails   = payload.Data.GovernmentData.Data.Bvn
	)

	switch documentType {
	case PASSPORT:
		{
			identificationDetails.Passport.FirstName = idData.FirstName
			identificationDetails.Passport.LastName = idData.LastName
			identificationDetails.Passport.Number = idData.DocumentNumber
			identificationDetails.Passport.Expiration = idData.ExpiryDate
			identificationDetails.Passport.IssueDate = idData.DateIssued
			identificationDetails.Passport.Nationality = idData.Nationality
			identificationDetails.Passport.DateOfBirth = idData.DateOfBirth
			identificationDetails.Passport.FrontImage = payload.Data.ID.Data.IDURL
			identificationDetails.Passport.BackImage = payload.Data.ID.Data.BackURL
		}
	case NATIONALID:
		{
			identificationDetails.NationalID.FirstName = idData.FirstName
			identificationDetails.NationalID.LastName = idData.LastName
			identificationDetails.NationalID.Number = idData.DocumentNumber
			identificationDetails.NationalID.Nationality = idData.Nationality
			identificationDetails.NationalID.DateOfBirth = idData.DateOfBirth
			identificationDetails.NationalID.FrontImage = payload.Data.ID.Data.IDURL
			identificationDetails.NationalID.BackImage = payload.Data.ID.Data.BackURL
		}
	case VOTERID:
		{
			identificationDetails.VoterID.FirstName = idData.FirstName
			identificationDetails.VoterID.LastName = idData.LastName
			identificationDetails.VoterID.Number = idData.DocumentNumber
			identificationDetails.VoterID.Expiration = idData.ExpiryDate
			identificationDetails.VoterID.IssueDate = idData.DateIssued
			identificationDetails.VoterID.Nationality = idData.Nationality
			identificationDetails.VoterID.DateOfBirth = idData.DateOfBirth
			identificationDetails.VoterID.FrontImage = payload.Data.ID.Data.IDURL
			identificationDetails.VoterID.BackImage = payload.Data.ID.Data.BackURL
		}
	case DRIVERLICENSE:
		{
			identificationDetails.DriversLicense.FirstName = idData.FirstName
			identificationDetails.DriversLicense.LastName = idData.LastName
			identificationDetails.DriversLicense.Number = idData.DocumentNumber
			identificationDetails.DriversLicense.Expiration = idData.ExpiryDate
			identificationDetails.DriversLicense.IssueDate = idData.DateIssued
			identificationDetails.DriversLicense.Nationality = idData.Nationality
			identificationDetails.DriversLicense.DateOfBirth = idData.DateOfBirth
			identificationDetails.DriversLicense.FrontImage = payload.Data.ID.Data.IDURL
			identificationDetails.DriversLicense.BackImage = payload.Data.ID.Data.BackURL
		}
	}

	// GovernmentID extract BVN for Nigerian users
	if bvnDetails.Entity.Bvn != "" {
		identificationDetails.BVN.FirstName = bvnDetails.Entity.FirstName
		identificationDetails.BVN.LastName = bvnDetails.Entity.LastName
		identificationDetails.BVN.MiddleName = bvnDetails.Entity.MiddleName
		identificationDetails.BVN.Number = bvnDetails.Entity.Bvn
		identificationDetails.BVN.PhoneNumber = bvnDetails.Entity.PhoneNumber1
		identificationDetails.BVN.Nationality = bvnDetails.Entity.Nationality
		identificationDetails.BVN.NationalIdentificationNumber = bvnDetails.Entity.Nin
		identificationDetails.BVN.DateOfBirth = bvnDetails.Entity.DateOfBirth
	}

	// if payload.Data.BusinessID.BusinessName != "" {

	// }

	return verification.QueueVerificationTask(identificationDetails)
}
