package security

import (
	"context"
	"fmt"

	linq "github.com/ahmetb/go-linq"
	"github.com/google/uuid"
)

type SecuredOrganisations []uuid.UUID

var UnlimitedOrganisations SecuredOrganisations = nil

// For resources with organisation
func CheckPermission(ctx *context.Context, recordType string, organisationId *uuid.UUID, requiredPermissions ...AuthoriseAction) error {

	if IsApplicationContext(*ctx) {
		return nil
	}

	if ctx == nil {
		return nil
	}

	if organisationId == nil {
		return nil
	}

	c := *ctx

	acls, ok := c.Value("acls").(AccessControlList)

	if !ok {
		return NewAuthError("no acls found in context")
	}

	for _, requiredPermission := range requiredPermissions {
		p := linq.From(acls).
			WhereT(func(e AccessControlListEntry) bool {
				return e.RecordType == recordType && e.Action == requiredPermission && e.OrganisationId == *organisationId
			}).Count()

		if p != 1 {
			return NewAuthError(fmt.Sprintf("unauthorised, to execute this action you need actions: %s on record type: %s for organisation %s",
				permissionsToString(requiredPermissions), recordType, organisationId))
		}
	}

	return nil
}

// For resources without organisation
func CheckPermissionForResourceWithoutOrganisation(ctx *context.Context, recordType string, requiredPermissions ...AuthoriseAction) error {

	if IsApplicationContext(*ctx) {
		return nil
	}

	if ctx == nil {
		return nil
	}

	c := *ctx

	acls, ok := c.Value("acls").(AccessControlList)

	if !ok {
		return NewAuthError("no acls found in context")
	}

	for _, requiredPermission := range requiredPermissions {
		p := linq.From(acls).
			WhereT(func(e AccessControlListEntry) bool {
				return e.RecordType == recordType && e.Action == requiredPermission
			}).Count()

		if p != 1 {
			return NewAuthError(fmt.Sprintf("unauthorised, to execute this action you need actions: %s on record type: %s",
				permissionsToString(requiredPermissions), recordType))
		}
	}

	return nil
}

func GetOrganisationsWithPermission(ctx *context.Context, recordType string, requiredPermission AuthoriseAction) (SecuredOrganisations, error) {

	res := []uuid.UUID{}

	if IsApplicationContext(*ctx) {
		return UnlimitedOrganisations, nil
	}

	if ctx == nil {
		return nil, nil
	}

	c := *ctx

	acls, ok := c.Value("acls").(AccessControlList)

	if !ok {
		return nil, NewAuthError("no acls found in context")
	}

	linq.From(acls).
		WhereT(func(e AccessControlListEntry) bool {
			return e.RecordType == recordType && e.Action == requiredPermission
		}).SelectT(func(elem AccessControlListEntry) uuid.UUID {
		return elem.OrganisationId
	}).ToSlice(&res)

	return res, nil
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
