package main

import (
	"context"
	"net/http"

	"github.com/ynori7/hulksmash"
)

const (
	url = "http://127.0.0.1"
)

func main() {
	requestBuilder := func(index int) *http.Request {
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		return req
	}

	hulksmash.NewSmasher().Smash(context.Background(), requestBuilder)
}
