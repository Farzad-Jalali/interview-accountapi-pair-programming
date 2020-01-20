package httpclient

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type ClientConfig struct {
	ClientId            string
	ClientSecret        string
	UserId              string
	HostUrl             *url.URL
	InitialToken        string
	authEndpoint        string
	UnderlyingTransport http.RoundTripper
}

type transport struct {
	config              *ClientConfig
	token               string
	underlyingTransport http.RoundTripper
}

type authJsonResponse struct {
	AccessToken string `json:"access_token"`
}

func NewClientConfig(clientId, clientSecret, userId string, hostUrl *url.URL) *ClientConfig {
	return &ClientConfig{
		ClientId:            clientId,
		ClientSecret:        clientSecret,
		UserId:              userId,
		HostUrl:             hostUrl,
		authEndpoint:        "/v1/oauth2/token",
		UnderlyingTransport: http.DefaultTransport,
	}
}

func (c *ClientConfig) WithInternalAuth() *ClientConfig {
	c.authEndpoint = "/securityapi/v1/oauth2/token"
	return c
}

func (c *ClientConfig) WithInitialToken(token string) *ClientConfig {
	c.InitialToken = token
	return c
}

func (c *ClientConfig) WithUnderlyingTransport(underlyingTransport http.RoundTripper) *ClientConfig {
	c.UnderlyingTransport = underlyingTransport
	return c
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}

	req.Header.Add("X-consumer-custom-id", t.config.UserId)

	res, err := t.underlyingTransport.RoundTrip(req)

	if res != nil && (res.StatusCode == 401 || res.StatusCode == 403) {

		authRequest, err := http.NewRequest("POST", t.config.HostUrl.String()+t.config.authEndpoint, nil)
		authRequest.SetBasicAuth(t.config.ClientId, t.config.ClientSecret)

		if err != nil {
			return nil, fmt.Errorf("could not build auth request, error: %v", err)
		}

		authResponse, err := t.underlyingTransport.RoundTrip(authRequest)

		if err != nil {
			return nil, fmt.Errorf("error authenticating, error: %v", err)
		}

		if authResponse.StatusCode != 200 {
			return nil, fmt.Errorf("non 200 status code getting token, status code: %d", authResponse.StatusCode)
		}

		authBody, err := ioutil.ReadAll(authResponse.Body)

		if err != nil {
			return nil, fmt.Errorf("could not read auth response body, error: %v", err)
		}

		var authJson authJsonResponse

		err = json.Unmarshal(authBody, &authJson)

		if err != nil {
			return nil, fmt.Errorf("could not parse auth json response, error: %v", err)
		}

		t.token = authJson.AccessToken

		req.Header.Del("Authorization")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

		return t.underlyingTransport.RoundTrip(req)
	}

	return res, err

}

func NewHttpClient(config *ClientConfig) *http.Client {
	transport := &transport{underlyingTransport: config.UnderlyingTransport, config: config}
	transport.token = config.InitialToken

	h := &http.Client{Transport: transport}

	return h
}
