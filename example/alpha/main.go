package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ynori7/hulksmash/sequence"
	"log"
	"net/http"

	"github.com/ynori7/hulksmash"
)

const (
	url        = "https://web-api-stage.goorange.sixt.com/v1/rental-testing/dummy"
	lastName   = "Doe"
	jsonFormat = `{"flightNumber": "%s", "lastName":"%s"}`
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	requestBuilder := func(item string) (*http.Request, error) {
		reqBody := fmt.Sprintf(jsonFormat, item, lastName)
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
		hulksmash.WithStartIndex(int(sequence.GetIndexForAlpha36("aaaaaa"))),
		hulksmash.WithIterations(100),
		hulksmash.WithWorkerCount(3),
		hulksmash.WithSuccessResponseCallback(successFunc),
		hulksmash.WithSequenceFunc(sequence.AlphaNumeric36),
	).Smash(ctx, requestBuilder)
}
