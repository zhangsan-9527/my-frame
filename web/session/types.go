package session

import (
	"context"
	"net/http"
)

var (
// ErrkeyNotFound sentinel error, 预定义错误
// ErrkeyNotFound = errors.New("")
)

// Store 管理 Session 本身
type Store interface {
	// Generate 生成
	// session 对应的 ID 谁来指定
	// 要不要在接口维度上设置超时时间, 以及要不要让 Store 内部去生成ID, 都是可以自由决策的
	Generate(ctx context.Context, id string) (Session, error)

	Refresh(ctx context.Context, id string) error

	Remove(ctx context.Context, id string) error

	// Get 获取
	Get(ctx context.Context, id string) (Session, error)

	// 这种设计意味着要拿到这个session 才能刷新, 上面拿到id就可以刷新
	// Refresh(ctx context.Context, sess Session) error
}

type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

type Propagator interface {

	// Inject 将 session id 注入到里面
	// 必须是等幂的
	Inject(id string, writer http.ResponseWriter) error
	// Extract 将 session id 从 http.Request 中提取出来
	// 例如从 cookie 中将 session id 提取出来
	Extract(req *http.Request) (string, error)

	// Remove 将 session id 从 http.ResponseWriter 中删除
	// 例如删除对应的 cookie
	Remove(writer http.ResponseWriter) error
}

/*
	面试要点 (一)
		Session 是什么? 一种在服务端维护用户状态的机制
		Cookie 和 Session 的对比? 其实两者都可以看做是维护用户状态的机制，只不过一个是在客户端个是在服务端
		什么时候刷新 Session? 用户活跃的时候就可以刷，前端定时刷，或者后端每次收到请求都刷。但是这里要注意，频繁刷新 Session 可能给 Session 存储的地方带来庞大的压力，例如 Redis
		怎么实现一个 Session? 简单来说就是要构建出 Session 和 Store 两个抽象，其它就随意了
		怎么生成 session id? 看业务需求了，最简单的就是 UUID，高级一点的就用特定的算法进行加密然后将业务需要的数据编码进去。
*/

/*
	Session 安全性

	实际上 Session 这种 session id 的认证还是比较弱的，如果没有做一些安全措施，那么不管是谁拿到 session id，服务器都认一认 session id 不认人。

	一些可行的保护 session id 的方案在使用 Cookie 的时候，同时设置 http only 和secure 选项，限制 Cookie 只能在 HTTPS 协议里面被传输。

	在 session id 编码的时候带上一些客户端信息如agent 信息、MAC 地址之类的。如果服务端检测到 sessionid所携带的这些信息发生了变化，就要求用户重新登录
*/

/*
	为什么你不管 session id 生牌Session 总结
		session id 在当下的生成策略可以说是五花八门，实在管不过来。session id 的生成策略可以要考虑:
		是否要包含业务信息:
		·是:
		> 编码什么业务信息，用户决定，接口难以设计
		> 编码用什么算法，用户决定，接口更难以设计 // jwt就是包含了必要信息的特殊(非敏感信息)的sessid
		否: 你都不用包含什么信息了，UUID 搞一下就可以了，用不着我来管

	所以要管的管不了，不要管的一那就不管

*/

/*

	Session 总结 - 什么时候刷新Session
		如果我们一直不刷新 Session，时间一到，Session就会直接过期。即便此时用户还在操作，也会导致用户直接退出登录状态

		最直接的做法就是每次收到一个请求都刷新。体量不大的时候就这么干，简单直接好维护。类似的做法还可以是前端定时心跳刷新，例如 5 秒刷新一次

		高端一点的就是维护长短两个 token，可以看做是两个过期时间不一样的 Session。每次先检查短过期时间的 token，找不到就去找长过期时间的token，这时会重新生成一个短 token。

*/

/*
代码获取测试
*/
