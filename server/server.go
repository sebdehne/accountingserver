package server

import (
	"github.com/kataras/iris"
	"strconv"
)

type Api struct {
	Version int
	Routes  []Route
}

type Route struct {
	Method      string
	PathPattern string
	Handler     iris.HandlerFunc
}

func RunServer(prefixPath string, apis ...Api) {
	i := iris.New()

	for _, api := range apis {
		for _, r := range api.Routes {
			i.HandleFunc(r.Method, prefixPath + "/v" + strconv.Itoa(api.Version) + r.PathPattern, r.Handler)
		}
	}

	i.Any("/", func(ctx *iris.Context) {
		ctx.Error("Not Found", iris.StatusNotFound)
	})

	i.Listen(":8081")
}

