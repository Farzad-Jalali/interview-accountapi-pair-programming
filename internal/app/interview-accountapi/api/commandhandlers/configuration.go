package commandhandlers

import (
	"github.com/form3tech/go-security/security"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/errors"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/executors"
)

func Configure() {
	errors.Must(executors.InMemoryCommandExecutor.RegisterCommandHandler(
		CreateAccountCommandHandler,
		security.AllowEveryone(),
	))
	errors.Must(executors.InMemoryCommandExecutor.RegisterCommandHandler(
		DeleteAccountCommandHandler,
		security.AllowEveryone(),
	))
}
