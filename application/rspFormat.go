package application

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

type RateRsp []dailyPrice

type dailyPrice struct {
	Price *int32 `json:"price"`
	Date  string `json:"date"`
}

func trimDateYmd(date string) (int, int, int) {
	year, _ := strconv.Atoi(date[0:4])
	month, _ := strconv.Atoi(date[5:7])
	day, _ := strconv.Atoi(date[8:10])
	return year, month, day
}

func strToDate(str string) time.Time {
	y, m, d := trimDateYmd(str)
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func buildPriceItem(date string, price sql.NullInt32) dailyPrice {
	var priceVal *int32 = nil
	if price.Valid {
		priceVal = &price.Int32
	}

	return dailyPrice{
		Price: priceVal,
		Date:  date,
	}
}

func composeRateRsp(cons []priceCons, startDate string) string {
	if len(cons) == 0 {
		return "[]"
	}

	curDate := strToDate(startDate)
	rsp := RateRsp{}
	for _, dailyPrice := range cons {
		thisRoundDate := strToDate(dailyPrice.Date)
		// if last written date is before than this round's date
		// then fill these days' price with null
		for curDate.Before(thisRoundDate) {
			strCurDate := curDate.Format("2006-01-02")
			rsp = append(rsp, buildPriceItem(strCurDate, sql.NullInt32{}))
			curDate = curDate.AddDate(0, 0, 1)
		}
		// append this round's data
		if thisRoundDate == curDate {
			curDate = curDate.AddDate(0, 0, 1)
			rsp = append(rsp, buildPriceItem(dailyPrice.Date, dailyPrice.Price))
		}
	}

	rspJson, _ := json.Marshal(rsp)
	return string(rspJson)
}
