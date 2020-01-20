package commandhandlers

import (
	"context"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/commands"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/storage"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/jmoiron/sqlx"
)

func DeleteAccountCommandHandler(ctx *context.Context, db *sqlx.DB, c commands.DeleteAccountCommand) error {
	log.
		WithContext(ctx).
		WithField("account_id", c.AccountId.String()).
		Debug("Deleting account...")

	dataRecord := &internalmodels.AccountRecord{
		Version: &c.Version,
		ID:      c.AccountId,
	}
	return storage.NewAccountStorage(db).Delete(dataRecord)
}
