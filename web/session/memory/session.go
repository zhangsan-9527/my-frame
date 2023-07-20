package memory

import (
	"context"
	"errors"
	"my-frame/web/session"
	"sync"
)

var (
	//ErrkeyNotFound sentinel error, 预定义错误
	errkeyNotFound = errors.New("session: 找不到key")
)

type Store struct {
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Remove(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	//TODO implement me
	panic("implement me")
}

// Session 基于内存实现
type Session struct {
	id string

	//mutex sync.RWMutex
	//values map[string]any

	valus sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.valus.Load(key)
	if !ok {
		//return nil, fmt.Errorf("%w, key %s", errkeyNotFound, key)
		return nil, errkeyNotFound
	}

	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	s.valus.Store(key, val)
	return nil

}

func (s *Session) ID() string {
	return s.id
}
