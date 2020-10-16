package security

import (
	"context"
	"fmt"

	"github.com/ahmetb/go-linq"

	"github.com/google/uuid"
)

type SecuredOrganisations []uuid.UUID

var UnlimitedOrganisations SecuredOrganisations = nil

// For resources with organisation
func CheckPermission(ctx *context.Context, recordType string, organisationId *uuid.UUID, requiredPermissions ...AuthoriseAction) error {
	if organisationId == nil || ctx == nil || IsApplicationContext(*ctx) {
		return nil
	}

	c := *ctx
	userID, ok := c.Value(contextKeyUserID).(string)
	if !ok {
		return NewAuthError("no user id found in context")
	}

	acls, ok := c.Value(contextKeyACLs).(AccessControl)
	if !ok {
		return NewAuthError("no acls found in context")
	}

	for _, requiredPermission := range requiredPermissions {
		if !acls.HasAccess(*organisationId, recordType, requiredPermission) {
			return NewAuthError(fmt.Sprintf("%s unauthorised, to execute this action you need actions: %s on record type: %s for organisation %s",
				permissionsToString(requiredPermissions), userID, recordType, organisationId))
		}
	}
	return nil
}

// For resources without organisation
func CheckPermissionForResourceWithoutOrganisation(ctx *context.Context, recordType string, requiredPermissions ...AuthoriseAction) error {
	if ctx == nil || IsApplicationContext(*ctx) {
		return nil
	}
	c := *ctx
	acls, ok := c.Value(contextKeyACLs).(AccessControl)
	if !ok {
		return NewAuthError("no acls found in context")
	}

	for _, requiredPermission := range requiredPermissions {
		if !acls.ContainsAcl(recordType, requiredPermission) {
			return NewAuthError(fmt.Sprintf("unauthorised, to execute this action you need actions: %s on record type: %s",
				permissionsToString(requiredPermissions), recordType))
		}
	}
	return nil
}

func GetOrganisationsWithPermission(ctx *context.Context, recordType string, requiredPermission AuthoriseAction) (SecuredOrganisations, error) {
	if ctx == nil {
		return nil, nil
	}
	c := *ctx
	if IsApplicationContext(c) {
		return UnlimitedOrganisations, nil
	}

	acls, ok := c.Value(contextKeyACLs).(AccessControl)
	if !ok {
		return nil, NewAuthError("no acls found in context")
	}
	return acls.AllowedOrganisations(recordType, requiredPermission), nil
}

func (s SecuredOrganisations) IsUnlimited() bool {
	return s == nil
}

func (s SecuredOrganisations) IntersectFilter(filtered []uuid.UUID) SecuredOrganisations {
	if len(filtered) > 0 {
		if s.IsUnlimited() {
			return filtered
		} else {
			var result []uuid.UUID
			linq.From(filtered).Intersect(linq.From(s)).ToSlice(&result)
			return result
		}
	}
	return s
}
