package middlewares

import (
	"time"

	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
)

func AccessLogger(ctx iris.Context) {
	start := time.Now().UTC()

	ctx.Next()

	requestId := ctx.Values().GetString("_x_request_id")
	method := ctx.Method()
	path := ctx.Path()
	address := ctx.RemoteAddr()
	statusCode := ctx.GetStatusCode()
	timeElapsed := float64(time.Since(start)) / float64(time.Millisecond)

	if path == "/" {
		return
	}

	log := logrus.WithFields(logrus.Fields{"method": method, "path": path, "status_code": statusCode, "time_elapsed": timeElapsed, "request_id": requestId, "address": address})
	log.Info("HTTP request")
}
