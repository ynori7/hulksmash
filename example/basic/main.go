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
	requestBuilder := func(index int) (*http.Request, error) {
		return http.NewRequest(http.MethodGet, url, nil)
	}

	hulksmash.NewSmasher().Smash(context.Background(), requestBuilder)
}
