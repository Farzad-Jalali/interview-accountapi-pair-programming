package security

import (
	"testing"

	"encoding/base64"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_EncodeAccessControlListWithSingleOrganisation(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778ae}Account[243:0]", aclList)

}

func Test_EncodeAccessControlListWithMultipleOrganisationsSamePermissions(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778a1960e070e878d442bb638d0aa33f778ae}Account[243:01]", aclList)

}

func Test_EncodeAccessControlListWithMultipleOrganisationsDifferentPermissions(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778a1960e070e878d442bb638d0aa33f778ae}Account[129:0|243:1]", aclList)
}

func Test_EncodeAccessControlListWithComplexPermissions(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: DELETE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: DELETE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: REJECT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: REJECT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Submission"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778aee0d2496a525c436fb8e0a8838075b11404e2d9b88d094c4e829f78b8170933544dce3bcdecf84f2bbe4e56d04aa292c0}Account[1023:0]Payment[1:23]Submission[64:1]", aclList)

}

func Test_EncodeAccessControlListWithNonAlphabeticalOrderRecordTypes(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Zoo"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Submission"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778aee0d2496a525c436fb8e0a8838075b11404e2d9b88d094c4e829f78b8170933544dce3bcdecf84f2bbe4e56d04aa292c0}Payment[1:23]Submission[64:1]Zoo[64:0]", aclList)
}

func Test_EncodeAccessControlListWithTenOrganisations(t *testing.T) {

	accessControlList := AccessControlList{
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9cee8dde-f1b1-46d9-b2a2-a5a95a8d0239"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("ce48c314-5fa0-4c1d-aa14-e3dc3f935528"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9af0d5f8-65e8-4746-a237-c746106c342c"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("e929e690-ca56-4520-bc25-174b29694b75"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9b2e011f-2ab8-4d3d-9aa6-1da8a12e3927"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("1a1ef895-1458-48ea-a0ca-472debc20bc7"), RecordType: "Payment"},
	}

	compressed, err := Encode(accessControlList)

	assert.Nil(t, err)

	aclList, err := decompress(compressed)

	assert.Nil(t, err)
	assert.Equal(t, "{960e070e878d442bb638d0aa33f778ae9af0d5f865e84746a237c746106c342c9b2e011f2ab84d3d9aa61da8a12e39279cee8ddef1b146d9b2a2a5a95a8d0239ce48c3145fa04c1daa14e3dc3f935528e0d2496a525c436fb8e0a8838075b114e929e690ca564520bc25174b29694b7504e2d9b88d094c4e829f78b8170933541a1ef895145848eaa0ca472debc20bc74dce3bcdecf84f2bbe4e56d04aa292c0}Payment[1:00010203040506070809]", aclList)

}

// - decode

func IgnoreError(b []byte, e error) []byte {
	if e != nil {
		panic(e)
	}
	return b
}

func Test_DecodeEmptyAccessControlList(t *testing.T) {

	acl := "{}"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)

}

func Test_DecodeAccessControlListWithMalformedOrgID(t *testing.T) {

	acl := "{960g070e878d442bb638d0aa33f778ae}Account[243:0]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)

}

func Test_DecodeAccessControlListWithSingleOrganisation(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778ae}Account[243:0]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)

}

func Test_DecodeAccessControlListWithMultipleOrganisationsSamePermissions(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778a1960e070e878d442bb638d0aa33f778ae}Account[243:01]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)
}

func Test_DecodeAccessControlListWithMultipleOrganisationsDifferentPermissions(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778a1960e070e878d442bb638d0aa33f778ae}Account[129:0|243:1]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778a1"), RecordType: "Account"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)
}

func Test_DecodeAccessControlListWithComplexPermissions(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778aee0d2496a525c436fb8e0a8838075b11404e2d9b88d094c4e829f78b8170933544dce3bcdecf84f2bbe4e56d04aa292c0}Account[1023:0]Payment[1:23]Submission[64:1]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: CREATE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: EDIT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: DELETE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: DELETE_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: REJECT, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},
		AccessControlListEntry{Action: REJECT_APPROVE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Account"},

		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Submission"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)
}

func Test_DecodeAccessControlListWithNonAlphabeticalOrderRecordTypes(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778aee0d2496a525c436fb8e0a8838075b11404e2d9b88d094c4e829f78b8170933544dce3bcdecf84f2bbe4e56d04aa292c0}Payment[1:23]Submission[64:1]Zoo[64:0]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Zoo"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: READ, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Submission"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)
}

func Test_DecodeAccessControlListWithTenOrganisations(t *testing.T) {

	acl := "{960e070e878d442bb638d0aa33f778ae9af0d5f865e84746a237c746106c342c9b2e011f2ab84d3d9aa61da8a12e39279cee8ddef1b146d9b2a2a5a95a8d0239ce48c3145fa04c1daa14e3dc3f935528e0d2496a525c436fb8e0a8838075b114e929e690ca564520bc25174b29694b7504e2d9b88d094c4e829f78b8170933541a1ef895145848eaa0ca472debc20bc74dce3bcdecf84f2bbe4e56d04aa292c0}Payment[1:04020103070508090600]"
	aclList, err := Decode(IgnoreError(compress(acl)))

	expectedAccessControlList := AccessControlList{
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("960e070e-878d-442b-b638-d0aa33f778ae"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("4dce3bcd-ecf8-4f2b-be4e-56d04aa292c0"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("e0d2496a-525c-436f-b8e0-a8838075b114"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("04e2d9b8-8d09-4c4e-829f-78b817093354"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9cee8dde-f1b1-46d9-b2a2-a5a95a8d0239"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("ce48c314-5fa0-4c1d-aa14-e3dc3f935528"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9af0d5f8-65e8-4746-a237-c746106c342c"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("e929e690-ca56-4520-bc25-174b29694b75"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("9b2e011f-2ab8-4d3d-9aa6-1da8a12e3927"), RecordType: "Payment"},
		AccessControlListEntry{Action: CREATE, OrganisationId: uuid.MustParse("1a1ef895-1458-48ea-a0ca-472debc20bc7"), RecordType: "Payment"},
	}

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedAccessControlList, aclList)
}

func Test_DecodeBase64EncodedAclsZippedByJava(t *testing.T) {

	acl := `eJx9kt1qwzAMhV8pPzSlvctgu11J2W7KCIqtDoNrB8splLF3X5rYdWLT3Vk+R9In2T+IRXZGhF1VsoxVvMphywssOWzOebdlvzVjelD2lGdFuc++PgiNP/dwu6CyLQ3dRRAJrdorSMHBjkc65bsimMgnNWgHo2ruMpyrwSsaApne391edZeHueT/3lg1k9oaJ3u849ARM6K3k3NGdFbwFaJRFvNSSFmXrRl6Kab9fKxoxT0GmzF4N9+gBMGS56BJ2HRXxwdG1CoISS/PmQ73JiS+Kt5rER57Jlv2mRgbLTHa1fO6jimU9E+5/kPP8tIhX0CCuq93sv0Bz5z5Xg==`

	decodedStr, _ := base64.StdEncoding.WithPadding('=').DecodeString(acl)

	aclList, err := Decode(decodedStr)

	assert.Nil(t, err)
	assert.Equal(t, 170, len(aclList))
}
