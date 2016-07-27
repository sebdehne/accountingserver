package server

import (
	"github.com/kataras/iris"
	"strconv"
)

type Api struct {
	Prefix  string
	Version int
	Routes  []Route
}

type Route struct {
	Method      string
	PathPattern string
	Handler     iris.HandlerFunc
}

func RunServer(listenExpr string, apis ...Api) {
	i := iris.New()

	for _, api := range apis {
		for _, r := range api.Routes {
			i.HandleFunc(r.Method, api.Prefix + "/v" + strconv.Itoa(api.Version) + r.PathPattern, r.Handler)
		}
	}

	i.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
		ctx.Error("Page " + ctx.PathString() + " not Found", iris.StatusNotFound)
	})

	i.Listen(listenExpr)
}

