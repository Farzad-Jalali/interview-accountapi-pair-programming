package security

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	linq "github.com/ahmetb/go-linq"
	"github.com/google/uuid"
)

type AuthoriseAction int

const (
	CREATE         AuthoriseAction = 1 << iota
	CREATE_APPROVE                 //2
	DELETE                         //4
	DELETE_APPROVE                 //8
	EDIT                           //16
	EDIT_APPROVE                   //32
	READ                           //64
	APPROVE                        //128
	REJECT                         //256
	REJECT_APPROVE                 //512
)

var authoriseActions = []AuthoriseAction{
	CREATE,
	CREATE_APPROVE,
	DELETE,
	DELETE_APPROVE,
	EDIT,
	EDIT_APPROVE,
	READ,
	APPROVE,
	REJECT,
	REJECT_APPROVE,
}

var authoriseActionsMap = map[AuthoriseAction]string{
	CREATE:         "CREATE",
	CREATE_APPROVE: "CREATE_APPROVE",
	DELETE:         "DELETE",
	DELETE_APPROVE: "DELETE_APPROVE",
	EDIT:           "EDIT",
	EDIT_APPROVE:   "EDIT_APPROVE",
	READ:           "READ",
	APPROVE:        "APPROVE",
	REJECT:         "REJECT",
	REJECT_APPROVE: "REJECT_APPROVE",
}

type AccessControlListEntry struct {
	OrganisationId uuid.UUID
	Action         AuthoriseAction
	RecordType     string
}

type AccessControlList []AccessControlListEntry

type encodedPermission struct {
	recordType        string
	actions           AuthoriseAction
	organisationIndex int
}

func uuidToSigBits(u uuid.UUID) (int64, int64) {
	msb := int64(0)
	lsb := int64(0)
	data := u[:]
	for i := 0; i < 8; i++ {
		msb = (msb << 8) | (int64(data[i]) & 0xff)
	}
	for i := 8; i < 16; i++ {
		lsb = (lsb << 8) | (int64(data[i]) & 0xff)
	}

	return msb, lsb
}

func compareUUIDs(a, b uuid.UUID) bool {
	msb1, lsb1 := uuidToSigBits(a)
	msb2, lsb2 := uuidToSigBits(b)

	if msb1 == msb2 {
		return lsb1 < lsb2
	}

	return msb1 < msb2
}

