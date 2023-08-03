package http

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	tls "github.com/refraction-networking/utls"
)

// NewClient returns a new http client with reasonable timeouts and using a randomized tcp hello fingerprint
func NewClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, address)
			},
			TLSHandshakeTimeout: 1 * time.Second,
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

				//initialize the tcp connection
				tcpConn, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}

				//initialize the conifg for tls
				config := tls.Config{
					ServerName: strings.Split(addr, ":")[0], //set the server name with the provided addr
				}

				//initialize a tls connection with the underlying tcp connection and config
				//only HelloRandomizedNoALPN seems to consistently work
				tlsConn := tls.UClient(tcpConn, &config, tls.HelloRandomizedNoALPN)

				//start the tls handshake between servers
				err = tlsConn.Handshake()
				if err != nil {
					tcpConn.Close()
					return nil, err
				}

				return tlsConn, nil
			},
			ForceAttemptHTTP2: true,
		},
	}
}
