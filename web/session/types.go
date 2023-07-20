package session

import (
	"context"
	"net/http"
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
