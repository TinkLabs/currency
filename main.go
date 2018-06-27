package main

import (
	"currency/routers"

	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	routers.Register(app)
	app.Run(iris.Addr(":8080"), iris.WithoutStartupLog, iris.WithoutVersionChecker)
}
