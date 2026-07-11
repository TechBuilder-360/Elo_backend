package business

import (
	"context"

	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/business"
	"github.com/Toflex/directory_v2/ent/businessdocument"
	"github.com/Toflex/directory_v2/ent/kybdocument"
	"github.com/Toflex/directory_v2/ent/manager"
	"github.com/Toflex/directory_v2/ent/user"
	"github.com/Toflex/directory_v2/pkg/util"
)

func (r *repository) GetBusinessByName(ctx context.Context, name string) (*businessResult, error) {
	b, err := r.db.Business.Query().
		Where(business.NameEqualFold(name)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		return nil, err
	}

	return &businessResult{*b}, nil
}

func (r *repository) GetBusinessByID(ctx context.Context, id string) (*businessResult, error) {
	b, err := r.db.Business.Query().
		Where(business.IDEQ(id)).
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		return nil, err
	}

	return &businessResult{*b}, nil
}

func (r *repository) Create(ctx context.Context, payload createbusiness) (*ent.Business, error) {
	b, err := r.db.Business.
		Create().
		SetName(payload.Name).
		SetAbout(payload.About).
		SetCategory(payload.Category).
		SetEmail(payload.Email).
		SetRegisteredBy(payload.RegisteredBy).
		SetOnSite(payload.OnSite).
		Save(ctx)

	return b, err
}

func (r *repository) CreateBusinssLocation(ctx context.Context, payload createBusinssLocation) error {
	_, err := r.db.BusinessLocation.
		Create().
		SetName(payload.Name).
		SetAddress(payload.Street).
		SetCity(payload.City).
		SetState(payload.State).
		SetCountry(payload.Country).
		SetZipCode(payload.ZipCode).
		SetNillableLatitude(payload.Latitude).
		SetNillableLongitude(payload.Longitude).
		SetIsHeadOffice(payload.IsHeadOffice).
		SetBusiness(payload.Business).
		Save(ctx)

	return err
}

// AddDocument implements [IRepository].
func (r *repository) AddDocument(ctx context.Context, business *ent.Business, payload AddBusinessDocument) error {
	_, err := r.db.BusinessDocument.Create().
		SetTitle(payload.Title).
		SetDescription(payload.Description).
		SetURL(payload.URL).
		SetType(businessdocument.Type(payload.Type)).
		SetBusiness(business).
		AddKybDocumentIDs(payload.DocumentID).
		Save(ctx)

	return err
}

// AddBusinessManager implements [IRepository].
func (r *repository) AddBusinessManager(ctx context.Context, businessID, userID string, isOwner bool) error {
	_, err := r.db.Manager.
		Create().
		SetBusinessID(businessID).
		SetUserID(userID).
		SetIsOwner(isOwner).
		Save(ctx)

	return err
}

// GetManager implements [IRepository].
func (r *repository) GetManager(ctx context.Context, u *ent.User, b *ent.Business) (*ent.Manager, error) {
	m, err := r.db.Manager.Query().
		Where(manager.HasUserWith(user.IDEQ(u.ID)), manager.HasBusinessWith(business.IDEQ(b.ID))).
		WithRole().
		Only(ctx)

	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		return nil, err
	}

	return m, err
}

// GetBusinessDocuments implements [IRepository].
func (r *repository) GetBusinessDocument(ctx context.Context, b *ent.Business, documentID string) (*ent.BusinessDocument, error) {
	doc, err := r.db.BusinessDocument.Query().
		Where(businessdocument.
			HasBusinessWith(business.IDEQ(b.ID)), businessdocument.IDEQ(documentID), businessdocument.DeletedAtIsNil(),
		).WithKybDocument().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return doc, err
}

func (r *repository) GetKYBDocuments(ctx context.Context, businessID string) ([]*ent.BusinessDocument, error) {
	return r.db.BusinessDocument.
		Query().
		Where(businessdocument.TypeEQ(businessdocument.TypeKYB),
			businessdocument.HasBusinessWith(
				business.IDEQ(businessID),
			)).
		WithKybDocument().
		All(ctx)
}

func (r *repository) KYBDocumentByID(ctx context.Context, id string) (*ent.KYBDocument, error) {
	doc, err := r.db.KYBDocument.Query().
		Where(kybdocument.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return doc, err
}

func (r *repository) DeleteDocument(ctx context.Context, id string) error {
	return r.db.BusinessDocument.DeleteOneID(id).Exec(ctx)
}

func (r *repository) KYBDOcuments(ctx context.Context) ([]*ent.KYBDocument, error) {
	return r.db.KYBDocument.Query().Where(kybdocument.ActiveEQ(true)).All(ctx)
}

func (r *repository) UpdateDcument(ctx context.Context, bixID, documentID, url string) error {
	err := r.db.BusinessDocument.UpdateOneID(documentID).
		Where(
			businessdocument.HasBusinessWith(business.IDEQ(bixID)),
		).
		SetURL(url).
		Exec(ctx)

	return err
}

// Update implements [IRepository].
func (r *repository) Update(ctx context.Context, b *ent.Business, payload UpdateBusiness) error {
	query := r.db.Business.UpdateOne(b)

	if payload.Name != nil && b.VerificationStatus != business.VerificationStatusVERIFIED {
		query.SetName(util.AddressToString(payload.Name))
	}

	if payload.About != nil {
		query.SetAbout(util.AddressToString(payload.About))
	}

	if payload.CountryOfIncorporation != nil {
		query.SetCountryOfIncorporation(util.AddressToString(payload.CountryOfIncorporation))
	}

	if payload.DateOfIncorporation != nil {
		query.SetDateOfIncorporation(util.AddressToString(payload.DateOfIncorporation))
	}

	if payload.Number != nil {
		query.SetRegistrationNumber(util.AddressToString(payload.Number))
	}

	if payload.Industry != nil {
		query.SetCategory(util.AddressToString(payload.Industry))
	}

	if payload.Website != nil {
		query.SetWebsite(util.AddressToString(payload.Website))
	}

	_, err := query.Save(ctx)

	return err
}

func (r *repository) WithTransaction(tx *ent.Tx) IRepository {
	return &repository{
		db: tx.Client(),
	}
}
