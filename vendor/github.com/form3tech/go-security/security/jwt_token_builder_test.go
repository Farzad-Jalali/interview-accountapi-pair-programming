package security

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"encoding/base64"
	"fmt"

	linq "github.com/ahmetb/go-linq"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_can_build_jwt_token(t *testing.T) {

	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)

	publicKey := key.PublicKey

	assert.Nil(t, err)

	userId := uuid.New()

	organisation1 := uuid.MustParse("35880d70-f983-4140-a7e0-6984e4ce683d")
	organisation2 := uuid.MustParse("ec529267-d2c7-49f9-961e-a65b7939b937")

	tokenString, err := NewJwtToken(key, userId).
		ForOrganisations(organisation1, organisation2).
		GivePermissions(CREATE, READ, EDIT_APPROVE).ToRecordType("Account").
		GivePermissions(DELETE, DELETE_APPROVE).ToRecordType("Payment").
		Build()

	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return &publicKey, nil
	})

	var claimUserId uuid.UUID
	var accessControlList AccessControlList

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimUserId = uuid.MustParse(claims["https://form3.tech/user_id"].(string))
		acls := claims["https://form3.tech/acl"].(string)
		bytes, err := base64.StdEncoding.WithPadding('=').DecodeString(acls)

		assert.Nil(t, err)

		accessControlList, err = Decode(bytes)

		assert.Nil(t, err)
	}

	assert.Equal(t, claimUserId, userId)
	assert.Equal(t, 10, len(accessControlList))

	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == CREATE && entry.OrganisationId == organisation1
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == READ && entry.OrganisationId == organisation1
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == EDIT_APPROVE && entry.OrganisationId == organisation1
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Payment" && entry.Action == DELETE && entry.OrganisationId == organisation1
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Payment" && entry.Action == DELETE_APPROVE && entry.OrganisationId == organisation1
	}).Count())

	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == CREATE && entry.OrganisationId == organisation2
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == READ && entry.OrganisationId == organisation2
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Account" && entry.Action == EDIT_APPROVE && entry.OrganisationId == organisation2
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Payment" && entry.Action == DELETE && entry.OrganisationId == organisation2
	}).Count())
	assert.Equal(t, 1, linq.From(accessControlList).WhereT(func(entry AccessControlListEntry) bool {
		return entry.RecordType == "Payment" && entry.Action == DELETE_APPROVE && entry.OrganisationId == organisation2
	}).Count())

}
