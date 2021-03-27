package anonymizer

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Anonymizer is an object used for getting randomized data to help make a request look less suspicious
type Anonymizer struct {
	rand *rand.Rand
}

// New returns a new anonymizer
func New() Anonymizer {
	return Anonymizer{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AnonymizeRequest accepts an http request and anonymizes it by adding a forwarded-for header and user agent string
func (a Anonymizer) AnonymizeRequest(r *http.Request) {
	r.Header.Set("X-Forwarded-For", a.GetRandomIp())
	r.Header.Set("User-Agent", a.GetRandomUserAgent())
	r.Header.Set("Cache-Control", "max-age=0")
	r.Header.Set("Upgrade-Insecure-Requests", "1")
}

// GetRandomIp returns a random IPv4 or IPv6 address
func (a Anonymizer) GetRandomIp() string {
	if a.rand.Intn(2) == 1 {
		return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x",
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
			a.rand.Intn(65535)+1,
		)
	}

	return fmt.Sprintf("%d.%d.%d.%d",
		a.rand.Intn(255)+1,
		a.rand.Intn(255)+1,
		a.rand.Intn(255)+1,
		a.rand.Intn(255)+1,
	)
}

// GetRandomUserAgent returns a randomized user agent string
func (a Anonymizer) GetRandomUserAgent() string {
	platform := platforms[a.rand.Intn(len(platforms))]
	userAgentFunc := userAgentFuncs[a.rand.Intn(len(userAgentFuncs))]
	return userAgentFunc(a, platform)
}

var platforms = []string{
	"Windows NT 6.1; Win64; x64",                 //windows
	"Macintosh; Intel Mac OS X x.y",              //mac
	"iPhone; CPU iPhone OS 13_5_1 like Mac OS X", //iOS
	"X11; Linux x86_64",                          //linux
	"Windows Phone OS 7.5",                       //windows phone
}

var userAgentFuncs = []func(a Anonymizer, platform string) string{
	// Firefox
	func(a Anonymizer, platform string) string {
		//Mozilla/5.0 (platform; rv:[40-50].0) Gecko/20[10-21][01-12][01-31] Firefox/[40-50].0
		format := "Mozilla/5.0 (%s; rv:%d.0) Gecko/%s Firefox/%d.%d"
		return fmt.Sprintf(format, platform, a.rand.Intn(10)+40, a.randomDate(), a.rand.Intn(10)+40, a.rand.Intn(100))
	},

	// Chrome
	func(a Anonymizer, platform string) string {
		//Mozilla/5.0 (platform) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/[40-55].0.[1000-3000].[0-200] Safari/[500-550].[0-100]
		format := "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/%d.%d"
		return fmt.Sprintf(format, platform, a.rand.Intn(15)+40, a.rand.Intn(2000)+1000, a.rand.Intn(200), a.rand.Intn(50)+500, a.rand.Intn(100))
	},

	// Opera
	func(a Anonymizer, platform string) string {
		//Mozilla/5.0 (platform) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/[40-55].0.[1000-3000].[0-200] Safari/[500-550].[0-100] OPR/[30-40].0.[1000-3000].[0-100]
		format := "Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.%d.%d Safari/%d.%d OPR/%d.0.%d.%d"
		return fmt.Sprintf(format, platform, a.rand.Intn(15)+40, a.rand.Intn(2000)+1000, a.rand.Intn(200),
			a.rand.Intn(50)+500, a.rand.Intn(100), a.rand.Intn(10)+30, a.rand.Intn(2000)+1000, a.rand.Intn(100))
	},

	// Safari
	func(a Anonymizer, platform string) string {
		//Mozilla/5.0 (platform) AppleWebKit/[550-650].1.[0-20] (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/[550-650].1
		format := "Mozilla/5.0 (%s) AppleWebKit/%d.1.%d (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/%d.1"
		return fmt.Sprintf(format, platform, a.rand.Intn(150)+400, a.rand.Intn(20), a.rand.Intn(150)+400)
	},
}

const dateFormat = "20060102"

func (a Anonymizer) randomDate() string {
	min := time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := a.rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format(dateFormat)
}
