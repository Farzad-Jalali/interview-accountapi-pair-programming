package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/form3tech/go-data/data"
	application_errors "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/go-openapi/strfmt"
	uuid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const unique_violation = "23505"

type genericRecord struct {
	ID             interface{}
	OrganisationID interface{}
	Version        int64
	CreatedOn      interface{}
	ModifiedOn     interface{}
	Record         interface{}
}

type Storage struct {
	db        *sqlx.DB
	tableName string
}

func (s *Storage) GetByID(ID uuid.UUID, result interface{}) error {
	sqlStmt, params, err := data.
		Select("*").
		From("\"" + s.tableName + "\"").
		Where(sq.Eq{"id": ID.String()}).ToSql()

	if err != nil {
		return err
	}

	err = s.db.Get(result, sqlStmt, params...)

	if err == sql.ErrNoRows {
		return application_errors.NewNotFoundError(fmt.Sprintf("record %v does not exist", ID))
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CheckExists(ID strfmt.UUID) (bool, error) {
	sqlStmt, params, err := data.
		Select("count(*)").
		From("\"" + s.tableName + "\"").
		Where(sq.Eq{"id": ID}).ToSql()

	if err != nil {
		return false, err
	}

	var result int64
	err = s.db.Get(&result, sqlStmt, params...)

	return result > 0, err
}

func (s *Storage) getGenericRecord(dataRecord interface{}) (*genericRecord, error) {
	dataRecordValue := reflect.ValueOf(dataRecord)
	for dataRecordValue.Kind() == reflect.Ptr {
		dataRecordValue = dataRecordValue.Elem()
	}

	recordField := dataRecordValue.FieldByName("Record")
	if !recordField.IsValid() {
		return nil, errors.New("no Record field")
	}
	record := recordField.Interface()

	IDField := dataRecordValue.FieldByName("ID")
	if !IDField.IsValid() {
		return nil, errors.New("no ID field")
	}
	ID := IDField.Interface()

	organisationIDField := dataRecordValue.FieldByName("OrganisationID")
	if !organisationIDField.IsValid() {
		return nil, errors.New("no OrganisationID field")
	}
	organisationID := organisationIDField.Interface()

	createdOnField := dataRecordValue.FieldByName("CreatedOn")
	if !createdOnField.IsValid() {
		return nil, errors.New("no CreatedOn field")
	}
	createdOn := createdOnField.Interface()

	modifiedOnField := dataRecordValue.FieldByName("ModifiedOn")
	if !modifiedOnField.IsValid() {
		return nil, errors.New("no ModifiedOn field")
	}
	modifiedOn := modifiedOnField.Interface()

	versionField := dataRecordValue.FieldByName("Version")
	if !versionField.IsValid() {
		return nil, errors.New("no Version field")
	}

	version := int64(0)
	if !versionField.IsNil() {
		for versionField.Kind() == reflect.Ptr {
			versionField = versionField.Elem()
		}
		version = versionField.Interface().(int64)
	}
	return &genericRecord{
		ID:             ID,
		OrganisationID: organisationID,
		Version:        version,
		CreatedOn:      createdOn,
		ModifiedOn:     modifiedOn,
		Record:         record,
	}, nil
}

func setVersion(dataRecord interface{}, version int64) error {
	dataRecordValue := reflect.ValueOf(dataRecord)
	for dataRecordValue.Kind() == reflect.Ptr {
		dataRecordValue = dataRecordValue.Elem()
	}

	versionField := dataRecordValue.FieldByName("Version")
	if !versionField.IsValid() || !versionField.CanSet() {
		return errors.New("no Version field")
	}
	versionField.Set(reflect.ValueOf(&version))
	return nil
}

func (s *Storage) Add(dataRecord interface{}) error {
	// TODO: Try sqlx here.
	recordDetails, err := s.getGenericRecord(dataRecord)
	if err != nil {
		return err
	}

	sqlStmt, params, err := data.Insert("\""+s.tableName+"\"").
		Columns("id", "organisationid", "version", "isdeleted", "islocked", "createdon", "modifiedon", "record").
		Values(recordDetails.ID, recordDetails.OrganisationID, 0, 0, 0, recordDetails.CreatedOn, recordDetails.ModifiedOn, recordDetails.Record).
		ToSql()

	if err != nil {
		return err
	}

	_, err = s.db.Exec(sqlStmt, params...)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == unique_violation {
			return application_errors.NewDuplicateError("Cannot insert duplicate record")
		}

		return fmt.Errorf("could not insert Report, error: %v", err)
	}
	_ = setVersion(dataRecord, 0)
	return nil
}

func (s *Storage) Update(dataRecord interface{}) error {
	recordDetails, err := s.getGenericRecord(dataRecord)
	if err != nil {
		return err
	}

	sqlStmt, params, err := data.Update("\""+s.tableName+"\"").
		Set("record", recordDetails.Record).
		Set("version", sq.Expr("version + 1")).
		Where(sq.Eq{"id": recordDetails.ID, "version": recordDetails.Version}).
		ToSql()

	if err != nil {
		return err
	}

	res, err := s.db.Exec(sqlStmt, params...)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok && pqError.Code == unique_violation {
			return application_errors.NewDuplicateError("Cannot insert duplicate record")
		}

		return fmt.Errorf("could not insert Report, error: %v", err)
	}

	rows, err := res.RowsAffected()
	if rows == 0 || err != nil {
		return application_errors.NewConflictError(fmt.Sprintf("unable to update expected version %d", recordDetails.Version))
	}

	if err := setVersion(dataRecord, recordDetails.Version+1); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteRecord(predicate interface{}) error {
	sqlStm, params, err := data.Delete(s.tableName).
		Where(predicate).
		ToSql()
	if err != nil {
		return err
	}
	_, err = s.db.Exec(sqlStm, params...)
	if err != nil {
		return fmt.Errorf("database error - failed to delete record: %s", err)
	}
	return nil
}