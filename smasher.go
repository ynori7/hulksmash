package hulksmash

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/ynori7/workerpool"
)

// SuccessResponse is sent to the success callback. This response contains the original request as well as the response.
type SuccessResponse struct {
	Request  *http.Request
	Response *http.Response
}

type (
	// BuildRequestFunc is a function which accepts a iteration index and returns an http request
	BuildRequestFunc        func(index int) *http.Request

	// SuccessResponseCallback is a callback function which is called after a successful http request is performed
	SuccessResponseCallback func(resp SuccessResponse)
)

type smasher struct {
	anonymizer Anonymizer

	client *http.Client

	iterations int
	startIndex int
	workers    int

	onError   func(err error)
	onSuccess SuccessResponseCallback
}

// NewSmasher returns a new smasher with the specified configuration
func NewSmasher(options ...SmasherOption) *smasher {
	s := &smasher{
		anonymizer: NewAnonymizer(),

		// Set defaults
		iterations: defaultIterations,
		startIndex: defaultStartIndex,
		client:     &http.Client{},
		workers:    defaultWorkerCount,
		onError:    defaultOnError,
		onSuccess:  defaultSuccessResponseCallback,
	}

	// apply options
	for _, opt := range options {
		opt(s)
	}

	return s
}

// Smash will perform the configured requests repeatedly based on the configuration
func (s *smasher) Smash(ctx context.Context, buildRequest BuildRequestFunc) {
	workerPool := workerpool.NewWorkerPool(
		func(result interface{}) { //will be a SuccessResponse
			resp := result.(SuccessResponse)
			s.onSuccess(resp)
		},
		s.onError,
		func(job interface{}) (result interface{}, err error) {
			index := job.(int)

			req := buildRequest(index)
			if req == nil {
				return nil, errors.New("invalid request")
			}
			s.anonymizer.AnonymizeRequest(req) //disguise the traffic

			resp, err := s.client.Do(req)
			if err != nil {
				return nil, err
			}

			return SuccessResponse{Request: req, Response: resp}, nil
		})

	list := makeRange(s.startIndex, s.startIndex+s.iterations)

	if err := workerPool.Work(
		ctx,
		s.workers, //The number of workers which should work in parallel
		list,      //The items to be processed
	); err != nil {
		log.Println(err.Error())
	}
}

// makeRange returns a list starting from min up to (but not including) max
func makeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}
