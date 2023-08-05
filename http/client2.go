package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

// ClientV2 is an http client which uses a randomized tcp hello fingerprint
type ClientV2 struct {
	roundTripper *roundTripper
}

// NewClientV2 returns a new http client with using a randomized tcp hello fingerprint
func NewClientV2() *ClientV2 {
	return &ClientV2{
		roundTripper: newRoundTripper(),
	}
}

// Do implements the http.Client interface
func (c *ClientV2) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.roundTripper.RoundTrip(req)
	if isConnectionBroken(err) {
		c.roundTripper.resetConnection()
		return c.roundTripper.RoundTrip(req)
	}
	return resp, err
}

func isConnectionBroken(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "connection broken") || strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "connection refused") {
		return true
	}

	return false
}

// SetClientHelloID sets the ClientHelloID to use when establishing a new connection in case you want to override the default
func (c *ClientV2) SetClientHelloID(id utls.ClientHelloID) {
	c.roundTripper.setHelloID(id)
}

type roundTripper struct {
	sync.Mutex

	transport http.RoundTripper
	conn      net.Conn
	dialer    *net.Dialer
	helloID   utls.ClientHelloID
}

// resetConnection resets the connection so that the next request will establish a new connection
func (rt *roundTripper) resetConnection() {
	rt.Lock()
	defer rt.Unlock()

	rt.transport = nil
	rt.conn = nil
}

// RoundTrip implements the http.RoundTripper interface
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport, err := rt.getTransport(req)
	if err != nil {
		return nil, err
	}

	return transport.RoundTrip(req)
}

func (rt *roundTripper) getTransport(req *http.Request) (http.RoundTripper, error) {
	rt.Lock()
	defer rt.Unlock()

	if rt.transport != nil {
		return rt.transport, nil
	}

	switch strings.ToLower(req.URL.Scheme) {
	case "http":
		rt.transport = http.DefaultClient.Transport
		return rt.transport, nil
	case "https":
	default:
		return nil, fmt.Errorf("invalid URL scheme: '%v'", req.URL.Scheme)
	}

	if _, err := rt.dialTLS("tcp", getAddrFromURL(req.URL)); err != nil {
		return nil, err
	}

	return rt.transport, nil
}

func (rt *roundTripper) dialTLS(network, addr string) (net.Conn, error) {
	// Check if we have a connection already
	if conn := rt.conn; conn != nil {
		rt.conn = nil
		return conn, nil
	}

	rawConn, err := rt.dialer.DialContext(context.Background(), network, addr)
	if err != nil {
		return nil, err
	}

	var host string
	if host, _, err = net.SplitHostPort(addr); err != nil {
		host = addr
	}

	// Initialize the tls connection
	conn := utls.UClient(rawConn, &utls.Config{ServerName: host}, rt.helloID)
	if err = conn.Handshake(); err != nil {
		conn.Close()
		rt.resetConnection()
		return nil, err
	}

	if rt.transport != nil {
		return conn, nil
	}

	// Check which protocol was negotiated
	switch conn.ConnectionState().NegotiatedProtocol {
	case http2.NextProtoTLS:
		// The remote peer is speaking HTTP 2 + TLS.
		rt.transport = &http2.Transport{DialTLS: rt.dialTLSHTTP2}
	default:
		// Assume the remote peer is speaking HTTP 1.x + TLS.
		rt.transport = &http.Transport{DialTLS: rt.dialTLS}
	}

	// Save the connection for next time
	rt.conn = conn

	return nil, nil
}

func getAddrFromURL(u *url.URL) string {
	host, port, err := net.SplitHostPort(u.Host)
	if err == nil {
		return net.JoinHostPort(host, port)
	}

	return net.JoinHostPort(u.Host, u.Scheme)
}

func (rt *roundTripper) dialTLSHTTP2(network, addr string, cfg *tls.Config) (net.Conn, error) {
	return rt.dialTLS(network, addr)
}

func (rt *roundTripper) setHelloID(id utls.ClientHelloID) {
	rt.Lock()
	defer rt.Unlock()
	rt.helloID = id
}

func newRoundTripper() *roundTripper {
	return &roundTripper{
		dialer: &net.Dialer{
			Timeout:   1 * time.Second,
			KeepAlive: 30 * time.Second,
		},

		//this one is the most stable, and it seems to use HTTP 1.1 which gets detected with a higher botscore
		helloID: utls.HelloRandomizedNoALPN,
	}
}

func init() {
	utls.EnableWeakCiphers()
}
