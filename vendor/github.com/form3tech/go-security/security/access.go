package security

import (
	"context"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
)

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

		acls, ok := c.Value("acls").(AccessControlList)

		if !ok {
			return nil, NewAuthError("no acls found in context")
		}

		var allowedOrganisationIds []uuid.UUID

		linq.From(acls).
			WhereT(func(e AccessControlListEntry) bool { return e.RecordType == recordType && e.Action == READ }).
			SelectT(func(e AccessControlListEntry) uuid.UUID { return e.OrganisationId }).
			ToSlice(&allowedOrganisationIds)

		return allowedOrganisationIds, nil
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
