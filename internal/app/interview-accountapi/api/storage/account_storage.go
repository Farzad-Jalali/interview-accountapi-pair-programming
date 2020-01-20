package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"
	"github.com/form3tech/go-data/data"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const accountTableName = `"Account"`
const duplicationErrorCode = "23505"

type AccountStorage struct {
	Storage
}

func NewAccountStorage(db *sqlx.DB) *AccountStorage {
	return &AccountStorage{
		Storage{
			db:        db,
			tableName: accountTableName,
		},
	}
}

func (a *AccountStorage) Create(record *internalmodels.AccountRecord) error {
	sqlStmt, params, err := data.Insert(a.tableName).
		Columns("id", "organisation_id", "version", "is_deleted", "is_locked", "created_on", "modified_on", "record").
		Values(
			record.ID,
			record.OrganisationID,
			record.Version,
			record.IsDeleted,
			record.IsLocked,
			record.CreatedOn,
			record.ModifiedOn,
			record.Record,
		).
		ToSql()
	if err != nil {
		return err
	}
	_, err = a.db.Exec(sqlStmt, params...)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok && dbErr.Code == duplicationErrorCode {
			return errors.NewDuplicateError("Account cannot be created as it violates a duplicate constraint")
		}
	}
	return err
}

func (a *AccountStorage) Delete(record *internalmodels.AccountRecord) error {
	pred := squirrel.Eq{
		"id":      record.ID,
		"version": record.Version,
	}
	return a.DeleteRecord(pred)
}
