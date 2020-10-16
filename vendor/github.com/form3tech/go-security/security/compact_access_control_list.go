package security

import (
	"github.com/google/uuid"
)

type RecordControl struct {
	_       struct{}
	Name    string
	Control AuthoriseAction
}

type CompactAccessControl struct {
	_               struct{}
	OrganisationIDs []uuid.UUID
	RecordControl   []RecordControl
	// should be accessed using [accessControlIdx][]{organisationIdx1, organisationIdx2, ... }
	AccessControlMap map[int][]int
}

func (a CompactAccessControl) findOrganisationPosition(organisationID uuid.UUID) int {
	for idx, organisation := range a.OrganisationIDs {
		if organisation == organisationID {
			return idx
		}
	}
	return -1
}

func (a CompactAccessControl) findAccessControlPosition(accessControl RecordControl) []int {
	var result []int
	for idx, access := range a.RecordControl {
		hasAccess := access.Control&accessControl.Control > 0
		if hasAccess && access.Name == accessControl.Name {
			result = append(result, idx)
		}
	}
	return result
}

func (a CompactAccessControl) ContainsAcl(actionName string, control AuthoriseAction) bool {
	for _, actionControl := range a.RecordControl {
		hasAccess := actionControl.Control&control > 0
		if hasAccess && actionControl.Name == actionName {
			return true
		}
	}
	return false
}

func (a CompactAccessControl) HasAccess(organisationID uuid.UUID, actionName string, control AuthoriseAction) bool {
	orgIdx := a.findOrganisationPosition(organisationID)
	if orgIdx == -1 {
		return false
	}
	actionIdxList := a.findAccessControlPosition(RecordControl{Name: actionName, Control: control})
	if len(actionIdxList) == 0 {
		return false
	}
	for _, actionIdx := range actionIdxList {
		for _, idx := range a.AccessControlMap[actionIdx] {
			if idx == orgIdx {
				return true
			}
		}
	}
	return false
}

func (a CompactAccessControl) AllowedOrganisations(actionName string, control AuthoriseAction) []uuid.UUID {
	actionIdxList := a.findAccessControlPosition(RecordControl{Name: actionName, Control: control})
	if len(actionIdxList) == 0 {
		return []uuid.UUID{}
	}

	result := []uuid.UUID{}
	for _, actionIdx := range actionIdxList {
		for _, organisationIdx := range a.AccessControlMap[actionIdx] {
			result = append(result, a.OrganisationIDs[organisationIdx])
		}
	}
	return result
}

var _ AccessControl = &CompactAccessControl{}
