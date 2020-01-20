package api

import (
	"context"
	"fmt"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/commands"
	"net/http"
	"strconv"

	"github.com/form3tech/go-form3-web/web"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/executors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/queries"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/convert"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/swagger-client/interview-accountapi/models"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

func getLogger(ctx *context.Context, c *gin.Context) log.Logger {
	logger := log.WithContext(ctx)
	if correlationID, ok := c.Get("correlation-id"); ok {
		logger = logger.WithField("correlation-id", correlationID)
	}
	return logger
}

func HandleCreateAccount(ctx *context.Context, c *gin.Context) error {
	getLogger(ctx, c).Debugf("Handling create account for %+v", c.Params)

	newAccount := &models.AccountCreation{}
	if err := c.BindJSON(newAccount); err != nil {
		return errors.NewIllegalArgumentError(err.Error())
	}
	if err := newAccount.Validate(strfmt.NewFormats()); err != nil {
		return errors.NewIllegalArgumentError(err.Error())
	}

	dataRecord, err := convert.ToAccountDataRecord(newAccount)
	if err != nil {
		return errors.NewIllegalArgumentError(err.Error())
	}

	err = executors.InMemoryCommandExecutor.Execute(ctx, &dataRecord.OrganisationID, commands.CreateAccountCommand{
		DataRecord: dataRecord,
	})
	if err != nil {
		return err
	}

	result := &queries.GetAccountByIdResult{}

	err = executors.QueryExecutor.Execute(ctx, queries.GetAccountByIdCriteriaBuilder(dataRecord.ID), &result)
	if err != nil {
		return err
	}
	response := toAccountDetailsResponse(c, result.DataRecord)
	c.JSON(http.StatusCreated, response)
	return nil
}

func HandleGetAccountById(ctx *context.Context, c *gin.Context) error {
	getLogger(ctx, c).Debugf("Handling get account for %+v", c.Params)

	id := c.Param("id")
	if !strfmt.IsUUID(id) {
		return errors.NewIllegalArgumentError(fmt.Sprintf("id is not a valid uuid"))
	}
	accountID, err := convert.ToUUID(strfmt.UUID(id))
	if err != nil {
		return err
	}

	result := &queries.GetAccountByIdResult{}
	err = executors.QueryExecutor.Execute(ctx, queries.GetAccountByIdCriteriaBuilder(accountID), &result)
	if err != nil {
		return err
	}
	response := toAccountDetailsResponse(c, result.DataRecord)
	c.JSON(http.StatusOK, response)
	return nil
}

func HandleDeleteAccount(ctx *context.Context, c *gin.Context) error {
	getLogger(ctx, c).Debugf("Handling delete account for %+v", c.Params)

	accountId, err := convert.ToUUID(strfmt.UUID(c.Param("id")))
	if err != nil {
		return errors.NewIllegalArgumentError(fmt.Sprintf("id is not a valid uuid"))
	}
	version, err := strconv.ParseInt(c.Query("version"), 10, 0)
	if err != nil {
		return errors.NewIllegalArgumentError(fmt.Sprintf("invalid version number"))
	}

	result := &queries.GetAccountByIdResult{}
	err = executors.QueryExecutor.Execute(ctx, queries.GetAccountByIdCriteriaBuilder(accountId), &result)
	if err != nil {
		if _, ok := err.(*errors.NotFoundError); ok {
			c.Status(http.StatusNoContent)
			return nil
		}
		return err
	}
	if result.DataRecord == nil {
		c.Status(http.StatusNoContent)
		return nil
	}
	if *result.DataRecord.Version != version {
		return errors.NewNotFoundError("invalid version")
	}
	err = executors.InMemoryCommandExecutor.Execute(ctx, &result.DataRecord.OrganisationID, commands.DeleteAccountCommand{
		AccountId: accountId,
		Version:       version,
	})
	if err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}

func HandleListAccounts(ctx *context.Context, c *gin.Context) error {
	getLogger(ctx, c).Debugf("Handling list accounts for %+v", c.Params)

	organisationIds, err := convert.ToUUIDs(c.QueryArray("filter[organisation_id]"))
	if err != nil {
		return errors.NewIllegalArgumentError(err.Error())
	}
	criteria := queries.NewListAccountsCriteriaBuilder().
		WithPageCriteria(web.BuildPageCriteria(c)).
		WithFilterByOrganisationId(organisationIds).
		Build()

	result := &queries.ListAccountsResult{}
	if err := executors.QueryExecutor.Execute(ctx, criteria, &result); err != nil {
		return err
	}

	var accounts []*models.Account
	for _, record := range result.DataRecords {
		report := convert.FromAccountDataRecord(record)
		accounts = append(accounts, report)
	}

	response := &models.AccountDetailsListResponse{
		Data: accounts,
	}
	if result.PageResults.TotalRecords != 0 {
		links := web.BuildListLinks(c, result.PageResults)
		response.Links = convert.FromLinks(links)
	}
	c.JSON(http.StatusOK, response)
	return nil
}

func toAccountDetailsResponse(c *gin.Context, data *internalmodels.AccountRecord) (response *models.AccountDetailsResponse) {
	account := convert.FromAccountDataRecord(data)
	links := web.BuildItemLinks(c, data.ID.String())
	return &models.AccountDetailsResponse{
		Data:  account,
		Links: convert.FromLinks(links),
	}
}