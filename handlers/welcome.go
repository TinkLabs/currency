package handlers

import (
	"github.com/kataras/iris"
)

func Welcome(ctx iris.Context) {
	ctx.Text("Welcome to currency!!")
}
