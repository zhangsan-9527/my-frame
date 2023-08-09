package orm

import (
	"context"
	"strings"
)

type Selector[T any] struct {
	table  string
	where  []Predicate
	having []Predicate
	builder
}

// 不允许在方法中引入泛型
//func (db *DB)NewSelector[T any]() *Selector[T]{}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: builder{
			sb: &strings.Builder{},
			db: db,
		},
	}
}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.model, err = s.db.r.Register(new(T))
	if err != nil {
		return nil, err
	}
	sb := s.sb
	sb.WriteString("SELECT * FROM ")
	// 我怎么把表名拿到
	if s.table == "" {
		sb.WriteByte('`')
		sb.WriteString(s.model.tableName)
		sb.WriteByte('`')
	} else {
		//segs := strings.Split(s.table, ".")
		//sb.WriteByte('`')
		//sb.WriteString(segs[0])
		//sb.WriteByte('`')
		//sb.WriteByte('`')
		//sb.WriteByte('.')
		//sb.WriteByte('`')
		//sb.WriteByte('`')
		//sb.WriteString(segs[1])
		//sb.WriteByte('`')
		sb.WriteString(s.table)
	}

	if len(s.where) > 0 {
		sb.WriteString(" WHERE ")
		if err = s.buildPredicates(s.where); err != nil {
			return nil, err
		}
		//p := s.where[0]
		//for i := 1; i < len(s.where); i++ {
		//	p = p.And(s.where[i])
		//}
		//
		//if err = s.buildExpression(p); err != nil {
		//	return nil, err
		//}

	}

	sb.WriteByte(';')
	return &Query{
		SQL:  sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

// ids := []int{1, 2, 3}
// s.Where("id in (?, ?, ?)", ids)

// s.Where("id in (?, ?, ?)", ids...)
// golint-ci
//func (s *Selector[T]) Where(query string, args ...any) *Selector[T] {
//
//}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {

	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 在这里, 就是要发起查询, 并且处理结果集
	_, err = db.QueryContext(ctx, q.SQL, q.Args...)
	// 在这里, 继续处理结果集
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

/*

SELECT 起步面试要点
	GORM 是如何构造 SQL的?
		在GORM 里面主要有四个抽象: Builder、Expression、Clause 和Interface。简单一句话概括 GORM 的设计思路就是 50L 的不同部分分开构造，最后再拼接在一起。

	什么是 Builder 模式? 能用来干什么?
		用我们的 ORM 的例子就可以，Builder 模式尤其适合用于构造复杂多变的对象。

	在 ORM 框架使用泛型有什么优点?
		能用来约束用户传入的参数或者用户希望得到的返回值，加强类型安全。

	另外有些时候面试官可能会让你手写 SOL，你需要顺便记住一下 SELECT 语句常见的部分是怎么写的

*/
