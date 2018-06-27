package middlewares

import (
	"time"

	"github.com/kataras/iris"
)

const DateQueryStringFormat = "2006-01-02"

func PraseDates(ctx iris.Context) {
	dateStr := ctx.FormValue("date")
	dateStartStr := ctx.FormValue("date_start")
	dateEndStr := ctx.FormValue("date_end")

	if dateStr != "" {
		date, err := time.Parse(DateQueryStringFormat, dateStr)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		ctx.Values().Set("_date", &date)
	}

	if dateStartStr != "" {
		dateStart, err := time.Parse(DateQueryStringFormat, dateStartStr)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		ctx.Values().Set("_date_start", &dateStart)
	}

	if dateEndStr != "" {
		dateEnd, err := time.Parse(DateQueryStringFormat, dateEndStr)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}
		ctx.Values().Set("_date_end", &dateEnd)
	}

	ctx.Next()
}
