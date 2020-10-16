package security

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type aclFullDecoder struct {
}

var _ Decoder = &aclFullDecoder{}

func NewAclFull() Decoder {
	return newDecompress(&aclFullDecoder{})
}

func NewAclFullWithBase64Decoder() Decoder {
	return newBase64Decoder(NewAclFull())
}

func splitSubN(s string, n int) []string {
	sub := ""
	var subs []string

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

var authoriseActions = []AuthoriseAction{
	CREATE,
	deprecated_CREATE_APPROVE,
	DELETE,
	deprecated_DELETE_APPROVE,
	EDIT,
	deprecated_EDIT_APPROVE,
	READ,
	deprecated_APPROVE,
	deprecated_REJECT,
	deprecated_REJECT_APPROVE,
}

func computeActions(actionSum int) []AuthoriseAction {

	var results []AuthoriseAction

	runningSum := actionSum

	for i := len(authoriseActions) - 1; i >= 0; i-- {

		actionValue := int(authoriseActions[i])
		if runningSum >= actionValue {
			runningSum -= actionValue
			results = append(results, authoriseActions[i])
		}
	}

	return results

}

func (a *aclFullDecoder) Decode(data []byte) (AccessControl, error) {
	results := AccessControlList{}

	payload := string(data)
	organisationIdRegex, _ := regexp.Compile(`{([a-f0-9]+)}`)
	matches := organisationIdRegex.FindStringSubmatch(payload)

	if len(matches) < 2 {
		return NilAccessControl, nil
	}

	organisationIdMap := make(map[int]uuid.UUID)

	for i, u := range splitSubN(matches[1], 32) {
		organisationIdMap[i], _ = uuid.Parse(u)
	}

	organisationStrLen := len(fmt.Sprintf("%d", len(organisationIdMap)))

	remainingAcl := strings.Replace(payload, matches[0], "", 1)

	actionRegex, _ := regexp.Compile(`(\w+)\[(\d+:\d+(\|\d+:\d+)*)]+`)

	for len(remainingAcl) > 0 {

		actionMatches := actionRegex.FindStringSubmatch(remainingAcl)

		for _, p := range strings.Split(actionMatches[2], "|") {
			split := strings.Split(p, ":")
			actionSum, _ := strconv.Atoi(split[0])
			actions := computeActions(actionSum)

			for i := 0; i < len(split[1]); i = i + organisationStrLen {

				index, _ := strconv.Atoi(split[1][i : i+organisationStrLen])

				for _, action := range actions {

					results = append(results, AccessControlListEntry{
						RecordType:     actionMatches[1],
						OrganisationId: organisationIdMap[index],
						Action:         action,
					})
				}
			}

		}

		remainingAcl = strings.Replace(remainingAcl, actionMatches[0], "", 1)

	}

	return results, nil
}
