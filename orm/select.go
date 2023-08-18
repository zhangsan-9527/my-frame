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
		sb.WriteString(s.model.TableName)
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

//func (s *Selector[T]) GetV1(ctx context.Context) (*T, error) {
//	q, err := s.Build()
//	// 这个是构造 SQL 失败错误
//	if err != nil {
//		return nil, err
//	}
//
//	db := s.db.db
//	// 在这里, 就是要发起查询, 并且处理结果集
//	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
//	// 这个是查询数据库时的错误
//	if err != nil {
//		return nil, err
//	}
//
//	// 要确认有没有数据
//	if !rows.Next() {
//		// 要不要返回 error?
//		// 要返回error, 和 sql 语义包保持一致 sql.ErrNoRows
//		return nil, ErrNoRows
//	}
//
//	// 在这里, 继续处理结果集
//
//	// 怎么知道 SELECT 出来了那些列?
//	// 拿到了 SELECT的列
//	cs, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//
//	var vals []any
//	tp := new(T)
//	// 起始地址
//	address := reflect.ValueOf(tp).UnsafePointer()
//
//	for _, c := range cs {
//		// c 是列名
//		fd, ok := s.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//
//		// 要计算字段的地址
//		// 字段的地址 = 起始地址 + 偏移量
//		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
//
//		// 反射在特定的低智商, 创建一个特定类型的实例
//		// 反射创建一个实例, 这里创建的实例是原本类型的指针类型
//		// 例如 fd.Type = int . 那么val 是 *int
//		val := reflect.NewAt(fd.Typ, fdAddress)
//		vals = append(vals, val.Interface())
//	}
//
//	err = rows.Scan(vals...)
//	return tp, err
//}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {

	q, err := s.Build()
	// 这个是构造 SQL 失败错误
	if err != nil {
		return nil, err
	}

	db := s.db.db
	// 在这里, 就是要发起查询, 并且处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args...)
	// 这个是查询数据库时的错误
	if err != nil {
		return nil, err
	}

	// 要确认有没有数据
	if !rows.Next() {
		// 要不要返回 error?
		// 要返回error, 和 sql 语义包保持一致 sql.ErrNoRows
		return nil, ErrNoRows
	}

	// 接口定义好之后, 就两件事, 一个是用新接口的方法改造上层,
	tp := new(T)

	val := s.db.creator(s.model, tp)
	err = val.SetColumns(rows)

	// 一个就是提供不同的实现

	return tp, err

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
