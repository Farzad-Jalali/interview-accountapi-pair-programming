package security

import (
	"github.com/google/uuid"
)

type AccessControlListEntry struct {
	OrganisationId uuid.UUID
	Action         AuthoriseAction
	RecordType     string
}

type AccessControlList []AccessControlListEntry

func (a AccessControlList) ContainsAcl(actionName string, control AuthoriseAction) bool {
	for _, entry := range a {
		if entry.RecordType == actionName && entry.Action&control > 0 {
			return true
		}
	}
	return false
}

func (a AccessControlList) HasAccess(organisationID uuid.UUID, actionName string, control AuthoriseAction) bool {
	for _, entry := range a {
		if entry.OrganisationId == organisationID && entry.RecordType == actionName && entry.Action&control > 0 {
			return true
		}
	}
	return false
}

func (a AccessControlList) AllowedOrganisations(actionName string, control AuthoriseAction) []uuid.UUID {
	results := []uuid.UUID{}
	for _, entry := range a {
		if entry.RecordType == actionName && entry.Action&control > 0 {
			results = append(results, entry.OrganisationId)
		}
	}
	return results
}

var _ AccessControl = AccessControlList{}
