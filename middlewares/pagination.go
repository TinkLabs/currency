package middlewares

import (
	"github.com/kataras/iris"
)

// Pagination parses query string for collection endpoints
// and sets to context
func Pagination(ctx iris.Context) {
	var err error

	// default skip
	skip, err := extractIntParam(ctx, 0, "skip")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	// default limit
	limit, err := extractIntParam(ctx, 10, "limit")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	orderBy := ctx.URLParam("order_by")

	ctx.Values().Set("_skip", skip)
	ctx.Values().Set("_limit", limit)
	ctx.Values().Set("_order_by", orderBy)

	ctx.Next()
}

func StartAfterPagination(ctx iris.Context) {
	startAfter := ctx.URLParam("start_after")

	ctx.Values().Set("_start_after", startAfter)

	ctx.Next()
}

func extractIntParam(ctx iris.Context, defaultVal int, queryStr string) (int, error) {
	var err error
	paramVal := defaultVal
	isParamExists := ctx.URLParamExists(queryStr)
	if isParamExists {
		paramVal, err = ctx.URLParamInt(queryStr)
		if err != nil {
			return paramVal, err
		}
	}
	return paramVal, nil
}
