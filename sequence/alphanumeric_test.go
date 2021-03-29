package sequence

import (
	"fmt"
	"testing"
)

func Test_GetAlphaForKey36(t *testing.T) {
	testdata := map[int]string{
		0:              "0",
		1:              "1",
		100:            "2s",
		10000000000000: "3jlxpt2ps",
	}

	for in, expected := range testdata {
		actual := GetAlphaForKey36(int64(in))
		if expected != actual {
			fmt.Printf("expected %s but got %s\n", expected, actual)
			t.Fail()
		}
	}
}

func Test_GetIndexForAlpha36(t *testing.T) {
	testdata := map[string]int64{
		"0":         0,
		"1":         1,
		"2s":        100,
		"3jlxpt2ps": 10000000000000,
	}

	for in, expected := range testdata {
		actual := GetIndexForAlpha36(in)
		if expected != actual {
			fmt.Printf("expected %d but got %d\n", expected, actual)
			t.Fail()
		}
	}
}
