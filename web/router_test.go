package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// 第一个步骤是构造路由树
	// 第二个步骤是验证路由树
	testRouters := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/aaa",
		},
		{
			method: http.MethodGet,
			path:   "/*/aaa/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "login",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "login////",
		//},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	for _, route := range testRouters {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 接下来就在这里断言两者相等
	// 在这里断言路由树和你预期的一模一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				chlidren: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						chlidren: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						chlidren: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
							},
						},
						starChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				chlidren: map[string]*node{
					"order": &node{
						path: "order",
						chlidren: map[string]*node{
							"create": &node{
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": &node{
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	// 不能用assert.Equal(t, wantRouter, r) 因为里面有HandleFunc, HandleFunc是不可比的
	msg, ok := wantRouter.equal(&r)
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web: 路径必须以 / 开头")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/ad/da/d/", mockHandler)
	}, "web: 路径必须以 / 结尾")
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/ad/da//d/", mockHandler)
	}, "web: 不能有连续的 /")

	// 根节点重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "web: 路由冲突, 重复注册[/]")

	// 不同节点重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "web: 路由冲突, 重复注册[/a/b/c]")

	// 可用的 http method 要不要校验
	// mockHandler 是 nil 呢? 要不要校验
	// r.addRoute("aaa", "/a/b/c", mockHandler)

}

// 比较只读  不需要指针
// 为什么返回一个str(返回一个错误信息, 帮助我们排查问题)
// bool 是代表是否真的相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		// v, dst 要相等
		msg, ok := v.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}

	if n.starChild != nil {
		msg, ok := n.starChild.equal(y.starChild)
		if !ok {
			return msg, ok
		}
	}

	if len(n.chlidren) != len(y.chlidren) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	// 比较 handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler 不相等"), false
	}

	for path, c := range n.chlidren {
		dst, ok := y.chlidren[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}

		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

func TestRoute_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodDelete,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodDelete,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "login",
		//},
		//{
		//	method: http.MethodPost,
		//	path:   "login////",
		//},
	}

	r := newRouter()
	var mockHandler = func(ctx *Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{
			// 方法都不存在
			name:      "method not found",
			method:    http.MethodOptions,
			path:      "/order/detail",
			wantFound: false,
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "detail",
			},
		},
		{
			// 命中了, 但是没有handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			wantNode: &node{
				path: "order",
				chlidren: map[string]*node{
					"detail": &node{
						handler: mockHandler,
						path:    "detail",
					},
				},
			},
		},
		{
			// 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
				chlidren: map[string]*node{
					"order": &node{
						path: "order",
						chlidren: map[string]*node{
							"detail": &node{
								handler: mockHandler,
								path:    "detail",
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			n, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.wantFound, found)
			if !found {
				return
			}
			// 只能直接全Node去比, 里面还有handle函数, 要一个个去比
			assert.Equal(t, testCase.wantNode.path, n.n.path)

			msg, ok := testCase.wantNode.equal(n.n)
			assert.True(t, ok, msg)

		})
	}

}
