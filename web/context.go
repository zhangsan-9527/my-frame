package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

//var (
//	jsonUserNumber            = true
//	jsonDisallowUnknownFields = true
//)

type Context struct {
	Req *http.Request

	// Resp 如果用户直接使用这个
	// 那个他们就绕开了 RespData 和 RespStatusCode 两个
	// name部分 middleware 无法运作
	Resp http.ResponseWriter

	// 这个如要是为了 middleware 读写用的
	// 响应数据
	RespData []byte
	// 响应状态码
	RespStatusCode int

	PathParams map[string]string

	//Ctx context.Context

	// 缓存的数据
	queryValues url.Values

	MatchedRoute string

	// cookieSameSite http.SameSite
}

//func (c *Context) ErrPage() {
//
//}

// SetCookie 设置Cookie
func (c *Context) SetCookie(ck *http.Cookie) {
	// 不推荐
	//ck.SameSite = c.cookieSameSite
	http.SetCookie(c.Resp, ck)
}

func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

// RespJSON 处理输出-JSON响应
func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespData = data
	c.RespStatusCode = status

	// 放到最后响应之前 (root(ctx)的之前)传入
	//c.Resp.WriteHeader(status)
	//c.Resp.Header().Set("Content-Type", "application/json")
	//c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	//n, err := c.Resp.Write(data)
	//if n != len(data) {
	//	return errors.New("web: 未写入全部数据")
	//}
	return err
}

// BindJSON 处理输入: 解决大多数人的需求(处理body 转成json)
func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web: 输入为 nil")
	}

	if c.Req.Body == nil {
		return errors.New("web: body 为 nil")
	}

	// 不要这样写
	//bs, _ := io.ReadAll(c.Req.Body)
	//json.Unmarshal(bs, val)
	decoder := json.NewDecoder(c.Req.Body)

	// useNumber => 数字就会用 Number 来表示
	// 否则默认是float64
	//if jsonUserNumber {
	//	decoder.UseNumber()
	//}

	// 如果要是有一个未知的字段, 就会报错
	// 比如 User 只有 Name 和 Email 两个字段, JSON 里面额外多了一个 Age 字段, 就会报错
	//decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

// FormValue 从表单里拿数据
func (c *Context) FormValue(key string) (string, error) {
	//r.PostForm == nil  ParseForm()方法中的  不用担心重复ParseForm
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	//vals, ok := c.Req.Form[key]
	//if !ok {
	//	return "", errors.New("web: key 不存在")
	//}
	//return vals[0], nil
	return c.Req.FormValue(key), nil
}

// QueryValue 处理输入-查询参数
// Query 和表单比起来, 它没有缓存
func (c *Context) QueryValue(key string) (string, error) {
	// 问题: 这个 ParseQuery 每次都会解析一遍查询串，在我们这里就是 name=xiaoming&age=18 字符串。
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok || len(vals) == 0 {
		return "", errors.New("web: key 不存在")
	}
	return vals[0], nil
	// 用户区别不出来是真的有值, 但是值恰好是空字符串
	// 还是没有值
	//return c.queryValues.Get(key), nil
}

func (c *Context) QueryValueV1(key string) StringValue {
	// 问题: 这个 ParseQuery 每次都会解析一遍查询串，在我们这里就是 name=xiaoming&age=18 字符串。
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok || len(vals) == 0 {
		return StringValue{
			err: errors.New("web: key 不存在"),
		}
	}
	return StringValue{
		val: vals[0],
	}
	// 用户区别不出来是真的有值, 但是值恰好是空字符串
	// 还是没有值
	//return c.queryValues.Get(key), nil
}

// PathValue 处理路径参数
func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web; key 不存在")
	}
	return val, nil
}

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web; key 不存在"),
		}
	}
	return StringValue{
		val: val,
	}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

