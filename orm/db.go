package orm

import "database/sql"

type DBOption func(db *DB)

// DB 是一个 sql.DB 的装饰器
type DB struct {
	r  *registry
	db *sql.DB
}

func Open(driver string, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

// OpenDB 为什么要拿出来一个OpenDB
//
//	从直觉上来说，我们可能只需要一个 Open方法，它会创建一个我们的 DB 实例。
//	实际上，因为用户可能自己创建了 sql.DB实例，所以我们要允许用户直接用 sql.DB来创建我们的 DB。
//	OpenDB 常用于测试，以及集成别的数据库中间件。我们会使用 sqlmock 来做单测试
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:  newRegistry(),
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func MustOpen(driver string, dataSourceName string, opts ...DBOption) *DB {
	res, err := Open(driver, dataSourceName, opts...)
	if err != nil {
		panic(err)
	}
	return res
}
