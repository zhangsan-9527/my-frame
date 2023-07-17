package web

import (
	"fmt"
	"net"
	"net/http"
)

// Server 对于特性上来说:
// 至少要提供三部分功能:
//    生命周期控制:即启动、关闭。如果在后期，我们还要考虑增加生命周期回调特性
//    路由注册接口:提供路由注册功能
//    作为 http包 到 Web框架的桥梁

type HandleFunc func(ctx *Context)

// 接口实现校验
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// addRoute 增加路由注册功能
	// method 是 HTTP 方法
	// path 是 路由
	// handleFunc 是 业务逻辑
	addRoute(method string, path string, handleFunc HandleFunc)

	// AddRoute1 这种提供多个(没必要)
	//AddRoute1(method string, path string, handle ...HandleFunc)
}

//type HTTPSServer struct {
//	HTTPServer
//}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	// addr string 创建的时候传递, 而不是Start接收  这个都是可以的
	//router
	router
	//r *router
	mdls []Middleware

	log func(msg string, args ...any)

	tplEngine TemplateEngine

	//Middleware []Middleware

}

//func NewHTTPServerV1(mdls ...Middleware) *HTTPServer {
//	return &HTTPServer{
//		router: newRouter(),
//		mdls:   mdls,
//	}
//}

// 第一个问题: 相对路径还是绝对路径?
// 第二个问题: json, yaml, xml
//func NewHTTPServerV2(cfgFilePath string) *HTTPServer {
//	// 你在这里加载配置, 解析, 然后初始化 HTTPServer
//}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithTemplateEngine(tpl TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tpl
	}
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

// ServeHTTP 处理请求的入口
// ServeHTTP 则是我们整个 Web 框架的核心入口。我们将在整个方法内部完成
// Context 构建
// 路由匹配
// 执行业务逻辑
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 你的框架代码就在这里
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: h.tplEngine,
	}

	// 接下来就是查找路由，并且执行命中的业务逻辑

	// 最后一个是这个
	root := h.serve

	// 然后这里就是利用最后一个不断往前回溯组装链条
	// 从后往前
	// 把后一个作为前一个的 next 构造好链条
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}

	// 这里执行的时候, 就是从前往后了

	// 这里, 最后一个步骤, 就是把 RespData 和 RespStatusCode 刷新到响应里面

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			// 就设置好了 RespData 和 RespStatusCode
			next(ctx)
			h.flashResp(ctx)
		}
	}

	root = m(root)

	root(ctx)
	//h.serve(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}

	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		h.log("写入响应数据失败 %v", err)
	}
}

func (h *HTTPServer) serve(ctx *Context) {
	// before route
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	// after route
	if !ok || info.n.handler == nil {
		// 路由没有命中, 就是 404
		ctx.RespStatusCode = 404
		ctx.RespData = []byte("NOT FOUND")
		return
	}
	ctx.PathParams = info.pathParams
	ctx.MatchedRoute = info.n.route
	// before execute
	info.n.handler(ctx)
	// after execute

}

func (h *HTTPServer) Start(addr string) error {
	// 也可以自己创建Server
	//http.Server{}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// 在这里, 可以让用户注册所谓的 after start 回调
	// 比如: 往你的 admin 注册一下自己这个实例
	// 在这里执行一些你业务所需的前置条件

	return http.Serve(l, h)
}

//func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
//	// 注册到路由树里面
//	//panic("implement me")
//}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}
func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}
func (h *HTTPServer) Put(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}
func (h *HTTPServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

//func (h *HTTPServer) AddRoute1(method string, path string, handle ...HandleFunc) {
//	panic("implement me")
//}

func (h *HTTPServer) Start1(addr string) error {

	return http.ListenAndServe(addr, h)
}

/*

面试要点
HTTP 服务器的生命周期? 一般来说就是启动、运行和关闭。在这三个阶段的前后都可以插入生命周期回调。一般来说，面试生命周期，多半都是为了问生命周期回调。例如说怎么做Web 服务的服务发现? 就是利用生命周期回调的启动后回调，将 Web 服务注册到服务中心

HTTP Server 功能? 记住在不同的框架里面有不同的叫法，比如说在 Gin 里面叫做Engine，它们的基本功能都是提供路由注册，生命周期控制以及作为与 http 包结合的杯梁。

*/
