package commands

import "github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/internalmodels"

type CreateAccountCommand struct {
	DataRecord *internalmodels.AccountRecord
}
