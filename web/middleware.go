package web

import "sync"

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

// AOP 方案在不同的框架, 不同的语言里都有不同的叫法
//type MiddlewareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//
//// Interceptor 拦截器模式
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}

//type HandleFuncV1 func(ctx *Context) (next bool)
//
//// Chain 集中式
//type Chain []HandleFuncV1
//
//type ChainV1 struct {
//	handles []HandleFuncV1
//}
//
//func (c ChainV1) Run(ctx *Context) {
//	for _, h := range c.handles {
//		next := h(ctx)
//		// 这种是中断执行
//		if !next {
//			return
//		}
//
//	}
//}

type Net struct {
	handlers []HandleFuncV1
}

func (c Net) Run(ctx *Context) {
	var wg sync.WaitGroup
	for _, hdl := range c.handlers {
		h := hdl
		if h.concurrent {
			wg.Add(1)
			go func() {
				h.Run(ctx)
				wg.Done()
			}()
		} else {
			h.Run(ctx)
		}
	}
	wg.Wait()
}

type HandleFuncV1 struct {
	concurrent bool
	handlers   *HandleFuncV1
	Net
}

func (h *HandleFuncV1) Run(ctx *Context) {
	var wg sync.WaitGroup
	for _, hdl := range h.Net.handlers {
		h := hdl
		if h.concurrent {
			wg.Add(1)
			go func() {
				h.Net.Run(ctx)
				wg.Done()
			}()
		}
	}

}

/*

AOP 方案
    面向切面编程。核心在于将横向关注点从业务中AOP(Aspect Oriented Programming)，剥离出来。
    横向关注点: 就是那些跟业务没啥关系，但是每个业务又必须要处理的。常见的有几类:
        可观测性:logging、metric和 tracing
        安全相关:登录、鉴权与权限控制
        错误处理:例如错误页面支持
        可用性保证: 熔断限流和降级等
    基本上Web 框架都会设计自己的 AOP 方案


面试要点
	什么是 AOP ?
		AOP 就是面切面编程，用于解决横向关注点问题，如可观测性问题、安全问题等
	什么是洋葱模式?
		形如洋葱，拥有一个核心，这个核心一般就是业务逻辑。而后在这个核心外面层层包裹，每一层就是一个 Middleware。一般用洋葱模式来无侵入式地增强核心功能，或者解决 AOP 问题。
	什么是责任链模式?
		不同的 Handler 组成一条链，链条上的每一环都有自己功能。一方面可以用责任链模式将复杂逻辑分成链条上的不同步骤，另外-府面也可以员活地在链条上添加新的
	Handelr.怎么实现?
		最简单的方案就是我们课程上进的这种函数式方案，还有一种是集中调度的模式


Middleware 例子- AccessLog
	在日常工作中，我们可能希望能够记录所有进来的请求以支持 DEBUG，也就是所谓的 AccessLog。
	Beego 里面支持了 AccessLog，Iris 也支持了AccessLog。
	这里我们提供一个简单的 AccessLog Middleware。

*/
