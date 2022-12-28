package main

import (
	//"io/ioutil"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zzuRingo/ratestask/application"
)

func create(x int32) *int32 {
	return &x
}

func equalInt32Pointer(a *int32, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

var testRatesRsp = application.RateRsp{
	{Price: create(1112), Date: "2016-01-01"},
	{Price: create(1112), Date: "2016-01-02"},
	{Price: nil, Date: "2016-01-03"},
	{Price: nil, Date: "2016-01-04"},
}

type rateRsp struct {
	Data application.RateRsp `json:"data"`
}

func TestRates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET",
		"/rates?date_from=2016-01-01&date_to=2016-01-04&origin=CNSGH&destination=north_europe_main", nil)

	// send requst & get rsp
	application.GetRates(c)
	b, _ := ioutil.ReadAll(w.Body)
	if w.Code != 200 {
		t.Error(w.Code, string(b))
		return
	}
	// tricky json format
	bTrimed := strings.TrimLeft(string(b), `"`)
	bTrimed = strings.TrimRight(bTrimed, `"`)
	fmt.Println(bTrimed)
	strRsp := `{\"data\":` + string(bTrimed) + "}"
	x, err := strconv.Unquote(`"` + strRsp + `"`)
	if err != nil {
		t.Errorf("rsp %s format error %v\n", string(bTrimed), err)
		return
	}
	var rsp rateRsp
	e := json.Unmarshal([]byte(x), &rsp)

	// check response
	// a. check if rsp is json
	if e != nil {
		t.Errorf("rsp %s is not json, err is %v\n", string(x), e)
		return
	}
	// b. check if rsp items is enough
	if len(rsp.Data) != len(testRatesRsp) {
		t.Errorf("len rsp = %d, it is not enough. rsp = %s",
			len(rsp.Data), string(b))
		return
	}
	// c. check all item's value
	for i := range testRatesRsp {
		// c1. check date, include 2016-01-04 in which day have no data
		if rsp.Data[i].Date != testRatesRsp[i].Date {
			t.Errorf("ret date %s want %s",
				rsp.Data[i].Date, testRatesRsp[i].Date)
		}
		// c2. check price
		//     -- when the val > 0, is it correct?
		//     -- when the val == nil, on this day, are the records < 3?
		if !equalInt32Pointer(rsp.Data[i].Price, testRatesRsp[i].Price) {
			t.Errorf("at date %s price = %+v want %+v",
				rsp.Data[i].Date, rsp.Data[i].Price, testRatesRsp[i].Price)
		}
	}
}
