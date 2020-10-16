package security

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

type Form3Claims struct {
	AclUrl *string `json:"https://form3.tech/acl_url"`
	UserId string  `json:"https://form3.tech/user_id"`
	jwt.StandardClaims
}

type jwtTokenBuilder struct {
	privateKey    *rsa.PrivateKey
	userId        uuid.UUID
	organisations []uuid.UUID
	permissions   []Permission
	aclUrl        *string
}

type Permission struct {
	AuthorisedActions []AuthoriseAction
	RecordType        string
}

func NewJwtToken(privateKey *rsa.PrivateKey, userId uuid.UUID) *jwtTokenBuilder {
	return &jwtTokenBuilder{
		privateKey:  privateKey,
		userId:      userId,
		permissions: []Permission{},
	}
}

func (b *jwtTokenBuilder) ForOrganisations(organisations ...uuid.UUID) *jwtTokenBuilder {
	b.organisations = organisations
	return b
}

func (b *jwtTokenBuilder) WithAclUrl(url string) *jwtTokenBuilder {
	b.aclUrl = &url
	return b
}

func (b *jwtTokenBuilder) Build() (string, error) {

	claims := Form3Claims{
		AclUrl: b.aclUrl,
		UserId: b.userId.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 600).Unix(),
			Issuer:    b.userId.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(b.privateKey)

	if err != nil {
		return "", fmt.Errorf("could not sign jwt token, error: %v", err)
	}

	return signedToken, nil
}

func EncodeAcls(organisations []uuid.UUID, permissions []Permission) (string, error) {
	accessControlList := AccessControlList{}

	for _, organisationId := range organisations {
		for _, permission := range permissions {
			for _, authorisedAction := range permission.AuthorisedActions {
				accessControlList = append(accessControlList, AccessControlListEntry{
					Action:         authorisedAction,
					OrganisationId: organisationId,
					RecordType:     permission.RecordType,
				})
			}
		}
	}

	buffer := &bytes.Buffer{}
	err := NewAclFullWithBase64Encoder().Encode(accessControlList, buffer)
	if err != nil {
		return "", fmt.Errorf("could not encode access control list, error: %v", err)
	}
	return buffer.String(), nil
}
