package security

import (
	"context"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
)

type AuthoriseAction int

const (
	CREATE                    AuthoriseAction = 1 << iota
	deprecated_CREATE_APPROVE                 //2
	DELETE                                    //4
	deprecated_DELETE_APPROVE                 //8
	EDIT                                      //16
	deprecated_EDIT_APPROVE                   //32
	READ                                      //64
	deprecated_APPROVE                        //128
	deprecated_REJECT                         //256
	deprecated_REJECT_APPROVE                 //512
)

var authoriseActionsMap = map[AuthoriseAction]string{
	CREATE:                    "CREATE",
	deprecated_CREATE_APPROVE: "CREATE_APPROVE",
	DELETE:                    "DELETE",
	deprecated_DELETE_APPROVE: "DELETE_APPROVE",
	EDIT:                      "EDIT",
	deprecated_EDIT_APPROVE:   "EDIT_APPROVE",
	READ:                      "READ",
	deprecated_APPROVE:        "APPROVE",
	deprecated_REJECT:         "REJECT",
	deprecated_REJECT_APPROVE: "REJECT_APPROVE",
}

func RestrictWithPermissions(recordType string, requiredPermissions ...AuthoriseAction) func(*context.Context, *uuid.UUID) error {
	return func(ctx *context.Context, organisationId *uuid.UUID) error {
		return CheckPermission(ctx, recordType, organisationId, requiredPermissions...)
	}
}

func QueryForRecordPermissions(recordType string) func(ctx *context.Context) ([]uuid.UUID, error) {
	return func(ctx *context.Context) ([]uuid.UUID, error) {
		if ctx == nil {
			return []uuid.UUID{}, nil
		}

		c := *ctx
		acls, ok := c.Value(contextKeyACLs).(AccessControl)

		if !ok {
			return nil, NewAuthError("no acls found in context")
		}

		return acls.AllowedOrganisations(recordType, READ), nil
	}
}

func permissionsToString(requiredPermissions []AuthoriseAction) string {

	var results []string

	linq.From(requiredPermissions).
		SelectT(func(a AuthoriseAction) string { return authoriseActionsMap[a] }).
		ToSlice(&results)

	return strings.Join(results, ", ")
}

func AllowEveryone() func(*context.Context, *uuid.UUID) error {
	return func(ctx *context.Context, organisationId *uuid.UUID) error {
		return nil
	}
}

type AuthError struct {
	message string
}

func ProduceAuthError(message string) error {
	return NewAuthError(message)
}

func NewAuthError(message string) *AuthError {
	return &AuthError{
		message: message,
	}
}

func (a *AuthError) Error() string {
	return a.message
}
