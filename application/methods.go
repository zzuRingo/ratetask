package application

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func getRates(c *gin.Context) {
	// 1. check param
	param := ProcessRatesParam(c.Request.URL.Query())
	if param == nil {
		c.JSON(http.StatusBadRequest, "error param")
		return
	}

	// 2. transform input slug to port code
	var wg sync.WaitGroup
	var errDest, errOri error = nil, nil
	destPortCodes, oriPortCodes := []string{}, []string{}
	// do the two query parallel
	wg.Add(2)
	go func() {
		destPortCodes, errDest = paramToPortCode(param.destination)
		wg.Done()
	}()
	go func() {
		oriPortCodes, errOri = paramToPortCode(param.origin)
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("code1 %+v\n code2 %+v\n", oriPortCodes, destPortCodes)
	if errOri != nil || errDest != nil {
		c.JSON(http.StatusInternalServerError, "error when query corresponding port codes")
	}

	// 3. fetch average price
	avgPrices, err := GetDailyAvgPrice(oriPortCodes, destPortCodes,
		param.dateFrom, param.dateTo)
	if err != nil {
		fmt.Printf("GetDailyAvgPrice error :%+v", err)
		c.JSON(http.StatusInternalServerError, "error when query avg prices")
	}

	// 4. compose rsp
	c.JSON(http.StatusOK, composeRateRsp(avgPrices, param.dateFrom))
}
