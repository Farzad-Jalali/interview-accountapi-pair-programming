package security

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/form3tech/go-logger/log"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	r "github.com/form3tech-oss/jwt-go/request"
)

var jwtStrategyDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name: "access_control_strategy_time",
	Help: "Time taken to process a jwt token",
})

type contextKey int

const (
	contextKeyUserID contextKey = iota
	contextKeyACLs
	contextKeyApplication
)

type tokenReader struct {
	publicKey *rsa.PublicKey
	decoder   Decoder
}

type TokenReaderOptions func(*tokenReader)

func WithDefaultAclDecoder() TokenReaderOptions {
	return func(reader *tokenReader) {
		reader.decoder = NewAclFullWithBase64Decoder()
	}
}

func WithAclCompactDecoder() TokenReaderOptions {
	return func(reader *tokenReader) {
		reader.decoder = NewAclCompactWithBase64Decoder()
	}
}

var oneTokenReader sync.Once
var tokenReaderInstance *tokenReader
var logRate float32 = 0

var defaultTokenReaderOpts = []TokenReaderOptions{
	WithDefaultAclDecoder(),
}

func init() {
	logRateString, exists := os.LookupEnv("ACCESS_CONTROL_STRATEGY_LOG_RATE")
	if !exists {
		return
	}

	rate, err := strconv.ParseFloat(logRateString, 32)
	if err != nil {
		return
	}
	logRate = float32(rate)
	logrus.New().Infof("ACCESS_CONTROL_STRATEGY_LOG_RATE=%f", logRate)
}

type JwtTokenParser interface {
	ParseTokenFromRequest(req *http.Request) (*context.Context, error)
}

func GetTokenReader(jwtPublicKey *rsa.PublicKey, opts ...TokenReaderOptions) JwtTokenParser {
	oneTokenReader.Do(func() {
		tokenReaderInstance = &tokenReader{
			publicKey: jwtPublicKey,
			decoder:   NewAclFullWithBase64Decoder(),
		}
		for _, opt := range append(defaultTokenReaderOpts, opts...) {
			opt(tokenReaderInstance)
		}
		rand.Seed(time.Now().UnixNano())
	})
	return tokenReaderInstance
}

func (t *tokenReader) ParseTokenFromRequest(req *http.Request) (*context.Context, error) {
	timer := prometheus.NewTimer(jwtStrategyDuration)
	defer timer.ObserveDuration()

	token, err := r.ParseFromRequest(req, r.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse jwt token, error: %v", err)
	}

	var accessControlList AccessControl
	var userId string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessControlList, err = aclsFromClaims(t.decoder, claims)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse out acl claim to access control list")
		}

		userId = claims["https://form3.tech/user_id"].(string)
	}

	ctx := context.WithValue(req.Context(), contextKeyACLs, accessControlList)
	ctx = context.WithValue(ctx, contextKeyUserID, userId)

	return &ctx, nil
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(contextKeyUserID).(string)
	if !ok {
		return "", fmt.Errorf("user id is not available in provided context")
	}
	return userID, nil
}

func ApplicationContext(ctx context.Context) *context.Context {
	result := context.WithValue(ctx, contextKeyApplication, true)
	return &result
}

func IsApplicationContext(ctx context.Context) bool {
	ok, isApp := ctx.Value(contextKeyApplication).(bool)
	return ok && isApp
}

func aclsFromClaims(dec Decoder, claims jwt.MapClaims) (AccessControl, error) {
	aclUrlClaim, ok := claims["https://form3.tech/acl_url"]
	if !ok || aclUrlClaim == nil {
		return nil, fmt.Errorf("claim 'https://form3.tech/acl_url' is missing")
	}

	aclUrl, ok := aclUrlClaim.(string)
	if !ok {
		return nil, fmt.Errorf("claim 'https://form3.tech/acl_url' is invalid type")
	}

	return fetchAcls(dec, aclUrl)
}

func fetchAcls(dec Decoder, urlString string) (AccessControl, error) {
	aclUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to parse acl_url '%s'", urlString))
	}

	randomLog(fmt.Sprintf("Fetching ACLs from URI '%s'", aclUrl))

	resp, err := http.Get(aclUrl.String())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("request for acl_url failed '%s'", urlString))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return handleSuccessAclUrlResponse(dec, resp.Body)
	case http.StatusNotFound:
		return NilAccessControl, nil
	}
	return nil, fmt.Errorf("non 2xx response calling acl_url %s: %d", urlString, resp.StatusCode)
}

func handleSuccessAclUrlResponse(dec Decoder, body io.Reader) (AccessControl, error) {
	acls, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read request body")
	}
	return dec.Decode(acls)
}

func randomLog(message string) {
	if logRate <= 0 {
		return
	}
	randMax := int(1 / logRate)
	/* #nosec G404 */
	if rand.Intn(randMax) == 0 {
		log.Info(message)
	}
}
