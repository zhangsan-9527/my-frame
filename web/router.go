package web

import (
	"fmt"
	"strings"
)

// router 用来支持对路由树的操作
// 代表路由树(森林)
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// GET 有一棵树, POST也有一棵树

	// http method => 路由树根节点
	trees map[string]*node
}

//type tree struct {
//	root *node
//}

/*
全静态匹配 -接口设计
关键类型:
router:维持住了所有的路由树，它是整个路由注册和查找的总入口。router 里面维护了一个 map，是按照HTTP方法来组织路由树的
node: 代表的是节点。它里面有一个 children 的map 结构，使用 map 结构是为了快速查找到子节点.
*/
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 加一些限制使用户不能乱传

// 已经注册了的路由，无法被覆盖，例如 /user/home注册两次，会冲突
// path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id和 /user/* 冲突
// 同名路径参数，在路由匹配的时候，值会被覆盖，例如/user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456

// path  必须以 / 开头, 不能以 / 结尾, 中间也不能有连续的 //
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	// path为空的校验(限制)
	if path == "" {
		panic("web: 路径不能为空字符串")
	}

	// path开头结尾的校验
	if path[0] != '/' {
		panic("web: 路径必须以 / 开头")
	}

	// path结尾校验
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路径不能以 / 结尾")
	}

	// 中间连续 //, 可以用 strings.contanins("//")

	// 首先找到树
	root, ok := r.trees[method]
	if !ok {
		// 说明还没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 根节点特殊处理一下
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突, 重复注册[/]")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}

	// 切割这个path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		// 中间连续 //, 可以用 strings.contanins("//")
		if seg == "" {
			panic("web: 不能有连续的 /")
		}
		// 递归下去, 找准位置
		// 如果中途有节点不存在, 你就要创建出来
		chlid := root.childOfCreate(seg)
		root = chlid
	}

	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突, 重复注册[%s]", path))
	}
	root.handler = handleFunc
	root.route = path

}

func (n *node) childOfCreate(seg string) *node {

	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配, 已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配, 已有路径参数")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}

	if n.chlidren == nil {
		n.chlidren = map[string]*node{}
	}
	res, ok := n.chlidren[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.chlidren[seg] = res
	}
	return res
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	// 基本上是不是也是沿着树深度查找下去?
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 对根节点进行处理
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	// 这里把前置和后置的 / 都去掉
	path = strings.Trim(path, "/")

	// 按照斜杠切割
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.chlidOf(seg)
		if !found {
			return nil, false
		}

		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path 是 :id 这种形式, 把: 去掉
			pathParams[child.path[1:]] = seg
		}

		root = child
	}

	// 代表我确实有这个节点, 但是节点是不是用户注册的 有handler的  就不一定了
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true

	// 代表我确实有这个节点, 并且判断是否有handler
	//return root, root.handler != nil

}

// 优先考虑静态匹配, 匹配不上, 再考虑通配符匹配
// 第一个返回值是子节点
// 第二个是标记是否是路径参数
// 第三个标记命中了没有
func (n *node) chlidOf(path string) (*node, bool, bool) {
	if n.chlidren == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}

	child, ok := n.chlidren[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok

}

type node struct {
	route string

	path string

	// 静态节点
	// 子 path 到子节点的映射
	chlidren map[string]*node

	// 通配符匹配 *
	starChild *node

	// 路径参数
	paramChild *node

	// 缺一个代表用户注册的业务逻辑
	handler HandleFunc
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

/*
面试要点(1):
路由树算法?
核心就是前缀树。前缀的意思就是，两个节点共同的前缀，将会被抽取出来作为父亲节点.在我们的实现里面，是按照 / 来切割，每一段作为一个节点。

路由匹配的优先级?
本质上这是和 Web 框架相关的。在我们的设计里面是静态匹配 > 路径参数>通配符配。


路由查找会回溯吗?
这也是和 Web 框架相关的，我们在课程上是不支持的。在这里可以简单描述可回溯和不可回溯之间的区别，可以是田课程体cser/123/home 和 /user/*\/*。我这里不支持是因为这个特性非常鸡肋。

Web 框架是怎么组织路由树的?
一个 HTTP 方法一颗路由树，也可以考虑一颗路由树，每个节点标记自己支持的 HTTP 方法。在课程中可以看到，前者是比较主流的
*/

/*
面试要点(2):
路由查找的性能受什么影响? 或者说怎么评估路由查找的性能?
核心是看路由树的高度，次要因素是路由树的宽度 (想想我们的 children 字段)。

路由树是线程安全的吗?
严格来说也是跟 Web 框架相关的。大多数都不是线程安全的，这是为了性能。所以才要求大家一定要先注册路由，后启动 Web 服务器。如果你有运行期间动态添加路由的需求，只需要利用装饰器模式，就可以将一个线程不安全的封装为线程安全的路由树。

具体匹配方式的实现原理。
课程上我们讨论了静态匹配、通配符匹配和路径匹配，64e业里面要求大家照着实现一个正则匹配。其实核心就是划定优先级，然后一种种匹配方式挨个匹配过去。


*/
