package internalmodels

import "testing"

func TestAccount_Scan(t *testing.T) {
	json := `{
  "country": "GB",
  "base_currency": "GBP",
  "account_number": "41426819",
  "bank_id": "400300",
  "bank_id_code": "GBDSC",
  "bic": "NWBKGB22",
  "iban": "GB11NWBK40030041426819",
  "title": "Ms",
  "first_name": "Samantha",
  "bank_account_name": "Samantha Holder",
  "alternative_bank_account_names": [
    "Sam Holder"
  ],
  "account_classification": "Personal",
  "joint_account": false,
  "account_matching_opt_out": false,
  "secondary_identification": "A1B2C3D4"
}`
	r := &Account{}
	if err := r.Scan([]byte(json)); err != nil {
		t.Errorf("Account.Scan() error = %v", err)
	}

}
