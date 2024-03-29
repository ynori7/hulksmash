package hulksmash

import (
	"log"

	"github.com/ynori7/hulksmash/sequence"
)

// SmasherOption is a functional option for overriding the default configuration
type SmasherOption func(s *smasher)

// WithClient allows you to override the default http client, for example if you want to add your own round-tripper or timeouts
func WithClient(c HttpClient) SmasherOption {
	return func(s *smasher) {
		s.client = c
	}
}

// WithWorkerCount overrides the default worker count (default is 1, so a single thread)
func WithWorkerCount(c int) SmasherOption {
	return func(s *smasher) {
		s.workers = c
	}
}

// WithErrorFunc overrides the default error callback, which simply logs the error to stdout. This function is called in case
// of an error while performing the request (not in case of an error response)
func WithErrorFunc(f func(err error)) SmasherOption {
	return func(s *smasher) {
		s.onError = f
	}
}

// WithSuccessResponseCallback overrides the default success callback which handles all http responses. The default simply logs the response body
func WithSuccessResponseCallback(f SuccessResponseCallback) SmasherOption {
	return func(s *smasher) {
		s.onSuccess = f
	}
}

// WithIterations sets the number of calls the smasher should make (default is 1)
func WithIterations(i int) SmasherOption {
	return func(s *smasher) {
		s.iterations = i
	}
}

// WithStartIndex overrides the default start index for each iteration in case, for example, you want to resume from a previous position
func WithStartIndex(i int) SmasherOption {
	return func(s *smasher) {
		s.startIndex = i
	}
}

// WithAnonymizeRequests overrides the default flag to indicate whether reqeusts should be anonymized (by adding additional headers)
func WithAnonymizeRequests(anonymize bool) SmasherOption {
	return func(s *smasher) {
		s.anonymizeRequets = anonymize
	}
}

func WithSequenceFunc(f sequence.SequenceFunc) SmasherOption {
	return func(s *smasher) {
		s.sequenceFunc = f
	}
}

var (
	defaultIterations        = 1
	defaultStartIndex        = 0
	defaultWorkerCount       = 1
	defaultAnonymizeRequests = true
	defaultOnError           = func(err error) {
		log.Println(err.Error())
	}
	defaultSuccessResponseCallback = func(resp SuccessResponse) {
		log.Printf("Status %d, Body %s\n", resp.StatusCode, string(resp.ResponseBody))
	}
)
