package sequence

import "fmt"

// Numeric returns a slice of strings containing integers starting from min up to (but not including) max
func Numeric(min, max int) []string {
	if min >= max {
		return []string{}
	}

	a := make([]string, max-min)
	for i := range a {
		a[i] = fmt.Sprintf("%d", min+i)
	}
	return a
}
