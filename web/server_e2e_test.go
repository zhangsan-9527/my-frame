//go:build e2e

package web

import (
	"fmt"
	"strconv"
	"testing"
)

func TestServer(t *testing.T) {
	// 这个 Handler 就是我们跟 http包 的结合点
	//http.ListenAndServe(":8080", "要在这里搞一个 handler")

	//var s Server // NewServer方法
	//var h Server = &HTTPServer{}
	h := NewHTTPServer()

	//handle1 := func(ctx *Context) {
	//	fmt.Println("处理第一件事")
	//}
	//
	//handle2 := func(ctx *Context) {
	//	fmt.Println("处理第二件事")
	//}

	//// 用户自己去管这种
	//h.addRoute(http.MethodGet, "/user/home", func(ctx *Context) {
	//	handle1(ctx)
	//	handle2(ctx)
	//
	//})

	//h.addRoute(http.MethodGet, "/order/detail", func(ctx *Context) {
	//	ctx.Resp.Write([]byte("Hello order detail"))
	//
	//})

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello order detail"))

	})

	h.Get("/order/abc", func(ctx *Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))

	})

	h.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %s", ctx.Req.URL.Path)))

	})

	h.Post("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", id)))

	})
	h.Get("/values/:id", func(ctx *Context) {
		idStr, err := ctx.PathValue("id")
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("id 输入不对"))
			return
		}

		ctx.Resp.Write([]byte(fmt.Sprintf("hello, %d", id)))

	})

	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	h.Get("/user/123", func(ctx *Context) {
		ctx.RespJSON(202, User{
			Name: "zhangsan",
			Age:  18,
		})
	})

	// 完全委托给 http 包
	//h.AddRoute1(http.MethodGet, "/user", handle1, handle2)
	// 没有意义
	//h.AddRoute1(http.MethodGet, "/user")

	//// 用法一 完全委托给 http 包
	//http.ListenAndServe(":8899", h)
	//http.ListenAndServeTLS(":443", "", "", h)

	// 用法二 自己手动管
	h.Start(":8081")

}

//type SafeContext struct {
//	Context
//	mutex sync.RWMutex
//}
//
//func (c *SafeContext) RespJSONOK(val any) error {
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//	return c.Context.RespJSON(http.StatusOK, val)
//}
