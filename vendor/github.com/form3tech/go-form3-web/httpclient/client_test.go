package httpclient

import (
	"os"
	"testing"
)

var clientId, clientSecret, apiHost string

func Test_accessing_api_that_requires_auth(t *testing.T) {

	AcceptanceTest(t, func() {
		given, when, then := HttpClientTest(t)

		given.
			client_credentials_with_access_to_the_api(clientId, clientSecret, apiHost)

		when.
			the_api_is_accessed()

		then.
			result_should_be_ok().and().
			no_error_should_be_raised()
	})

}

func Test_accessing_api_twice_that_requires_auth(t *testing.T) {

	AcceptanceTest(t, func() {
		given, when, then := HttpClientTest(t)

		given.
			client_credentials_with_access_to_the_api(clientId, clientSecret, apiHost)

		when.
			the_api_is_accessed().and().
			the_api_is_accessed()

		then.
			result_should_be_ok().and().
			no_error_should_be_raised()
	})

}

func AcceptanceTest(t *testing.T, test func()) {
	var ok bool
	clientId, ok = os.LookupEnv("CLIENT_ID")
	if !ok {
		t.Errorf("must provide env variable CLIENT_ID to run test")
	}

	clientSecret, ok = os.LookupEnv("CLIENT_SECRET")
	if !ok {
		t.Errorf("must provide env variable CLIENT_SECRET to run test")
	}

	apiHost, ok = os.LookupEnv("API_HOST")
	if !ok {
		t.Errorf("must provide env variable API_HOST to run test")
	}
	test()

}
