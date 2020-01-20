package interview_accountapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/convert"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/swagger-client/interview-accountapi/client/account_api"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/swagger-client/interview-accountapi/models"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type getAccountStage struct {
	t                                     *testing.T
	accountNumber                         string
	bankID                                string
	organisationId                        uuid.UUID
	createAccountsRequestModel            *models.NewAccount
	client                                *account_api.Client
	postOrganisationAccountsCreatedResult *account_api.PostOrganisationAccountsCreated
	error                                 error
	getAccountResult                      *account_api.GetOrganisationAccountsIDOK
}

func GetAccountTest(t *testing.T) (*getAccountStage, *getAccountStage, *getAccountStage) {
	stage := &getAccountStage{
		t:              t,
		organisationId: uuid.New(),
		bankID:         "400300",
		accountNumber:  "41426819",
	}
	return stage, stage, stage
}

func (s *getAccountStage) and() *getAccountStage {
	return s
}

func (s *getAccountStage) an_authorized_service_user() *getAccountStage {
	s.client = NewAccountAPIClient(ServerPort)
	return s
}

func (s *getAccountStage) an_account_with_number_and_bank_id(accountNumber string, bankID string) *getAccountStage {
	s.accountNumber = accountNumber
	s.bankID = bankID

	s.createAccountsRequestModel = &models.NewAccount{
		ID:             strfmt.UUID(uuid.New().String()),
		OrganisationID: convert.FromUUID(s.organisationId),
		Type:           string(models.ResourceTypeAccounts),
		Attributes: &models.AccountAttributes{
			AccountNumber: accountNumber,
			BankID:        bankID,
			Country:       convert.StringToPtr("GB"),
		},
	}

	s.postOrganisationAccountsCreatedResult, s.error = s.client.PostOrganisationAccounts(&account_api.PostOrganisationAccountsParams{
		Context: context.Background(),
		CreationRequest: &models.AccountCreation{
			Data: s.createAccountsRequestModel,
		},
	})
	return s
}

func (s *getAccountStage) fetching_an_account_by_id() *getAccountStage {
	s.getAccountResult, s.error = s.client.GetOrganisationAccountsID(&account_api.GetOrganisationAccountsIDParams{
		Context: context.Background(),
		ID:      s.postOrganisationAccountsCreatedResult.Payload.Data.ID,
	})
	return s
}

func (s *getAccountStage) fetching_an_account_by_a_non_existing_id() *getAccountStage {
	s.getAccountResult, s.error = s.client.GetOrganisationAccountsID(&account_api.GetOrganisationAccountsIDParams{
		Context: context.Background(),
		ID:      convert.FromUUID(uuid.New()),
	})
	return s
}

func (s *getAccountStage) the_account_should_be_found() *getAccountStage {
	assert.NoError(s.t, s.error)
	account := s.getAccountResult.Payload.Data

	assert.NotNil(s.t, account)
	assert.Equal(s.t, models.ResourceTypeAccounts, account.Type)
	assert.Equal(s.t, convert.FromUUID(s.organisationId), account.OrganisationID)
	assert.Equal(s.t, s.accountNumber, account.Attributes.AccountNumber)
	assert.Equal(s.t, s.bankID, account.Attributes.BankID)
	return s
}

func (s *getAccountStage) the_status_code_is_404_not_found() *getAccountStage {
	assert.Error(s.t, s.error)
	assert.IsType(s.t, &account_api.GetOrganisationAccountsIDNotFound{}, s.error)
	return s
}
