package middlewares

import (
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
)

const requestIdKey = "X-Request-ID"

func AddRequestId(ctx iris.Context) {
	requestId := ctx.GetHeader(requestIdKey)

	if requestId == "" {
		u, err := uuid.NewV4()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		requestId = u.String()

		request := ctx.Request()
		request.Header.Add(requestIdKey, requestId)

		responseWriter := ctx.ResponseWriter()
		responseWriter.Header().Set(requestIdKey, requestId)

		ctx.Values().Set("_x_request_id", requestId)
	}

	ctx.Next()
}