package hulksmash

import (
	"context"
	"io/ioutil"
	"log"
	nethttp "net/http"

	"github.com/ynori7/hulksmash/anonymizer"
	"github.com/ynori7/hulksmash/http"
	"github.com/ynori7/workerpool"
)

// SuccessResponse is sent to the success callback. This response contains the original request as well as the response.
type SuccessResponse struct {
	StatusCode   int
	RequestBody  []byte
	ResponseBody []byte

	RawRequest  *nethttp.Request
	RawResponse *nethttp.Response
}

type (
	// BuildRequestFunc is a function which accepts a iteration index and returns an http request
	BuildRequestFunc func(index int) (*nethttp.Request, error)

	// SuccessResponseCallback is a callback function which is called after a successful http request is performed
	SuccessResponseCallback func(resp SuccessResponse)
)

type smasher struct {
	anonymizeRequets bool
	anonymizer       anonymizer.Anonymizer

	client *nethttp.Client

	iterations int
	startIndex int
	workers    int

	onError   func(err error)
	onSuccess SuccessResponseCallback
}

// NewSmasher returns a new smasher with the specified configuration
func NewSmasher(options ...SmasherOption) *smasher {
	s := &smasher{
		anonymizer: anonymizer.New(),

		// Set defaults
		iterations:       defaultIterations,
		startIndex:       defaultStartIndex,
		client:           http.NewClient(),
		workers:          defaultWorkerCount,
		onError:          defaultOnError,
		onSuccess:        defaultSuccessResponseCallback,
		anonymizeRequets: defaultAnonymizeRequests,
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

			req, err := buildRequest(index)
			if err != nil {
				return nil, err
			}

			if s.anonymizeRequets {
				s.anonymizer.AnonymizeRequest(req) //disguise the traffic
			}

			resp, err := s.client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			successResp := SuccessResponse{
				StatusCode:  resp.StatusCode,
				RawRequest:  req,
				RawResponse: resp,
			}

			if req.GetBody != nil {
				b, _ := req.GetBody()
				successResp.RequestBody, _ = ioutil.ReadAll(b)
			}
			successResp.ResponseBody, _ = ioutil.ReadAll(resp.Body)

			return successResp, nil
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
