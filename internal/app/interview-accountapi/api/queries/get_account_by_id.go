package queries

import (
	"fmt"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"

	"github.com/Masterminds/squirrel"
	"github.com/form3tech/go-data/data"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GetAccountByIdResult struct {
	OrganisationId uuid.UUID
	DataRecord     *internalmodels.AccountRecord
}

type GetAccountByIdCriteria struct {
	AccountId uuid.UUID
}

func GetAccountByIdCriteriaBuilder(id uuid.UUID) GetAccountByIdCriteria {
	return GetAccountByIdCriteria{
		AccountId: id,
	}
}

func GetAccountByIdQuery(db *sqlx.DB, q GetAccountByIdCriteria) (*GetAccountByIdResult, error) {
	dataRecord := &internalmodels.AccountRecord{}

	sqlStmt, params, err := data.
		Select("*").
		From(`"Account"`).
		Where(squirrel.Eq{"id": q.AccountId}).
		ToSql()
	if err != nil {
		return nil, err
	}
	if err := db.Get(dataRecord, sqlStmt, params...); err != nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("record %v does not exist", q.AccountId))
	}
	return &GetAccountByIdResult{
		DataRecord:     dataRecord,
		OrganisationId: dataRecord.OrganisationID,
	}, nil
}
