package commandhandlers

import (
	"context"
	"time"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/commands"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/storage"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/jmoiron/sqlx"
)

func CreateAccountCommandHandler(ctx *context.Context, db *sqlx.DB, c commands.CreateAccountCommand) error {
	log.
		WithContext(ctx).
		WithField("organisation_id", c.DataRecord.OrganisationID.String()).
		Debug("Creating account...")

	record := c.DataRecord
	var defaultVersion int64 = 0
	now := time.Now().UTC()
	record.Version = &defaultVersion
	record.CreatedOn = now
	record.ModifiedOn = now
	record.IsLocked = false
	record.IsDeleted = false

	return storage.NewAccountStorage(db).Create(record)
}
