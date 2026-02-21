# Hulk Smash [![GoDoc](https://godoc.org/github.com/ynori7/hulksmash?status.png)](https://godoc.org/github.com/ynori7/hulksmash) [![Go Report Card](https://goreportcard.com/badge/ynori7/hulksmash)](https://goreportcard.com/report/github.com/ynori7/hulksmash)

![HulkSmash Logo](hulksmash.png)

This is a very easy-to-use library for building a custom brute-force requester for QA purposes. This
tool can be useful, for example, for performance and load testing or for testing your rate-limiter. This tool
automatically adds randomized headers to anonymize the request such as `X-Forwarded-For` and `User-Agent`.

This tool is to be used only for benign purposes! 

## Usage

### Smasher
The smasher is the bulk requester which will perform repeated requests. It can be configured to perform requests in parallel and it has options for generating sequences. 

To use it, simply import `"github.com/ynori7/hulksmash"`. Then construct a request builder, which is a function to get the request you want to perform. This request builder accepts an index in case you want to send requests to a variety of endpoints, with varying payloads, or with a cachebreaker. Then you simply create your Smasher instance and tell it to start smashing.

```go
requestBuilder := func(item string) (*http.Request, error) {
    return http.NewRequest(http.MethodGet, url, nil)
}

hulksmash.NewSmasher(hulksmash.WithIterations(5)).Smash(context.Background(), requestBuilder)
```

The Smasher comes with some configurable options with safe defaults. Here is a list of the options:

| Option        | Description           | Default  |
| ------------- |:-------------| -----|
| WithClient     | Allows you to override the HTTP Client | hulksmash's http.NewClientV2() |
| WithWorkerCount      | Sets the number of workers which will send requests in parallel      |   1 |
| WithErrorFunc | Function which is called in case of an error while performing the request      |    Simply logs it to stdout |
| WithSuccessResponseCallback | Function which is called in case of a successful request | Simply logs the http status code and response body |
| WithIterations | The number of calls to make | 1 |
| WithStartIndex | The start index to use when iterating. Can be useful if you want to resume a previous experiment | 0 |
| WithAnonymizeRequests | Can be used to disable the logic to automatically add headers to make reqeusts look more organic | true |
| WithSequenceFunc | Can be used to specify the way the iteration sequence should be built, for example numeric or alphanumeric. Some presets are available in the sequence package | sequence.Numeric |

A basic and advanced example can be found in [example](example). Internally, the Smasher automatically uses hulksmash's HTTP ClientV2 (unless you override it with a custom client) and the request Anonymizer (if enabled).

### Anonymizer
The anonymizer is used to add randomized headers to HTTP requests to make them look more organic and less like bot traffic. It generates random IP addresses (both IPv4 and IPv6), user agent strings for various browsers, and accept-language headers.

```go
import "github.com/ynori7/hulksmash/anonymizer"

// Create a new anonymizer with a seed for reproducibility
anon := anonymizer.New(time.Now().Unix())

// Anonymize an existing request (adds headers automatically)
req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
anon.AnonymizeRequest(req)

// Or generate random values independently
randomIP := anon.GetRandomIp()  // Returns a random IPv4 or IPv6 address
randomUA := anon.GetRandomUserAgent()  // Returns a random user agent string
randomLang := anon.GetRandomAcceptLanguage()  // Returns a random accept-language header

// You can also specify a specific browser for the user agent
chromeUA := anon.GetRandomUserAgentWithBrowser(anonymizer.Chrome)
firefoxUA := anon.GetRandomUserAgentWithBrowser(anonymizer.Firefox)
```

The `AnonymizeRequest` method automatically sets the following headers:
- `X-Forwarded-For`: Random IP address
- `User-Agent`: Random browser user agent string
- `Cache-Control`: max-age=0
- `Upgrade-Insecure-Requests`: 1
- `Accept-Language`: Random English language preference

Supported browsers for `GetRandomUserAgentWithBrowser`:
- `anonymizer.Firefox`
- `anonymizer.Chrome`
- `anonymizer.Opera`
- `anonymizer.Safari`

### ClientV2
ClientV2 is an enhanced HTTP client that uses randomized TCP/TLS fingerprints to make requests appear more like they're coming from real browsers. It leverages the `refraction-networking/utls` library to randomize the TLS ClientHello fingerprint.

```go
import "github.com/ynori7/hulksmash/http"

// Create a new ClientV2
client := http.NewClientV2()

// Use it like a standard http.Client
req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
resp, err := client.Do(req)

// Optionally override the default ClientHello fingerprint
client.SetClientHelloID(utls.HelloChrome_Auto)
```

ClientV2 features:
- Randomized TLS fingerprints to avoid detection
- Automatic HTTP/2 support negotiation
- Connection reuse and automatic reconnection on errors
- Compatible with standard `http.Client` interface

For testing with `httptest`, use `NewClientV2ForTests`:

```go
// In your tests
server := httptest.NewServer(handler)
client := http.NewClientV2ForTests(server.Client().Transport)
```

The ClientV2 is automatically used by the Smasher, but you can also use it standalone for individual requests that need to appear more organic.

## Attribution
Icon from [Sujud.icon](https://www.iconfinder.com/MUHrist) ([CC BY 3.0](https://creativecommons.org/licenses/by/3.0/))

Uses `refraction-networking/utls` to randomize the TCP Hello fingerprint.
