package application

import (
	"testing"
)

type TestIsStubCase struct {
	in  string
	out bool
}

var isStubTestCases = []TestIsStubCase{
	{"CNSGH", false},
	{"CNSGh", true},
	{"norway_south_east", true},
	{"Åhus", true},
}

func TestIsStub(t *testing.T) {
	for i := range isStubTestCases {
		test := isStubTestCases[i]
		out := isStub(test.in)
		if test.out != out {
			t.Errorf("isStub(%s) = %v want %v",
				test.in, out, test.out)
		}
	}
}
