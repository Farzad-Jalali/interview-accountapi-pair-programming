package queries

import (
	"github.com/form3tech/go-cqrs/cqrs"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/executors"
)

func Configure() {
	errors.Must(executors.QueryExecutor.RegisterQuery(
		GetAccountByIdQuery,
		cqrs.WithNoFilter()),
	)
	errors.Must(executors.QueryExecutor.RegisterQuery(
		ListAccountsQuery,
		cqrs.WithNoFilter()),
	)
}
