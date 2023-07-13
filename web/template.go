package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {
	// Render 渲染页面
	// tplName 模板的名字, 按名索引
	// data 渲染页面用的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	// 选人页面, 数据写入到 writer 里面
	// Render(ctx, "aa", map[]{}, repsonseWriter)优点: 不用返回数据   缺点: 不好测试 repsonseWriter 直接写入http
	//	缺点 保持住了 RespData 的语义, 这意味着其他中间件可以篡改这个页面, 例如直接替换为错误页面等(但是下边这个就不行)
	//Render(ctx context.Context, tplName string, data any, writer io.Writer) error

	// 不需要, 让具体实现管自己的模板
	//AddTemplate(tlpName string, tpl []byte) error

	// 用这个Context 没有问题
	//Render(ctx Context)
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {
	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

// 没必要二次封装
//func (g *GoTemplateEngine) ParseGlob(pattern string) error {
//	var err error
//	g.T, err = template.ParseGlob(pattern)
//	return err
//}

/**
面试要点(一)
模板在日常工作中非常好用，面试的时候则主要聚焦在模板的语法上:
	模板的基本语法:变量声明、方法调用、循环、条件判断、操作符，以及一个比较常见的，i就是怎么在模板里面实现 for ...i... 的循环。
	什么是前缀表达式 (+ b c)?也就是模板里面的那种语法，和中缀表达式 (b+c)比起来，它更加贴近计算机的计算原理。所以模板用了前缀表达式，能够简化模板引擎的设计和实现。
	http/template 和 text/template 有什么区别: 前者多了对 HTTP 的支持，加强了安全性，例如特殊字符转义等。http/template 能够满足绝大多数页面渲染的要求
	模板中的pipeline是什么? 一串命令，pipeline之间可以通过|连在一起组成更加复杂的pipeline。单个命令可以是声明变量，也可以是方法调用。


面试要点(二)
	怎么支持错误页面?也就是如果响应码是 500 之类的，返回一个默认的错误页面。如果 Web 框架不支持算改响应，那么就毫无办法，只能自己在业务代码里面处理。如果支持篡改，那么就可以利用Middleware 来统一处理。
	template 是怎么实现的? template 有点像是解释执行，也就是模板引警在读到模板语法的时候就开始解析，并返回结果

*/