func Encode(accessControlList AccessControlList) ([]byte, error) {

	var organisationIds []string

	linq.From(accessControlList).
		SelectT(func(a AccessControlListEntry) uuid.UUID { return a.OrganisationId }).
		Distinct().
		SortT(compareUUIDs).
		SelectT(func(a uuid.UUID) string { return strings.Replace(a.String(), "-", "", 4) }).
		ToSlice(&organisationIds)

	organisationMap := make(map[string]int)

	for i, organisationId := range organisationIds {
		organisationMap[organisationId] = i
	}

	result := fmt.Sprintf("{%s}", strings.Join(organisationIds, ""))

	results := make(map[string]map[string]encodedPermission)

	linq.From(accessControlList).
		GroupByT(func(a AccessControlListEntry) string { return a.RecordType }, func(a AccessControlListEntry) AccessControlListEntry { return a }).
		ToMapBy(&results,
			func(k interface{}) interface{} { return k.(linq.Group).Key },
			func(v interface{}) interface{} {
				r := make(map[string]encodedPermission)
				linq.From(v.(linq.Group).Group).
					GroupByT(func(a AccessControlListEntry) string { return strings.Replace(a.OrganisationId.String(), "-", "", 4) },
						func(a AccessControlListEntry) AccessControlListEntry { return a }).
					ToMapBy(&r,
						func(ke interface{}) interface{} { return ke.(linq.Group).Key },
						func(va interface{}) interface{} {

							t := linq.From(va.(linq.Group).Group).SelectT(func(a AccessControlListEntry) AuthoriseAction { return a.Action }).
								AggregateT(func(a, b AuthoriseAction) AuthoriseAction { return a + b })

							rt := linq.From(va.(linq.Group).Group).First().(AccessControlListEntry)

							return encodedPermission{
								recordType:        rt.RecordType,
								actions:           t.(AuthoriseAction),
								organisationIndex: organisationMap[va.(linq.Group).Key.(string)],
							}
						})
				return r
			})

	var permissions []encodedPermission

	for _, encodedPermissions := range results {
		for _, encodedPermission := range encodedPermissions {
			permissions = append(permissions, encodedPermission)
		}
	}

	f := make(map[string]map[AuthoriseAction]string)

	organisationIdPrintFormat := "%0" + strconv.Itoa(len(strconv.Itoa(len(organisationIds)))) + "d"

	linq.From(permissions).
		GroupByT(func(e encodedPermission) string { return e.recordType }, func(e encodedPermission) encodedPermission { return e }).
		ToMapBy(&f,
			func(k interface{}) interface{} { return k.(linq.Group).Key },
			func(v interface{}) interface{} {
				r := make(map[AuthoriseAction]string)
				linq.From(v.(linq.Group).Group).
					GroupByT(func(a encodedPermission) AuthoriseAction { return a.actions },
						func(a encodedPermission) encodedPermission { return a }).
					ToMapBy(&r,
						func(ke interface{}) interface{} { return ke.(linq.Group).Key },
						func(va interface{}) interface{} {

							t := linq.From(va.(linq.Group).Group).
								SortT(func(a, b encodedPermission) bool { return a.organisationIndex < b.organisationIndex }).
								SelectT(func(a encodedPermission) string { return fmt.Sprintf(organisationIdPrintFormat, a.organisationIndex) }).
								AggregateT(func(a, b string) string { return a + b })

							return t
						})
				return r
			})

	sortedRecordTypes := make([]string, len(f))
	i := 0
	for k := range f {
		sortedRecordTypes[i] = k
		i++
	}
	sort.Strings(sortedRecordTypes)

	type kv struct {
		Key   AuthoriseAction
		Value string
	}

	for _, recordType := range sortedRecordTypes {
		p := f[recordType]
		result += fmt.Sprintf("%s[", recordType)
		j := 0

		var ss []kv
		for k, v := range p {
			ss = append(ss, kv{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value < ss[j].Value
		})

		for _, kv := range ss {
			result += fmt.Sprintf("%d:%s", kv.Key, kv.Value)
			j++
			if j != len(p) {
				result += "|"
			}
		}
		result += "]"
	}

	return compress(result)
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

func Decode(bytes []byte) (AccessControlList, error) {
	acl, err := decompress(bytes)

	if err != nil {
		return nil, err
	}

	results := AccessControlList{}

	organisationIdRegex, _ := regexp.Compile(`\{([a-f0-9]+)\}`)
	matches := organisationIdRegex.FindStringSubmatch(acl)

	if len(matches) < 2 {
		return results, nil
	}

	organisationIdMap := make(map[int]uuid.UUID)

	for i, u := range splitSubN(matches[1], 32) {
		organisationIdMap[i], _ = uuid.Parse(u)
	}

	organisationStrLen := len(fmt.Sprintf("%d", len(organisationIdMap)))

	remainingAcl := strings.Replace(acl, matches[0], "", 1)

	actionRegex, _ := regexp.Compile(`(\w+)\[(\d+:\d+(\|\d+:\d+)*)\]+`)

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

func compress(data string) ([]byte, error) {

	var b bytes.Buffer
	z := zlib.NewWriter(&b)

	if _, err := z.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := z.Flush(); err != nil {
		return nil, err
	}
	if err := z.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decompress(data []byte) (string, error) {

	rdata := bytes.NewReader(data)
	r, err := zlib.NewReader(rdata)

	if err != nil {
		return "", err
	}

	s, err := ioutil.ReadAll(r)

	if err != nil {
		return "", err
	}

	return string(s), nil
}
