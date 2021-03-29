package anonymizer

import (
	"testing"
)

func Test_GetRandomIP(t *testing.T) {
	a := New(1)

	expected1 := "1273:a2d8:3cc3:bbdc:2971:d3a0:e6b5:8bb2"
	expected2 := "120.152.43.60"

	actual1 := a.GetRandomIp()
	actual2 := a.GetRandomIp()

	if actual1 != expected1 || actual2 != expected2 {
		t.Fail()
	}
}

func Test_GetRandomUserAgent(t *testing.T) {
	a := New(1)

	expected1 := "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y) AppleWebKit/447.1.19 (KHTML, like Gecko) Version/13.1.1 Mobile/15E148 Safari/431.1"
	expected2 := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.1456.100 Safari/544.11"

	actual1 := a.GetRandomUserAgent()
	actual2 := a.GetRandomUserAgent()

	if actual1 != expected1 || actual2 != expected2 {
		t.Fail()
	}
}

func Test_GetRandomUserAgent_UserAgentFuncs(t *testing.T) {
	a := New(1)

	foundAgents := map[string]struct{}{}

	for _, f := range userAgentFuncs {
		foundAgents[f(a, "test")] = struct{}{}
	}

	// assert that each user agent was unique
	if len(foundAgents) != len(userAgentFuncs) {
		t.Fail()
	}
}