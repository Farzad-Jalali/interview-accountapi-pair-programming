package security

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
)

type encodedPermission struct {
	recordType        string
	actions           AuthoriseAction
	organisationIndex int
}

type aclFullEncoder struct{}

var _ Encoder = &aclFullEncoder{}

func NewAclFullEncoder() Encoder {
	return newCompressEncoder(&aclFullEncoder{})
}

func NewAclFullWithBase64Encoder() Encoder {
	return newBase64Encoder(NewAclFullEncoder())
}

func (a *aclFullEncoder) Encode(acls AccessControlList, writer io.Writer) error {

	var organisationIds []string

	linq.From(acls).
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

	linq.From(acls).
		GroupByT(func(a AccessControlListEntry) string { return a.RecordType }, func(a AccessControlListEntry) AccessControlListEntry { return a }).
		ToMapBy(&results,
			func(k interface{}) interface{} { return k.(linq.Group).Key },
			func(v interface{}) interface{} {
				r := make(map[string]encodedPermission)
				linq.From(v.(linq.Group).Group).
					GroupByT(func(a AccessControlListEntry) string {
						return strings.Replace(a.OrganisationId.String(), "-", "", 4)
					},
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

	_, err := writer.Write([]byte(result))
	return err
}

func compareUUIDs(a, b uuid.UUID) bool {
	msb1, lsb1 := uuidToSigBits(a)
	msb2, lsb2 := uuidToSigBits(b)

	if msb1 == msb2 {
		return lsb1 < lsb2
	}
	return msb1 < msb2
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
