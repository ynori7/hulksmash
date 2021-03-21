package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ynori7/hulksmash"
)

const (
	startId    = 38656780
	url        = "https://127.0.0.1/onlineCheckIn"
	lastName   = "Doe"
	jsonFormat = `{"flightNumber": "%d", "lastName":"%s"}`
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	requestBuilder := func(index int) *http.Request {
		reqBody := fmt.Sprintf(jsonFormat, index, lastName)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(reqBody)))
		return req
	}

	successFunc := func(resp hulksmash.SuccessResponse) {
		defer resp.Response.Body.Close()
		if resp.Response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Response.Body)
			log.Printf("%d: Found one %s", resp.Response.StatusCode, string(body))
			cancel()
		} else {
			b, _ := resp.Request.GetBody()
			reqBody, _ := ioutil.ReadAll(b)
			log.Printf("%d: Not valid %s", resp.Response.StatusCode, string(reqBody))
		}
	}

	hulksmash.NewSmasher(
		hulksmash.WithStartIndex(startId),
		hulksmash.WithIterations(100),
		hulksmash.WithWorkerCount(3),
		hulksmash.WithSuccessResponseCallback(successFunc),
	).Smash(ctx, requestBuilder)
}
