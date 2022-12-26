package application

import (
	"fmt"
	"net/url"
	"strconv"
)

type ratesParam struct {
	origin 			string
	destination 	string
	dateFrom 		string
	dateTo 			string
}

// bit val when param is complete
const fullParamBit = 0b1111

// ratesParam type
const (
	RatesParamOrigin int = iota
	RatesParamDestination
	RatesParamDateFrom
	RatesParamDateTo
)

// add paramBit's corresponding bit
func FlagParamBit(paramBit *int, ratesParamType int) {
	*paramBit = *paramBit | (1 << ratesParamType)
}

// if str is yyyy-mm-dd format, return true
// otherwise return false
func isYyyyMmDd(str string) bool {
	if len(str) != 10 {
		return false
	}

	if str[4] != '-' || str[7] != '-' {
		return false
	}

	year := str[0:4]
	month := str[5:7]
	day := str[8:10]
	if _, err := strconv.Atoi(year); err != nil {
		return false
	}
	if _, err := strconv.Atoi(month); err != nil {
		return false
	}
	if _, err := strconv.Atoi(day); err != nil {
		return false
	}
	return true
}

func ProcessRatesParam(param url.Values) *ratesParam {
	fmt.Println("full : ", fullParamBit)
	paramBit := 0
	rParam := &ratesParam{}
	for k, v := range param {
		if len(v) != 1 {
			fmt.Printf("k:%s, v:%+v", k, v)
			return nil
		}
		fmt.Println("param k:%s, v:%+v", k, v)
		if k == "origin" {
			FlagParamBit(&paramBit, RatesParamOrigin)
			rParam.origin = v[0]
			continue
		}
		if k == "destination" {
			FlagParamBit(&paramBit, RatesParamDestination)
			rParam.destination = v[0]
			continue
		}
		if k == "date_from" {
			if !isYyyyMmDd(v[0]) {
				fmt.Println("is not valid date")
				return nil
			}
			FlagParamBit(&paramBit, RatesParamDateFrom)
			rParam.dateFrom = v[0]
			continue
		}
		if k == "date_to" {
			if !isYyyyMmDd(v[0]) {
				fmt.Println("is not valid date")
				return nil
			}
			FlagParamBit(&paramBit, RatesParamDateTo)
			rParam.dateTo = v[0]
			continue
		}
	}
	if paramBit != fullParamBit {
		fmt.Println("param missed: cur bit%b", paramBit)
		return nil
	}
	if rParam.dateFrom > rParam.dateTo {
		fmt.Println("rParam.dateFrom > rParam.dateTo")
		return nil
	}
	return rParam
}