package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ChromeClient is an http client which uses a headless chrome browser to make requests
// Note that this client requires a chrome browser to be installed on the machine
type ChromeClient struct{}

// Do performs the http request using a headless chrome browser
func (c *ChromeClient) Do(req *http.Request) (*http.Response, error) {
	if req.Method != http.MethodGet {
		return nil, fmt.Errorf("only GET requests are supported")
	}

	dir, err := os.MkdirTemp("", "chromedp-tmp")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.UserAgent(req.UserAgent()),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel1()

	// also set up a custom logger
	taskCtx, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel2()

	// Set up a timeout to prevent the script from running indefinitely
	ctx, cancel3 := context.WithTimeout(taskCtx, 10*time.Second)
	defer cancel3()

	// Set up custom headers
	headers := make(map[string]interface{})
	for k, v := range req.Header {
		headers[k] = strings.Join(v, ",")
	}

	//Listen for the response
	var (
		statusCode      int64
		responseHeaders network.Headers
	)
	url := req.URL.String()

	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch responseReceivedEvent := event.(type) {
		case *network.EventResponseReceived:
			response := responseReceivedEvent.Response
			if response.URL == url {
				statusCode = response.Status
				responseHeaders = response.Headers
			}
		}
	})

	// Navigate to the website
	var body string
	err = chromedp.Run(ctx,
		chromedp.Tasks{
			network.Enable(),
			network.SetExtraHTTPHeaders(network.Headers(headers)),
			chromedp.Navigate(url),

			&myQueryAction{&body, &responseHeaders},
		})
	if err != nil {
		return nil, err
	}

	return buildResponse(statusCode, responseHeaders, body), nil
}

// buildResponse builds an http response from the given status code, headers, and body
func buildResponse(statusCode int64, h network.Headers, body string) *http.Response {
	headers := make(map[string][]string)
	for k, v := range h {
		headers[k] = []string{v.(string)}
	}

	resp := &http.Response{
		StatusCode:    int(statusCode),
		Header:        headers,
		Body:          io.NopCloser(strings.NewReader(string(body))),
		ContentLength: int64(len(body)),
	}

	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	return resp
}

func isJsonResp(h network.Headers) bool {
	for k, v := range h {
		if strings.ToLower(k) == "content-type" {
			return strings.Contains(strings.ToLower(v.(string)), "json")
		}
	}
	return false
}

type myQueryAction struct {
	body            *string
	responseHeaders *network.Headers
}

func (m *myQueryAction) Do(ctx context.Context) error {
	if isJsonResp(*m.responseHeaders) {
		return chromedp.InnerHTML(`pre`, m.body, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
	}
	return chromedp.OuterHTML(`html`, m.body, chromedp.NodeVisible, chromedp.ByQuery).Do(ctx)
}

// GetCookies performs a GET request using a headless Chrome browser, waits for any
// JavaScript challenge (e.g. AWS WAF Bot Control) to complete, and returns all cookies
// set in the browser after the challenge resolves.
func (c *ChromeClient) GetCookies(req *http.Request) ([]*http.Cookie, error) {
	if req.Method != http.MethodGet {
		return nil, fmt.Errorf("only GET requests are supported")
	}

	dir, err := os.MkdirTemp("", "chromedp-tmp")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.UserAgent(req.UserAgent()),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel1()

	taskCtx, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel2()

	// Longer timeout to accommodate JS challenge + page reload
	ctx, cancel3 := context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel3()

	headers := make(map[string]interface{})
	for k, v := range req.Header {
		headers[k] = strings.Join(v, ",")
	}

	var networkCookies []*network.Cookie
	err = chromedp.Run(ctx,
		chromedp.Tasks{
			network.Enable(),
			network.SetExtraHTTPHeaders(network.Headers(headers)),
			chromedp.Navigate(req.URL.String()),
			// Wait for WAF JS challenge to run and trigger window.location.reload()
			chromedp.Sleep(5 * time.Second),
			chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				networkCookies, err = network.GetCookies().Do(ctx)
				return err
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	cookies := make([]*http.Cookie, 0, len(networkCookies))
	for _, nc := range networkCookies {
		cookies = append(cookies, &http.Cookie{
			Name:  nc.Name,
			Value: nc.Value,
		})
	}
	return cookies, nil
}
