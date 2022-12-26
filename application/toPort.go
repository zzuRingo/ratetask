package application

import (
	"fmt"
	"unicode"
)

func isStub(s string) bool {
	return !onlyCapitals(s)
}

// return param corresponding ports
func paramToPortCode(param string) ([]string, error) {
	// if param is not stub, then assume it as port code
	if !isStub(param) {
		return []string{param}, nil
	}

	// find stub's corresponding ports in db
	portCodes, err := GetStubCorrespondingPort(param)
	if err != nil {
		fmt.Println("GetStubCorrespondingPort error :%+v", err)
	}
	return portCodes, nil
}

// if s only contains capital alphabets, return true
func onlyCapitals(s string) bool {
	for _, ch := range s {
		if !unicode.IsUpper(ch) {
			return false
		}
	}
	return true
}