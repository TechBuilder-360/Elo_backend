package business

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Toflex/directory_v2/ent"
	kyb "github.com/Toflex/directory_v2/pkg/business"
	"github.com/Toflex/directory_v2/pkg/constant"
	"github.com/Toflex/directory_v2/pkg/errors"
	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (s *service) AddDocument(ctx context.Context, payload UploadDocumentRequest, logger log.Entry) error {
	if !payload.User.Verified {
		return errors.New(errors.ErrFailed, "user not verified, please present a valid means of identification")
	}

	err := util.ValidateDataURI(payload.File, []string{string(constant.PDF), string(constant.JPG), string(constant.JPEG)})
	if err != nil {
		logger.WithError(err).Error("failed to validate data uri")
		return errors.New(errors.ErrValidation, "document is invalid")
	}

	if payload.DocumentID == "" {
		logger.WithField("user", payload.User.ID).WithField("business", payload.Business.ID).Error("document id is required")
		return errors.New(errors.ErrFailed, "Document ID is required")
	}

	kybdoc, err := s.repo.KYBDocumentByID(ctx, payload.DocumentID)
	if err != nil || kybdoc == nil {
		logger.WithError(err).WithField("user", payload.User.ID).WithField("business", payload.Business.ID).WithField("document", payload.DocumentID).Error("failed to fetch document")
		return errors.New(errors.ErrFailed, "document type not found")
	}

	if kybdoc.Name != string(kyb.OtherDocument) {
		payload.Description = kybdoc.Name
	}

	var doc *ent.BusinessDocument

	if kybdoc.Name != string(kyb.OtherDocument) {
		doc, err = s.repo.GetBusinessDocument(ctx, payload.Business, payload.DocumentID)
		if err != nil {
			logger.WithError(err).WithField("user", payload.User.ID).WithField("business", payload.Business.ID).WithField("document", payload.DocumentID).Error("failed to fetch document")
			return errors.New(errors.ErrFailed, "request failed")
		}
	}

	url, err := s.cloudinary.UploadBusinessFile(ctx, payload.File, payload.Description, payload.Business.Name, "KYB")
	if err != nil {
		logger.WithError(err).
			WithField("business", payload.Business.Name).
			WithField("registered by", payload.Business.RegisteredBy).
			WithField("document", kybdoc.Name).
			Error("failed to upload document")
		return errors.New(errors.ErrFailed, "failed to upload document")
	}

	if doc == nil {
		bizDoc := AddBusinessDocument{
			Title:       kybdoc.Name,
			Description: payload.Description,
			URL:         url,
			Type:        "KYB",
			DocumentID:  kybdoc.ID,
		}

		// Add the document to the business
		err = s.repo.AddDocument(ctx, payload.Business, bizDoc)
		if err != nil {
			logger.WithError(err).
				WithField("business", payload.Business.Name).
				WithField("registered by", payload.Business.RegisteredBy).
				WithField("document", payload.Description).
				Error("failed to add business document")
			return errors.New(errors.ErrFailed, "failed to add business document")
		}
	} else {
		err = s.repo.UpdateDcument(ctx, payload.Business.ID, doc.ID, url)
		if err != nil {
			logger.WithError(err).
				WithField("business", payload.Business.Name).
				WithField("registered by", payload.Business.RegisteredBy).
				WithField("document", payload.Description).
				Error("failed to update business document")
			return errors.New(errors.ErrFailed, "failed to add business document")
		}
	}

	return nil
}

func (s *service) DeleteDocument(ctx context.Context, business *ent.Business, documentID string, logger log.Entry) error {
	err := s.repo.DeleteDocument(ctx, documentID)
	if err != nil {
		logger.WithError(err).Error("failed to delete business document")
		return errors.New(errors.ErrFailed, "request failed")
	}

	return nil
}

func (s *service) BusinessDetail(ctx context.Context, user *ent.User, business *ent.Business, payload *BusinessDetailRequest, logger log.Entry) error {
	// Add random string to business name until they gets verified
	if payload.Name != nil {
		name := fmt.Sprintf("%s_%s", util.AddressToString(payload.Name), strconv.Itoa(util.RandomInt()))
		payload.Name = &name
	}

	err := s.repo.Update(ctx, business, UpdateBusiness{
		Name:                   payload.Name,
		About:                  payload.About,
		Industry:               payload.Industry,
		Website:                payload.Website,
		Number:                 &payload.RegistrationDetail.Number,
		CountryOfIncorporation: &payload.RegistrationDetail.CountryOfIncorporation,
		DateOfIncorporation:    &payload.RegistrationDetail.DateOfIncorporation,
	})
	if err != nil {
		logger.WithError(err).WithField("business", business.ID).Error("failed to update business detail")
		return errors.New(errors.ErrFailed, "request failed")
	}

	return nil
}

func (s *service) GetDocuments(ctx context.Context) ([]DocumentResult, error) {
	logger := log.LoggerInContext(ctx)

	result, err := s.repo.KYBDOcuments(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch kyb documents")
		return nil, errors.New(errors.ErrFailed, "request failed")
	}

	documents := make([]DocumentResult, 0)
	for _, doc := range result {
		documents = append(documents, DocumentResult{
			ID:       doc.ID,
			Name:     doc.Name,
			Required: doc.Required,
		})
	}

	return documents, nil
}

func (s *service) GetKYBDocuments(ctx context.Context, b *ent.Business) ([]KYBDocument, error) {
	logger := log.LoggerInContext(ctx)
	data, err := s.repo.GetKYBDocuments(ctx, b.ID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch business KYB documents")
		return nil, errors.New(errors.ErrFailed, "failed document to fetch documents")
	}

	docs := make([]KYBDocument, 0)

	for _, doc := range data {
		var kybDocumentID string
		if len(doc.Edges.KybDocument) > 0 {
			kybDocumentID = doc.Edges.KybDocument[0].ID
		}

		docs = append(docs, KYBDocument{
			ID:          doc.ID,
			DocumentID:  kybDocumentID,
			Description: doc.Description,
			File:        doc.URL,
		})
	}

	return docs, nil
}
