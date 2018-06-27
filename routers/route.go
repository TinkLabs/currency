package routers

import (
	_ "currency/core/logger"
	hd "currency/handlers"
	mid "currency/middlewares"

	"github.com/kataras/iris"
)

func Register(app *iris.Application) {
	app.Use(mid.AccessLogger)
	app.Use(mid.AddRequestId)

	app.Get("/", hd.Welcome)
	currency(app)
}