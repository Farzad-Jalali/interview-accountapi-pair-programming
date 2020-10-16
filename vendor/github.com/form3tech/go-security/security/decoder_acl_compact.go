package security

import (
	"bytes"
	"math"

	"github.com/google/uuid"
)

const (
	aclOrganisationBlockStartByte      = '{'
	aclOrganisationBlockEndByte        = '}'
	aclListBlockEndByte                = ']'
	aclListBlockBeginByte              = '['
	aclControlAndOrganisationSeparator = ':'
	aclControlListSeparator            = '|'
	uuidLengthInBytes                  = 32
)

type aclCompact struct{}

var _ Decoder = &aclCompact{}

func NewAclCompact() Decoder {
	return newDecompress(&aclCompact{})
}

func NewAclCompactWithBase64Decoder() Decoder {
	return newBase64Decoder(NewAclCompact())
}

func newAccessControlList(nOrganisations int) *CompactAccessControl {
	return &CompactAccessControl{
		OrganisationIDs:  make([]uuid.UUID, 0, nOrganisations),
		RecordControl:    make([]RecordControl, 0),
		AccessControlMap: make(map[int][]int),
	}
}

func addOrganisationID(a *CompactAccessControl, organisationId uuid.UUID) {
	a.OrganisationIDs = append(a.OrganisationIDs, organisationId)
}

func addAccessControlForOrganisations(a *CompactAccessControl, actionName string, actionControl AuthoriseAction, orgIdxs []int) {
	a.RecordControl = append(a.RecordControl, RecordControl{
		Name:    actionName,
		Control: actionControl,
	})
	accessControlIdx := len(a.RecordControl) - 1
	a.AccessControlMap[accessControlIdx] = orgIdxs
}

func (a *aclCompact) Decode(payload []byte) (AccessControl, error) {
	orgBlockStartIdx := bytes.IndexByte(payload, aclOrganisationBlockStartByte)
	orgBlockEndIdx := bytes.IndexByte(payload, aclOrganisationBlockEndByte)
	nOrganisations := (orgBlockEndIdx - orgBlockStartIdx - 1) >> 5
	if orgBlockStartIdx == -1 || orgBlockEndIdx == -1 || nOrganisations == 0 {
		return NilAccessControl, nil
	}

	nOrganisationByteSize := countDigitsFor(nOrganisations)

	result := newAccessControlList(nOrganisations)
	for offset := orgBlockStartIdx + 1; offset < orgBlockEndIdx; offset += uuidLengthInBytes {
		organisationId, err := uuid.ParseBytes(payload[offset : offset+uuidLengthInBytes])
		if err != nil {
			return NilAccessControl, nil
		}
		addOrganisationID(result, organisationId)
	}

	payload = payload[orgBlockEndIdx+1:]
	for len(payload) > 0 {
		aclListBlockStartIdx := bytes.IndexByte(payload, aclListBlockBeginByte)
		if aclListBlockStartIdx == -1 {
			break
		}
		aclListBlockEndIdx := bytes.IndexByte(payload, aclListBlockEndByte)
		if aclListBlockEndIdx == -1 {
			break
		}

		actionName := string(payload[:aclListBlockStartIdx])
		aclControlList := payload[aclListBlockStartIdx+1 : aclListBlockEndIdx]
		payload = payload[aclListBlockEndIdx+1:]

		for len(aclControlList) > 0 {
			orgSeparatorIdx := bytes.IndexByte(aclControlList, aclControlAndOrganisationSeparator)
			listSeparatorIdx := bytes.IndexByte(aclControlList, aclControlListSeparator)
			if listSeparatorIdx == -1 {
				listSeparatorIdx = len(aclControlList)
			}

			actionControl := mustParseAuthoriseAction(aclControlList[:orgSeparatorIdx])
			nOrganisationIdxList := (listSeparatorIdx - orgSeparatorIdx - 1) / nOrganisationByteSize
			organisationIdxs := make([]int, 0, nOrganisationIdxList)

			for offset := orgSeparatorIdx + 1; offset < listSeparatorIdx; offset += nOrganisationByteSize {
				organisationIdx := mustParseToInt(aclControlList[offset : offset+nOrganisationByteSize])
				organisationIdxs = append(organisationIdxs, organisationIdx)
			}
			addAccessControlForOrganisations(result, actionName, actionControl, organisationIdxs)
			if len(aclControlList) <= listSeparatorIdx {
				break
			}
			aclControlList = aclControlList[listSeparatorIdx+1:]
		}
	}
	return result, nil
}

func mustParseAuthoriseAction(data []byte) AuthoriseAction {
	return AuthoriseAction(mustParseToInt(data))
}

func mustParseToInt(data []byte) int {
	val := 0
	for _, b := range data {
		val = val*10 + int(b-'0')
	}
	return val
}

func countDigitsFor(val int) int {
	return (int)(math.Log10(math.Abs(float64(val))) + 1)
}
