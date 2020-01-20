package accountapi

type AccountClassification string

const (
	AccountClassificationPersonal AccountClassification = "Personal"
	AccountClassificationBusiness AccountClassification = "Business"
)

type AccountListData struct {
	Data []Account `json:"data"`
}

type AccountData struct {
	Data Account `json:"data"`
}

type Account struct {
	Type           string            `json:"type"`
	ID             string            `json:"id"`
	OrganisationID string            `json:"organisation_id"`
	Version        int               `json:"version"`
	Attributes     AccountAttributes `json:"attributes"`
}

type AccountAttributes struct {
	Country                     string                `json:"country"`
	BaseCurrency                string                `json:"base_currency"`
	AccountNumber               string                `json:"account_number"`
	BankID                      string                `json:"bank_id"`
	BankIDCode                  string                `json:"bank_id_code"`
	Bic                         string                `json:"bic"`
	IBAN                        string                `json:"iban"`
	Title                       string                `json:"title"`
	FirstName                   string                `json:"first_name"`
	BankAccountName             string                `json:"bank_account_name"`
	AlternativeBankAccountNames []string              `json:"alternative_bank_account_names"`
	AccountClassification       AccountClassification `json:"account_classification"`
	JointAccount                bool                  `json:"joint_account"`
	AccountMatchingOptOut       bool                  `json:"account_matching_opt_out"`
	SecondaryIdentification     string                `json:"secondary_identification"`
}

type ErrorResponse struct {
	Message string `json:"error_message"`
}