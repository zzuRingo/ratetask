package application

import (
	"testing"
)

type YyyyMmDdTest struct {
	in  string
	out bool
}

var YyyyMmDdParamTest = []YyyyMmDdTest{
	{"2021-01-02", true},
	{"2021-1-02", false},
	{"201-01-02", false},
	{"2021-01-2", false},
	{"2021/01/02", false},
}

func TestIsYyyyMmDd(t *testing.T) {
	for i := range YyyyMmDdParamTest {
		test := YyyyMmDdParamTest[i]
		out := isYyyyMmDd(test.in)
		if test.out != out {
			t.Errorf("isYyyyMmDd(%s) = %v want %v",
				test.in, out, test.out)
		}
	}
}
