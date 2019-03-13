package routers

import (
	"time"

	_ "currency/core/logger"
	hd "currency/handlers"
	mid "currency/middlewares"

	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
)

func Register(app *iris.Application) {
	app.Use(mid.AccessLogger)
	app.Use(mid.AddRequestId)

	opts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type"},
		AllowedMethods: []string{"GET", "POST", "HEAD"},
		MaxAge:         int((24 * time.Hour).Seconds()),
		// Debug:          true,
	}
	app.Use(cors.New(opts))
	app.AllowMethods(iris.MethodOptions)

	app.Get("/", hd.Welcome)
	currency(app)
}