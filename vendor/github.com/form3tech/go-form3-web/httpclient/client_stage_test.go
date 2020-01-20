package httpclient

import (
	"context"
	"net/url"
	"testing"

	"github.com/form3tech/go-form3-web/httpclient/support/swagger-client/paymentapi/client"
	"github.com/form3tech/go-form3-web/httpclient/support/swagger-client/paymentapi/client/payments"
	rc "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

type httpClientStage struct {
	t                   *testing.T
	clientConfig        *ClientConfig
	getPaymentsResponse *payments.GetPaymentsOK
	error               error
}

func HttpClientTest(t *testing.T) (*httpClientStage, *httpClientStage, *httpClientStage) {

	stage := &httpClientStage{
		t: t,
	}

	return stage, stage, stage
}

func newPaymentClient(c *ClientConfig) *payments.Client {

	config := client.DefaultTransportConfig().
		WithHost(c.HostUrl.Host).
		WithBasePath("/v1/transaction").
		WithSchemes([]string{c.HostUrl.Scheme})

	rt := rc.NewWithClient(config.Host, config.BasePath, config.Schemes, NewHttpClient(c))

	return payments.New(rt, strfmt.Default)
}

func (s *httpClientStage) client_credentials_with_access_to_the_api(clientId, clientSecret, apiHost string) *httpClientStage {

	u, _ := url.Parse(apiHost)

	userId := "40dbef38-b747-411f-b7bf-021c818aa0d1"
	s.clientConfig = NewClientConfig(clientId, clientSecret, userId, u)

	return s
}

func (s *httpClientStage) the_api_is_accessed() *httpClientStage {

	paymentClient := newPaymentClient(s.clientConfig)

	s.getPaymentsResponse, s.error = paymentClient.GetPayments(&payments.GetPaymentsParams{Context: context.Background()})

	return s
}

func (s *httpClientStage) result_should_be_ok() *httpClientStage {

	assert.True(s.t, len(s.getPaymentsResponse.Payload.Data) >= 0)

	return s
}

func (s *httpClientStage) and() *httpClientStage {
	return s
}

func (s *httpClientStage) no_error_should_be_raised() *httpClientStage {

	assert.Nil(s.t, s.error)

	return s
}
