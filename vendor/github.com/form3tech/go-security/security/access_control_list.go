package security

import (
	"github.com/google/uuid"
)

type AccessControl interface {
	ContainsAcl(actionName string, control AuthoriseAction) bool
	HasAccess(organisationID uuid.UUID, actionName string, control AuthoriseAction) bool
	AllowedOrganisations(actionName string, control AuthoriseAction) []uuid.UUID
}

var NilAccessControl AccessControl = nilAccessControl{}

type nilAccessControl struct {
}

func (n nilAccessControl) ContainsAcl(_ string, _ AuthoriseAction) bool {
	return false
}

func (n nilAccessControl) HasAccess(_ uuid.UUID, _ string, _ AuthoriseAction) bool {
	return false
}

func (n nilAccessControl) AllowedOrganisations(_ string, _ AuthoriseAction) []uuid.UUID {
	return []uuid.UUID{}
}
