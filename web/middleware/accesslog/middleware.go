package accesslog

import (
	"encoding/json"
	"my-frame/web"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 记录请求
			defer func() {
				l := accesslog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))

			}()
			next(ctx)
		}
	}
}

type accesslog struct {
	Host string `json:"host,omitempty"`
	// 命中的路由
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}

/*
Middleware 例子- AccessLog 潜在问题
这种实现的问题
	太死板了，只记录寥寥几个字段，如果用户要记录更多的数据怎么办?
	固定采用了json 作为序列化的方式，要是用户想要用别的序列化方式，例如 protobuf，怎么办?
	不怎么办，让用户自己写一个 AccessLog !

这是我和别人设计理念上一个很大的差异。在这种用户完全可以自己提供扩展实现的地方，我是不愿意花费很多时间去设计精巧的结构和机制，去支持各种千奇百怪的用法的。我坚持的只有一个: 你有特殊需求，你就自己写。
默认提供的实现是给大多数的普通用户使用的，也相当于一个例子，有需要的用户可以参考这个实现写自己的实现。
*/

/*
思考-要不要考虑时机问题?
	例如，要不要设计类似 Beego 那种复杂的在不同阶段运行的Filter?
	这个问题和我们之前讨论的其它问题都不太一样，是因为用户完全没办法自己支持，只能依赖于框架支持 (也就是必须侵入式地修改框架)
	理论上来说是需要考虑的，但是我们可以推迟到用户真正需要的时候再来评估。
	因为大多数场景都是不需要考虑的，已有的设计完全能够满足


思考-Middleware 要不要考虑顺序问题

	理论上来说，每一个 Middleware 都应该不依赖于其它的 Middleware
	但这只是一个美好的希望，比如在我们已经实现的几个 Middleware 里面，Panic 很显然应该在最外层，也就是紧接着 flashResp 那里，错误处理应该在可观测性之后
	又比如，从业务上来看，鉴权应该在很靠前的位置，限流可以在鉴权前面，也可以在鉴权后面，取决于业务......

思考一中间件要不要考虑分路由问题
	前面所有的中间件都是对所有请求生效的
	个很常见的场景
	我们希望区分不同的路由，进行不同的处理。例如公开页面，用户不需要登录，但是有一些页面，用户就需要登录。
	这是你们的作业



面试要点
	什么是可观测性?
		也就是 logging、metrics 和 tracing。

	常用的可观测性的框架有哪些?
		你举例自己公司用的，开源的 OpenTelemetry、SkyWalking.Prometheus.

	怎么集成可观测性框架?
		一般都是利用 Middleware 机制，不仅仅是 Web 框架，几乎所有的框架都有类似 Middleware 的机制。

	Prometheus 的 Histogam和 Summary?
	全链路追踪 (tracing) 的几个概念?
	解释一下 tracer、tracing 和 span 的概念tracing 是怎么构建的?
		核心在于解释清楚 tracing 进程内和跨进程的运作，我们将在微服务框架里面看到跨进程是怎么处理的。

	HTTP 应该观测一些什么数据?
		也就是我们OpenTelemetrv 和 Promtheus 两个 Middleware 里面写的那些指标。

	什么是 99 线、999线?
		就是响应的比例，99% 的响应、99.9 %的响应的响应时间在多少以内。
*/
