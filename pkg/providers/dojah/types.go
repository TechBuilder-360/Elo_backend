package dojah

type documentType string

type WebhookPayload struct {
	Aml                Aml                `json:"aml"`
	Data               WebhookPayloadData `json:"data"`
	Value              string             `json:"value"`
	IDURL              string             `json:"id_url"`
	Status             bool               `json:"status"`
	IDType             string             `json:"id_type"`
	Message            string             `json:"message"`
	BackURL            string             `json:"back_url"`
	Metadata           Metadata           `json:"metadata"`
	SelfieURL          string             `json:"selfie_url"`
	ReferenceID        string             `json:"reference_id"`
	VerificationURL    string             `json:"verification_url"`
	VerificationMode   string             `json:"verification_mode"`
	VerificationType   string             `json:"verification_type"`
	VerificationValue  string             `json:"verification_value"`
	VerificationStatus string             `json:"verification_status"`
}

type Aml struct {
	Status bool `json:"status"`
}

type WebhookPayloadData struct {
	ID             ID             `json:"id"`
	Email          Email          `json:"email"`
	Index          Index          `json:"index"`
	Selfie         Selfie         `json:"selfie"`
	Countries      Countries      `json:"countries"`
	UserData       UserData       `json:"user_data"`
	GovernmentData GovernmentData `json:"government_data"`
	BusinessID     BusinessID     `json:"business_id"`
}

type BusinessID struct {
	ImageURL         string `json:"image_url"`
	BusinessName     string `json:"business_name"`
	BusinessType     string `json:"business_type"`
	BusinessNumber   string `json:"business_number"`
	BusinessAddress  string `json:"business_address"`
	RegistrationDate string `json:"registration_date"`
}

type Countries struct {
	Data    CountriesData `json:"data"`
	Status  bool          `json:"status"`
	Message string        `json:"message"`
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

type ID struct {
	Data    IDData `json:"data"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type IDData struct {
	IDURL   string      `json:"id_url"`
	IDData  IDDataClass `json:"id_data"`
	BackURL string      `json:"back_url"`
}

type IDDataClass struct {
	Extras         string       `json:"extras"`
	LastName       string       `json:"last_name"`
	FirstName      string       `json:"first_name"`
	MrzStatus      string       `json:"mrz_status"`
	DateIssued     string       `json:"date_issued"`
	ExpiryDate     string       `json:"expiry_date"`
	MiddleName     string       `json:"middle_name"`
	Nationality    string       `json:"nationality"`
	DateOfBirth    string       `json:"date_of_birth"`
	DocumentType   documentType `json:"document_type"`
	DocumentNumber string       `json:"document_number"`
}

type Index struct {
	Data    IndexData `json:"data"`
	Status  bool      `json:"status"`
	Message string    `json:"message"`
}

type IndexData struct {
}

type Selfie struct {
	Data    SelfieData `json:"data"`
	Status  bool       `json:"status"`
	Message string     `json:"message"`
}

type SelfieData struct {
	SelfieURL string `json:"selfie_url"`
}

type UserData struct {
	Data    UserDataData `json:"data"`
	Status  bool         `json:"status"`
	Message string       `json:"message"`
}

type UserDataData struct {
	Dob       string `json:"dob"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}

type Metadata struct {
	Ipinfo     Ipinfo `json:"ipinfo"`
	DeviceInfo string `json:"device_info"`
}

type Ipinfo struct {
	As         string  `json:"as"`
	ISP        string  `json:"isp"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Org        string  `json:"org"`
	Zip        string  `json:"zip"`
	City       string  `json:"city"`
	Proxy      bool    `json:"proxy"`
	Query      string  `json:"query"`
	Mobile     bool    `json:"mobile"`
	Status     string  `json:"status"`
	Country    string  `json:"country"`
	Hosting    bool    `json:"hosting"`
	District   string  `json:"district"`
	Timezone   string  `json:"timezone"`
	RegionName string  `json:"region_name"`
}

type GovernmentData struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type Data struct {
	Bvn Bvn `json:"bvn"`
}

type Bvn struct {
	Entity Entity `json:"entity"`
}

type Entity struct {
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
