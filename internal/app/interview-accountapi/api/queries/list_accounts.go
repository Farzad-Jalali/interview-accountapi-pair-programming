package queries

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/form3tech/go-data/data"
	"github.com/form3tech/go-form3-web/web"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ListAccountsCriteria struct {
	pageCriteria            web.PageCriteria
	filteredOrganisationIds []uuid.UUID
}

type ListAccountCriteriaBuilder struct {
	data ListAccountsCriteria
}

func NewListAccountsCriteriaBuilder() *ListAccountCriteriaBuilder {
	return &ListAccountCriteriaBuilder{}
}
func (b *ListAccountCriteriaBuilder) WithPageCriteria(pageCriteria web.PageCriteria) *ListAccountCriteriaBuilder {
	b.data.pageCriteria = pageCriteria
	return b
}
func (b *ListAccountCriteriaBuilder) WithFilterByOrganisationId(organisationIds []uuid.UUID) *ListAccountCriteriaBuilder {
	b.data.filteredOrganisationIds = organisationIds
	return b
}
func (b *ListAccountCriteriaBuilder) Build() ListAccountsCriteria {
	return b.data
}

func (c ListAccountsCriteria) buildQuery(builder data.PagedBuilder) data.PagedBuilder {
	whereClause := squirrel.And{}

	if len(c.filteredOrganisationIds) > 0 {
		whereClause = append(whereClause, squirrel.Eq{"organisation_id": c.filteredOrganisationIds})
	}
	if len(whereClause) > 0 {
		return builder.Where(whereClause)
	}
	return builder
}

type ListAccountsResult struct {
	PageResults web.PageResults
	DataRecords []*internalmodels.AccountRecord
}

func ListAccountsQuery(ctx *context.Context, db *sqlx.DB, criteria ListAccountsCriteria) (*ListAccountsResult, error) {
	result := ListAccountsResult{}

	query := data.
		Paged("*").
		From(`"Account"`)
	query = criteria.buildQuery(query)

	countSqlStmt, params, err := query.ToCount().ToSql()
	if err != nil {
		return nil, err
	}
	rowCount := 0
	err = db.Get(&rowCount, countSqlStmt, params...)
	if err != nil {
		return nil, err
	}
	result.PageResults.TotalRecords = rowCount
	result.PageResults.CurrentPage = criteria.pageCriteria.GetPageNumber(rowCount)
	result.PageResults.PageSize = criteria.pageCriteria.PageSize

	sqlStmt, params, err := query.ToSelect(result.PageResults.CurrentPage, result.PageResults.PageSize).
		OrderBy("pagination_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = db.Select(&result.DataRecords, sqlStmt, params...)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
