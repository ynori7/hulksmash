package main

import (
	"bytes"
	"context"
	"fmt"
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

	requestBuilder := func(index int) (*http.Request, error) {
		reqBody := fmt.Sprintf(jsonFormat, index, lastName)
		return http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(reqBody)))
	}

	successFunc := func(resp hulksmash.SuccessResponse) {
		if resp.StatusCode == 200 {
			log.Printf("%d: Found one %s", resp.StatusCode, string(resp.ResponseBody))
			cancel()
		} else {
			log.Printf("%d: Not valid %s", resp.StatusCode, string(resp.RequestBody))
		}
	}

	hulksmash.NewSmasher(
		hulksmash.WithStartIndex(startId),
		hulksmash.WithIterations(100),
		hulksmash.WithWorkerCount(3),
		hulksmash.WithSuccessResponseCallback(successFunc),
	).Smash(ctx, requestBuilder)
}
