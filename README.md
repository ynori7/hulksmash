# Hulk Smash

![HulkSmash Logo](hulksmash.png)

This is a very easy-to-use library for building a custom brute-force requester for QA purposes. This
tool can be useful, for example, for performance and load testing or for testing your rate-limiter. This tool
automatically adds randomized headers to anonymize the request such as `X-Forwarded-For` and `User-Agent`.

This tool is to be used only for benign purposes! 

## Usage
To use it, simply import `"github.com/ynori7/hulksmash"`. Then construct a request builder, which is a 
function to get the request you want to perform. This request builder accepts an index in case you want to 
send requests to a variety of endpoints, with varying payloads, or with a cachebreaker. Then you simply 
create your smasher instance and tell it to start smashing.

```go
requestBuilder := func(index int) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, url, nil)
    return req
}

hulksmash.NewSmasher().Smash(context.Background(), requestBuilder)
```

The smasher comes with some configurable options with safe defaults. Here is a list of the options:

| Option        | Description           | Default  |
| ------------- |:-------------| -----|
| WithClient     | Allows you to override the HTTP Client | &http.Client{} |
| WithWorkerCount      | Sets the number of workers which will send requests in parallel      |   1 |
| WithErrorFunc | Function which is called in case of an error while performing the request      |    Simply logs it to stdout |
| WithSuccessResponseCallback | Function which is called in case of a successful request | Simply logs the http status code and response body |
| WithIterations | The number of calls to make | 1 |
| WithStartIndex | The start index to use when iterating. Can be useful if you want to resume a previous experiment | 0 |

A basic and advanced example can be found in [exaple](example).

## Attribution
Icon from [Sujud.icon](https://www.iconfinder.com/MUHrist) ([CC BY 3.0](https://creativecommons.org/licenses/by/3.0/))
