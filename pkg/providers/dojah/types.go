package dojah

type documentType string
type WebhookPayload struct {
	Metadata           Metadata           `json:"metadata"`
	Data               WebhookPayloadData `json:"data"`
	IDType             string             `json:"id_type"`
	Value              string             `json:"value"`
	IDURL              string             `json:"id_url"`
	BackURL            string             `json:"back_url"`
	Message            string             `json:"message"`
	ReferenceID        string             `json:"reference_id"`
	WidgetID           string             `json:"widget_id"`
	VerificationMode   string             `json:"verification_mode"`
	VerificationType   string             `json:"verification_type"`
	VerificationValue  string             `json:"verification_value"`
	VerificationURL    string             `json:"verification_url"`
	SelfieURL          string             `json:"selfie_url"`
	Status             bool               `json:"status"`
	Aml                Aml                `json:"aml"`
	VerificationStatus string             `json:"verification_status"`
}

type Aml struct {
	Status bool `json:"status"`
}

type WebhookPayloadData struct {
	Index          Index          `json:"index"`
	Email          Email          `json:"email"`
	Countries      Countries      `json:"countries"`
	GovernmentData GovernmentData `json:"government_data"`
	ID             ID             `json:"id"`
}

type Countries struct {
	Data    CountriesData `json:"data"`
	Message string        `json:"message"`
	Status  bool          `json:"status"`
}

type CountriesData struct {
	Country string `json:"country"`
}

type Email struct {
	Data    EmailData `json:"data"`
	Status  bool      `json:"status"`
	Message string    `json:"message"`
}

type EmailData struct {
	Email string `json:"email"`
}

type GovernmentData struct {
	Data    GovernmentDataData `json:"data"`
	Message string             `json:"message"`
	Status  bool               `json:"status"`
}

type GovernmentDataData struct {
	Nin Nin `json:"nin"`
	Bvn Bvn `json:"bvn"`
}

type Nin struct {
	Entity NINEntity `json:"entity"`
}

type NINEntity struct {
	Customer              string      `json:"customer"`
	AppID                 interface{} `json:"app_id"`
	Nin                   string      `json:"nin"`
	FirstName             string      `json:"first_name"`
	LastName              string      `json:"last_name"`
	MiddleName            string      `json:"middle_name"`
	Gender                string      `json:"gender"`
	DateOfBirth           string      `json:"date_of_birth"`
	PhoneNumber           interface{} `json:"phone_number"`
	ImageURL              string      `json:"image_url"`
	Email                 interface{} `json:"email"`
	EmploymentStatus      interface{} `json:"employment_status"`
	MaritalStatus         string      `json:"marital_status"`
	BirthCountry          string      `json:"birth_country"`
	BirthLGA              interface{} `json:"birth_lga"`
	BirthState            interface{} `json:"birth_state"`
	EducationalLevel      interface{} `json:"educational_level"`
	MaidenName            string      `json:"maiden_name"`
	NspokenLang           interface{} `json:"nspoken_lang"`
	Profession            interface{} `json:"profession"`
	Religion              interface{} `json:"religion"`
	ResidenceAddressLine1 interface{} `json:"residence_address_line_1"`
	ResidenceAddressLine2 interface{} `json:"residence_address_line_2"`
	ResidenceStatus       interface{} `json:"residence_status"`
	ResidenceTown         interface{} `json:"residence_town"`
	ResidenceLGA          interface{} `json:"residence_lga"`
	ResidenceState        string      `json:"residence_state"`
	OspokenLang           interface{} `json:"ospoken_lang"`
	OriginLGA             string      `json:"origin_lga"`
	OriginPlace           string      `json:"origin_place"`
	OriginState           interface{} `json:"origin_state"`
	Height                interface{} `json:"height"`
	PFirstName            interface{} `json:"p_first_name"`
	PMiddleName           interface{} `json:"p_middle_name"`
	PLastName             interface{} `json:"p_last_name"`
	Sc                    bool        `json:"sc"`
	CreatedAt             string      `json:"createdAt"`
	UpdatedAt             string      `json:"updatedAt"`
	Firstname             string      `json:"firstname"`
	Surname               string      `json:"surname"`
	Birthdate             string      `json:"birthdate"`
	Telephoneno           interface{} `json:"telephoneno"`
	Middlename            string      `json:"middlename"`
}

type ID struct {
	Data    IDData `json:"data"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type IDData struct {
	IDURL   string      `json:"id_url"`
	BackURL string      `json:"back_url"`
	IDData  IDDataClass `json:"id_data"`
}

type IDDataClass struct {
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	MiddleName     string       `json:"middle_name"`
	Nationality    string       `json:"nationality"`
	MrzStatus      string       `json:"mrz_status"`
	ExpiryDate     string       `json:"expiry_date"`
	DocumentType   documentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`
	DateOfBirth    string       `json:"date_of_birth"`
	DateIssued     string       `json:"date_issued"`
	Extras         string       `json:"extras"`
}

type Index struct {
	Data    IndexData `json:"data"`
	Message string    `json:"message"`
	Status  bool      `json:"status"`
}

type IndexData struct {
}

type Metadata struct {
	Ipinfo     Ipinfo `json:"ipinfo"`
	DeviceInfo string `json:"device_info"`
}

type Ipinfo struct {
	Status     string  `json:"status"`
	Country    string  `json:"country"`
	City       string  `json:"city"`
	District   string  `json:"district"`
	Zip        string  `json:"zip"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Timezone   string  `json:"timezone"`
	ISP        string  `json:"isp"`
	Org        string  `json:"org"`
	As         int64   `json:"as"`
	Mobile     bool    `json:"mobile"`
	Proxy      bool    `json:"proxy"`
	Hosting    bool    `json:"hosting"`
	Query      string  `json:"query"`
	RegionName string  `json:"region_name"`
}
type Bvn struct {
	Entity BvnEntity `json:"entity"`
}

type BvnEntity struct {
	Customer           string      `json:"customer"`
	AppID              interface{} `json:"app_id"`
	Bvn                string      `json:"bvn"`
	FirstName          string      `json:"first_name"`
	LastName           string      `json:"last_name"`
	MiddleName         string      `json:"middle_name"`
	Gender             string      `json:"gender"`
	DateOfBirth        string      `json:"date_of_birth"`
	PhoneNumber1       string      `json:"phone_number1"`
	PhoneNumber2       string      `json:"phone_number2"`
	ImageURL           interface{} `json:"image_url"`
	Email              string      `json:"email"`
	EnrollmentBank     string      `json:"enrollment_bank"`
	EnrollmentBranch   string      `json:"enrollment_branch"`
	LevelOfAccount     string      `json:"level_of_account"`
	LGAOfOrigin        string      `json:"lga_of_origin"`
	LGAOfResidence     string      `json:"lga_of_residence"`
	MaritalStatus      string      `json:"marital_status"`
	NameOnCard         string      `json:"name_on_card"`
	Nationality        string      `json:"nationality"`
	Nin                string      `json:"nin"`
	RegistrationDate   string      `json:"registration_date"`
	ResidentialAddress string      `json:"residential_address"`
	StateOfOrigin      string      `json:"state_of_origin"`
	StateOfResidence   string      `json:"state_of_residence"`
	Title              string      `json:"title"`
	Type               string      `json:"type"`
	Xc                 interface{} `json:"xc"`
	Sc                 bool        `json:"sc"`
	WatchListed        string      `json:"watch_listed"`
	CreatedAt          string      `json:"createdAt"`
	UpdatedAt          string      `json:"updatedAt"`
}
