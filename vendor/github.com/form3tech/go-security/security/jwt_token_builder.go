package security

import (
	"crypto/rsa"
	"fmt"

	"encoding/base64"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Form3Claims struct {
	Acl    string `json:"https://form3.tech/acl"`
	UserId string `json:"https://form3.tech/user_id"`
	jwt.StandardClaims
}

type jwtTokenBuilder struct {
	privateKey    *rsa.PrivateKey
	userId        uuid.UUID
	organisations []uuid.UUID
	permissions   []permission
}

type permission struct {
	authorisedActions []AuthoriseAction
	recordType        string
}

type permissionBuilder struct {
	jwtTokenBuilder   *jwtTokenBuilder
	authorisedActions []AuthoriseAction
}

func NewJwtToken(privateKey *rsa.PrivateKey, userId uuid.UUID) *jwtTokenBuilder {
	return &jwtTokenBuilder{
		privateKey:  privateKey,
		userId:      userId,
		permissions: []permission{},
	}
}

func (b *jwtTokenBuilder) ForOrganisations(organisations ...uuid.UUID) *jwtTokenBuilder {
	b.organisations = organisations
	return b
}

func (b *jwtTokenBuilder) GivePermissions(authorisedActions ...AuthoriseAction) *permissionBuilder {
	return &permissionBuilder{
		jwtTokenBuilder:   b,
		authorisedActions: authorisedActions,
	}
}

func (p *permissionBuilder) ToRecordType(recordType string) *jwtTokenBuilder {

	p.jwtTokenBuilder.permissions = append(p.jwtTokenBuilder.permissions, permission{recordType: recordType, authorisedActions: p.authorisedActions})

	return p.jwtTokenBuilder
}

func (b *jwtTokenBuilder) Build() (string, error) {

	accessControlList := AccessControlList{}

	for _, organisationId := range b.organisations {
		for _, permission := range b.permissions {
			for _, authorisedAction := range permission.authorisedActions {
				accessControlList = append(accessControlList, AccessControlListEntry{
					Action:         authorisedAction,
					OrganisationId: organisationId,
					RecordType:     permission.recordType,
				})
			}
		}
	}

	encodedAcls, err := Encode(accessControlList)

	if err != nil {
		return "", fmt.Errorf("could not encode access control list, error: %v", err)
	}

	claims := Form3Claims{
		Acl:    base64.StdEncoding.WithPadding('=').EncodeToString(encodedAcls),
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
