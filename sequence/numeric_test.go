package sequence

import (
	"fmt"
	"testing"
)

func Test_Numeric(t *testing.T) {
	expected := []string{"0", "1", "2", "3"}
	actual := Numeric(0, 4)

	if len(expected) != len(actual) {
		fmt.Println("unexpected length")
		t.FailNow()
	}
	for k := range actual {
		if expected[k] != actual[k] {
			fmt.Printf("expected %s but got %s\n", expected, actual)
			t.Fail()
		}
	}
}

func Test_Numeric_WithMinValue(t *testing.T) {
	expected := []string{"5", "6", "7", "8"}
	actual := Numeric(5, 9)

	if len(expected) != len(actual) {
		fmt.Println("unexpected length")
		t.FailNow()
	}
	for k := range actual {
		if expected[k] != actual[k] {
			fmt.Printf("expected %s but got %s\n", expected, actual)
			t.Fail()
		}
	}
}

func Test_Numeric_MinGreaterThanMax(t *testing.T) {
	actual := Numeric(7, 4)

	if 0 != len(actual) {
		fmt.Println("unexpected length")
		t.FailNow()
	}

}
