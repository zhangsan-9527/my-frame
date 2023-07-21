package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"my-frame/web/session"
	"time"
)

var (
	errorSessionNotFound = errors.New("session 不存在")
)

type StoreOption func(store *Store)

// hset
//
//	sess_id(rKey)    key(key)    value(val)
//
// map[string]map[string]string
type Store struct {
	prefix     string
	client     redis.Cmdable
	expiration time.Duration
}

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		expiration: time.Minute * 15,
		client:     client,
		prefix:     "sessid",
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func StoreWitjPrefix(prefix string) StoreOption {
	return func(store *Store) {
		store.prefix = prefix
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	rKey := redisKey(s.prefix, id)
	_, err := s.client.HSet(ctx, rKey, id, id).Result()
	if err != nil {
		return nil, err
	}

	_, err = s.client.Expire(ctx, rKey, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		id:     id,
		rKey:   rKey,
		client: s.client,
	}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	rKey := redisKey(s.prefix, id)
	ok, err := s.client.Expire(ctx, rKey, s.expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("session: id 对应的 session 不存在")
	}

	return nil

}

func (s *Store) Remove(ctx context.Context, id string) error {
	rKey := redisKey(s.prefix, id)
	_, err := s.client.Del(ctx, rKey).Result()
	return err
	// 代表的是 id 对应的session 不存在, 你没有删任何东西
	//if cnt == 0 {
	//}
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	// 自由决策要不要提前把 session 存储的用户数据一并拿过来
	// 1. 都不拿
	// 2. 只拿高频数据(热点数据)
	// 3. 都拿
	rKey := redisKey(s.prefix, id)
	cnt, err := s.client.Exists(ctx, rKey).Result()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, errorSessionNotFound
	}

	return &Session{
		id:     id,
		rKey:   rKey,
		client: s.client,
	}, nil
}

type Session struct {
	id     string
	rKey   string
	prefix string
	client redis.Cmdable
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, err := s.client.HGet(ctx, s.rKey, key).Result()
	return val, err
}

func (s *Session) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`

	res, err := s.client.Eval(ctx, lua, []string{s.rKey}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errorSessionNotFound
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}

func redisKey(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
