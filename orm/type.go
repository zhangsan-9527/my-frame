package orm

import (
	"context"
	"database/sql"
)

// Querier 用于 Select 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)

	// 这种设计形态也可以
	//Get(ctx context.Context) (T, error)
	//GetMulti(ctx context.Context) ([]T, error)
}

// Executor 用于 INSERT, DELEIE 和 UPDATE
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)

	// 这样也可以(顾虑 在AOP的时候修改这个Query 修改不了)
	//Build() (Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
