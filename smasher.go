package hulksmash

import (
	"context"
	"io/ioutil"
	"log"
	nethttp "net/http"
	"time"

	"github.com/ynori7/hulksmash/anonymizer"
	"github.com/ynori7/hulksmash/http"
	"github.com/ynori7/hulksmash/sequence"
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
	// BuildRequestFunc is a function which accepts a iteration item and returns an http request
	BuildRequestFunc func(item string) (*nethttp.Request, error)

	// SuccessResponseCallback is a callback function which is called after a successful http request is performed
	SuccessResponseCallback func(resp SuccessResponse)
)

type smasher struct {
	anonymizeRequets bool
	anonymizer       anonymizer.Anonymizer

	client *nethttp.Client

	iterations   int
	startIndex   int
	sequenceFunc sequence.SequenceFunc
	workers      int

	onError   func(err error)
	onSuccess SuccessResponseCallback
}

// NewSmasher returns a new smasher with the specified configuration
func NewSmasher(options ...SmasherOption) *smasher {
	s := &smasher{
		anonymizer: anonymizer.New(time.Now().UnixNano()),

		// Set defaults
		iterations:       defaultIterations,
		startIndex:       defaultStartIndex,
		client:           http.NewClient(),
		workers:          defaultWorkerCount,
		onError:          defaultOnError,
		onSuccess:        defaultSuccessResponseCallback,
		anonymizeRequets: defaultAnonymizeRequests,
		sequenceFunc:     sequence.Numeric,
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
			item := job.(string)

			req, err := buildRequest(item)
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

	list := s.sequenceFunc(s.startIndex, s.startIndex+s.iterations)

	if err := workerPool.Work(
		ctx,
		s.workers, //The number of workers which should work in parallel
		list,      //The items to be processed
	); err != nil {
		log.Println(err.Error())
	}
}