/*
	处理输入要解决的问题:
		反序列化输入: 将 Body 字节流转换成一个具体的类型
		处理表单输入: 可以看做是一个和JSON 或者 XML 差不多的一种特殊序列化方式
		处理查询参数: 指从URL 中的查询参数中读取值，并且转化为对应的类型
		处理路径参数: 读取路径参数的值，并且转化为具体的类型
		重复读取body: http.Request 的 Body 默认是只能读取一次，不能重复读取的
		读取Header: 从Header 里面读取出来特定的值，并且转化为对应的类型
		模糊读取: 按照一定的顺序，尝试从 Body、Header、路径参数或者 Cookie 里面读取值，并且转化为特定类型

  处理输出要解决的问题:
		序列化输出: 按照某种特定的格式输出数据，例如JSON 或者 XML
		渲染页面:要考虑模板定位、命名和渲染的问题
		处理状态码:允许用户返回特定状态码的响应，例如 HTTP 404
		错误页面: 特定 HTTP Status 或者 error 的时候，能够重定向到一个错误页面，例如404 被重定向到首页
		设置 Cookie : 设置 Cookie 的值
		设置 Header: 往 Header 里面放一些东西



	处理输入-JSON 输入控制选项:
		严谨地说，如果用户有这种需求，那么他可能需要的是整个应用级别上控制 useNumber 或者 disableUnknownFields
		单一 HTTPServer 实例上控制 useNumber 或者disableUnknownFields
		特定路径下，例如在 /user/** 都为 true 或者都为false
		特定路由下，例如在 /user/details 下都为 true 或者false
		不同情况下框架支持的方式也不一样:
			整个应用级别:维持两个全局变量
			HTTPServer 级别:在HTTPServer 里面定义两个字段
			引入类型 BindJSONOpt 的方法, 用户灵活控制

		没必要(用户可以自己来)
	func BindJSON(ctx *Context, val any, userNumber bool, disallow bool) {

	}


	处理输入-表单输入:
	表单在 Go 的 http.Request 里面有两个Form: 一个是 URL 里面的查询参数和 PATCH、POST、PUT的表单数据
	PostForm:PATCH、POST 或者 PUT body 参数
	但是不管使用哪个，都要先使用 ParseForm 解析表单数据
		处理输入-Form 和PostForm
			表单在 Go 的 http.Request 里面有两个。
				Form:基本上可以认为，所有的表单数据都能拿到
				PostForm: 在编码是 X-www-form-urlencoded 的时候才能拿到
		实际中我是不建议大家使用表单的，我一般只在 Body 里面用JSON 通信，或者用 protobuf 通信。




	处理输入-查询参数:
	所谓的查询参数，就是指在 URL 问号之后的部分例如 URL:http://localhost:8081/form?name=xiaoming&age=18
	那么查询参数有两个:name=xiaoming 和age=18
	前面我们注意到，如果要是调用了 ParseForm，那么这部分也可以在 Form 里面找到
	问题: 这个 ParseQuery 每次都会解析一遍查询串，在我们这里就是 name=xiaoming&age=18 字符串。
*/

/*
	处理输出-JSON 响应:
	这种设计非常简单，就是帮助用户将 val 转化一下其它格式的输出也是类似的写法
	也可以考虑提供一个更加方便的 RespJSONOK 方法这个就是看个人喜好了

	这里有一个问题，如果 val 已经是 string 或者[]byte了，那么用户该怎么办?
		val是 string 或者[]byte 肯定不需要调用 RespJSON了，自己直接操作 Resp

	处理输出-需要支持错误页面吗?
	我们通常有一个需求，是如果一个响应返回了 404那么应该重定向到一个默认页面，比如说重定位到首页。

	那么该怎么处理?

	这里有一个很棘手的点: 不是所有的 404 都是要重定向的。比如说你是异步加载数据的 RESTful 请求，例如在打开页面之后异步加载用户详情，即便 404 了也不应该重定向

*/

/*
	Context面试题:
		Context 总结Context 是线程安全的吗?
			显然不是，和路由树不是线程安全的理由不太一样。
			Context 不需要保证线程安全，是因为在我们的预期里面，这个 Context 只会被用户在一个方法里面使用，而且不应该被多个 goroutine 操作。
			对绝达多数大来说，他们不需要一个线程安全的 Context。即便真要线程安全，我们也可以提供一个装饰器，让用户在使用前手动创建装饰器来转换一下
		Context 总结Context 为什么不设计为接口?
			目前来看，看不出来设计为接口的必要性
			Echo 设计为接口，但是只有一个实现，就足以说明设计为接口有点过度设计的感觉，
			即便 iris 设计为接口，而且允许用户提供自定义实现，但是看起来也不是那么有用

		Context 总结-Context 能不能用泛型?
			我们已经在好几个地方用过泛型了，在 Context 里面似乎也有使用泛型的场景，例如说处理表单数据、查询参数、路径参数......
			答案是不能。因为 Go 泛型有一个限制，结构体本身可以是泛型的，但是它不能声明泛型方法
			同样的道理，stringValue 也不能声明为泛型

	面试要点
		能不能重复读取 HTTP 协议的 Body 内容? 原生 API 是不可以的。但是我们可以通过封装来允许重复读取，核心步骤是我们将 Body 读取出来之后放到一个地方，后续都从这个地方读取。
		能不能修改 HTTP 协议的响应? 原生 API 也是不可以的。但是可以用我们的 RespData 这种机制，在最后再把数据刷新到网络中，在刷新之前，都可以修改。
		Form 和 PostForm 的区别。课程上讲过，一个无聊的问题，正常的情况下你的 API优先使用 Form 就不太可能出错。
		Web 框架是怎么支持路径参数的? 简单得很，We 框架在发现匹配上了某个路径参数之后，将这段路径记录下来作为路径参数的值，这个值默认是 string 类型，用户自己有需要就可以转化为不同的类型
*/
