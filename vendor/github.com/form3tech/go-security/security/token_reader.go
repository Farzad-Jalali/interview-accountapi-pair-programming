package security

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
	r "github.com/dgrijalva/jwt-go/request"
)

const applicationContext string = "security.application_context"

type tokenReader struct {
	publicKey *rsa.PublicKey
}

var oneTokenReader sync.Once
var tokenReaderInstance *tokenReader

type JwtTokenParser interface {
	ParseTokenFromRequest(req *http.Request) (*context.Context, error)
}

func GetTokenReader(jwtPublicKey *rsa.PublicKey) JwtTokenParser {

	oneTokenReader.Do(func() {
		tokenReaderInstance = &tokenReader{
			publicKey: jwtPublicKey,
		}
	})

	return tokenReaderInstance
}

func (t *tokenReader) ParseTokenFromRequest(req *http.Request) (*context.Context, error) {

	token, err := r.ParseFromRequest(req, r.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not parse jwt token, error: %v", err)
	}

	var accessControlList AccessControlList
	userId := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		acls := claims["https://form3.tech/acl"].(string)
		bytes, _ := base64.StdEncoding.WithPadding('=').DecodeString(acls)

		if len(bytes) == 0 {
			accessControlList = AccessControlList{}
		} else {
			accessControlList, err = Decode(bytes)
		}

		if err != nil {
			return nil, fmt.Errorf("could not parse out acl claim to access control list, error: %v", err)
		}

		userId = claims["https://form3.tech/user_id"].(string)
	}

	ctx := context.WithValue(req.Context(), "acls", accessControlList)
	ctx = context.WithValue(ctx, "user_id", userId)

	return &ctx, nil
}

func ApplicationContext(ctx context.Context) *context.Context {
	result := context.WithValue(ctx, applicationContext, true)
	return &result
}

func IsApplicationContext(ctx context.Context) bool {
	ok, isApp := ctx.Value(applicationContext).(bool)
	return ok && isApp
}
