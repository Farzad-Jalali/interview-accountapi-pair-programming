package convert

import (
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/swagger-client/interview-accountapi/models"
	"github.com/go-openapi/strfmt"

)

func ToAccountDataRecord(item *models.AccountCreation) (*internalmodels.AccountRecord, error) {
	id, err := ToUUID(item.Data.ID)
	if err != nil {
		return nil, err
	}
	organisationId, err := ToUUID(item.Data.OrganisationID)
	if err != nil {
		return nil, err
	}
	return &internalmodels.AccountRecord{
		ID:             id,
		OrganisationID: organisationId,
		Record: internalmodels.Account{
			AccountClassification: item.Data.Attributes.AccountClassification,
			AccountMatchingOptOut: item.Data.Attributes.AccountMatchingOptOut,
			AccountNumber: item.Data.Attributes.AccountNumber,
			AlternativeBankAccountNames: item.Data.Attributes.AlternativeBankAccountNames,
			BankAccountName: item.Data.Attributes.BankAccountName,
			BankID: item.Data.Attributes.BankID,
			BankIDCode: item.Data.Attributes.BankIDCode,
			BaseCurrency: item.Data.Attributes.BaseCurrency,
			Bic: item.Data.Attributes.Bic,
			Country: item.Data.Attributes.Country,
			CustomerID: item.Data.Attributes.CustomerID,
			FirstName: item.Data.Attributes.FirstName,
			Iban: item.Data.Attributes.Iban,
			JointAccount: item.Data.Attributes.JointAccount,
			SecondaryIdentification: item.Data.Attributes.SecondaryIdentification,
			Title: item.Data.Attributes.Title,
		},
	}, nil
}

func FromAccountDataRecord(record *internalmodels.AccountRecord) *models.Account {
	return &models.Account{
		ID:             FromUUID(record.ID),
		OrganisationID: FromUUID(record.OrganisationID),
		Type:           models.ResourceTypeAccounts,
		Version:        record.Version,
		ModifiedOn:     strfmt.DateTime(record.ModifiedOn),
		CreatedOn:      strfmt.DateTime(record.CreatedOn),
		Attributes: &models.AccountAttributes{
			AccountClassification: record.Record.AccountClassification,
			AccountMatchingOptOut: record.Record.AccountMatchingOptOut,
			AccountNumber: record.Record.AccountNumber,
			AlternativeBankAccountNames: record.Record.AlternativeBankAccountNames,
			BankAccountName: record.Record.BankAccountName,
			BankID: record.Record.BankID,
			BankIDCode: record.Record.BankIDCode,
			BaseCurrency: record.Record.BaseCurrency,
			Bic: record.Record.Bic,
			Country: record.Record.Country,
			CustomerID: record.Record.CustomerID,
			FirstName: record.Record.FirstName,
			Iban: record.Record.Iban,
			JointAccount: record.Record.JointAccount,
			SecondaryIdentification: record.Record.SecondaryIdentification,
			Title: record.Record.Title,
		},
	}
}
