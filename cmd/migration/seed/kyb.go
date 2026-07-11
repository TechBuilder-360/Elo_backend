package seed

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Toflex/directory_v2/ent"
	"github.com/Toflex/directory_v2/ent/kybdocument"
	"github.com/Toflex/directory_v2/pkg/business"
)

func seedKYBDocuments(ctx context.Context, db *ent.Client) error {
	var document = []*ent.KYBDocumentCreate{}

	documentTypes := business.GetKYBDocumentTypes()
	for _, docType := range documentTypes {
		document = append(document, db.KYBDocument.Create().SetName(docType.Name).SetRequired(docType.Required).SetActive(docType.Active))
	}

	return db.KYBDocument.CreateBulk(document...).OnConflict(
		// Specify the columns to check for conflicts
		sql.ConflictColumns(kybdocument.FieldName),
		// Specify the action to take on conflict (do nothing)
		sql.DoNothing(),
	).Exec(ctx)
}
