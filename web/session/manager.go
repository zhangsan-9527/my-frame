package session

import (
	"github.com/google/uuid"
	"my-frame/web"
)

type Manage struct {
	Propagator

	Store

	CtxSessKey string
}

func (m *Manage) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	//ctx.Req.Context().Value(m.CtxSessKey)
	val, ok := ctx.UserValues[m.CtxSessKey]
	if ok {
		return val.(Session), nil
	}
	// 尝试缓存住 session
	sessId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}

	sess, err := m.Get(ctx.Req.Context(), sessId)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.CtxSessKey] = sess
	//ctx.Req = ctx.Req.WithContext(context.WithValue(ctx.Req.Context(), m.CtxSessKey, sess))  // 复制问题 影响性能  还有 因为context.Context的一个特性
	return sess, err

}

func (m *Manage) InitSession(ctx *web.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 注入进去 HTTP 响应里面
	err = m.Inject(id, ctx.Resp)
	return sess, err
}

func (m *Manage) RefreshSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	return m.Refresh(ctx.Req.Context(), sess.ID())
}

func (m *Manage) RemoveSession(ctx *web.Context) error {
	sess, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), sess.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
