package memory

import (
	"context"
	"errors"
	"fmt"
	cache "github.com/patrickmn/go-cache"
	"my-frame/web/session"
	"sync"
	"time"
)

var (
	//ErrkeyNotFound sentinel error, 预定义错误
	errorKeyNotFound     = errors.New("session: 找不到key")
	errorSessionNotFound = errors.New("session: 找不到sess")
)

type Store struct {
	mutex      sync.RWMutex
	sessions   *cache.Cache
	expiration time.Duration
}

// 其他语言传一个int类型 的秒或者毫秒  加一个单位
//func NewStore(ms int, unit TimeUnit) {
//
//}

func NewStore(expiration time.Duration) *Store {
	return &Store{

		sessions:   cache.New(expiration, time.Second),
		expiration: expiration,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &Session{
		id: id,
	}
	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.sessions.Get(id)
	if !ok {
		return fmt.Errorf("session: 该id对应的session不存在 %s", id)
	}
	s.sessions.Set(id, val, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorSessionNotFound
	}

	return sess.(*Session), nil

}

// Session 基于内存实现
type Session struct {
	id string

	//mutex sync.RWMutex
	//values map[string]any

	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.values.Load(key)
	if !ok {
		//return nil, fmt.Errorf("%w, key %s", errkeyNotFound, key)
		return nil, errorKeyNotFound
	}

	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.values.Store(key, val)
	return nil

}

func (s *Session) ID() string {
	return s.id
}
