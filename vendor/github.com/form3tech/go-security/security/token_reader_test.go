package security

import (
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var keyPairs, _ = GenerateTestKeyPair()

func Test_read_token_from_context(t *testing.T) {

	userId := uuid.New()

	jwtToken, err := NewJwtToken(keyPairs.RsaPrivateKey, userId).
		ForOrganisations(uuid.New()).
		GivePermissions(CREATE, READ, EDIT, DELETE).ToRecordType(string("abcd")).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	reader := GetTokenReader(&keyPairs.RsaPrivateKey.PublicKey)

	ctx, err := reader.ParseTokenFromRequest(&http.Request{
		Header: http.Header{
			"Authorization": []string{"BEARER " + jwtToken},
		},
	})

	assert.Nil(t, err)

	accessControlList := (*ctx).Value("acls").(AccessControlList)

	assert.Equal(t, 4, len(accessControlList))
	assert.Equal(t, userId.String(), (*ctx).Value("user_id"))
}

func Test_read_token_with_no_acls(t *testing.T) {

	userId := uuid.New()

	claims := Form3Claims{
		Acl:    base64.StdEncoding.WithPadding('=').EncodeToString([]byte{}),
		UserId: userId.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 600).Unix(),
			Issuer:    userId.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtToken, err := token.SignedString(keyPairs.RsaPrivateKey)

	assert.NoError(t, err)

	reader := GetTokenReader(&keyPairs.RsaPrivateKey.PublicKey)

	ctx, err := reader.ParseTokenFromRequest(&http.Request{
		Header: http.Header{
			"Authorization": []string{"BEARER " + jwtToken},
		},
	})

	assert.NoError(t, err)

	accessControlList := (*ctx).Value("acls").(AccessControlList)

	assert.Equal(t, 0, len(accessControlList))
	assert.Equal(t, userId.String(), (*ctx).Value("user_id"))
}
