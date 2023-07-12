//go:build e2e

package accesslog

import (
	"fmt"
	"frame/web"
	"testing"
)

func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, zhangsan")
		ctx.Resp.Write([]byte("hello, zhangsan"))
	})

	server.Start(":8081")
}
