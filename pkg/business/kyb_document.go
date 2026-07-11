package business

type document string

const (
	// CertificateOfIncorporation is the document type for certificate of incorporation
	CertificateOfIncorporation document = "Certificate of Incorporation"

	// ArticlesOfAssociation is the document type for articles of association
	ArticlesOfAssociation document = "Articles of Association"

	// StatusCertificate is the document type for status certificate
	StatusCertificate document = "Status Certificate"

	// ProofOfAddress is the document type for proof of address
	ProofOfAddress document = "Proof of Address"

	// OtherDocument is the document type for other documents
	OtherDocument document = "Other Document"
)

type KYBDocumentType struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Active   bool   `json:"active"`
}

var kybDocument = make([]KYBDocumentType, 0)

func init() {
	kybDocument = append(kybDocument, KYBDocumentType{Name: string(CertificateOfIncorporation), Required: true, Active: true})
	kybDocument = append(kybDocument, KYBDocumentType{Name: string(ArticlesOfAssociation), Required: true, Active: true})
	kybDocument = append(kybDocument, KYBDocumentType{Name: string(StatusCertificate), Required: true, Active: true})
	kybDocument = append(kybDocument, KYBDocumentType{Name: string(ProofOfAddress), Required: true, Active: true})
	kybDocument = append(kybDocument, KYBDocumentType{Name: string(OtherDocument), Required: false, Active: true})
}

func GetKYBDocumentTypes() []KYBDocumentType {
	return kybDocument
}

func IsKYBDocumentType(docType string) bool {
	for _, doc := range kybDocument {
		if doc.Name == docType {
			return true
		}
	}
	return false
}
