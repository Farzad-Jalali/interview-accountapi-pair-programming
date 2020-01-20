# go-form3-web
Form3 http client for accessing the form3 api


## Usage
`httpclient` provides a package that gives you an authenticating http client.  The `NewHttpClient` func gives you a wrapped `*http.Client` that will authenticate (get a token) automatically if
the call is not authorised (401 or 403).  The client will then automatically make the original call again.  You don't have to worry about authentication with this client just set your client id
and client secret and be done.  The client will handle authentication.

First you will need to create a `*clientConfig` this is where you pass in your client id, secret and the host address of the api:

```
u, _ := url.Parse("https://api.form3.tech")

clientConfig := NewClientConfig("your client id", "your client secret", u)
```

Once you have `*clientConfig` you can simply use the `NewHttpClient` func to generate a new `*http.Client`:
```
client := NewHttpClient(clientConfig)
```

Here is an example of how you can plug it into a swagger generated client:
 ```
import rc "github.com/go-openapi/runtime/client"


func newPaymentClient(c *clientConfig) *payments.Client {

	var config *client.TransportConfig

	config = client.DefaultTransportConfig().
		WithHost(c.hostUrl.Host).
		WithBasePath("/v1/transaction").
		WithSchemes([]string{c.hostUrl.Scheme})

	rt := rc.NewWithClient(config.Host, config.BasePath, config.Schemes, NewHttpClient(c))

	return payments.New(rt, strfmt.Default)
}
```

## Running tests
To run the `httpclient` tests you need to have the following env variables defined:
```
CLIENT_ID       - your client id
CLIENT_SECRET   - your secret
API_HOST        - (including scheme)
```
