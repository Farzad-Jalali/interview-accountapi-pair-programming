package httpclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func NewConsulDnsHttpTransport(consulAddress string) (*http.Transport, error) {

	consulResolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (conn net.Conn, e error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "tcp", fmt.Sprintf("%s:8600", consulAddress))
		},
	}
	consulDialer := net.Dialer{
		Resolver: &consulResolver,
	}
	tlsConfig, err := trustForm3()
	if err != nil {
		return &http.Transport{}, err
	}
	tr := &http.Transport{
		DialContext:         consulDialer.DialContext,
		TLSHandshakeTimeout: 20 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	return tr, nil
}
