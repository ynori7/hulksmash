package sequence

import (
	"math/big"
)

const InvalidString = -1

// Gets the string value for the given index using lowercase letters and numbers
func GetAlphaForKey36(k int64) string {
	return big.NewInt(k).Text(36)
}

// Gets the index the given string using lowercase letters and numbers. If the string is not valid, returns -1. This is useful for selecting a start index.
func GetIndexForAlpha36(s string) int64 {
	n := new(big.Int)
	n, _ = n.SetString(s, 36)
	if n == nil {
		return InvalidString
	}
	return n.Int64()
}

// AlphaNumeric36 generates a sequence using lowercase letters and numbers
func AlphaNumeric36(min, max int) []string {
	if min >= max {
		return []string{}
	}

	a := make([]string, max-min)
	for i := range a {
		a[i] = GetAlphaForKey36(int64(min + i))
	}
	return a
}
