package orm

import (
	"context"
	"my-frame/orm/internal/errs"
	"strings"
)

type Selector[T any] struct {
	table string
	model *model
	where []Predicate
	sb    *strings.Builder
	args  []any
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb = &strings.Builder{}
	var err error
	s.model, err = parseModel(new(T))
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
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		if err := s.buildExpression(p); err != nil {
			return nil, err
		}

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

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
	case Predicate:
		// 在这里处理 p
		// p.left 构建好
		// p.op 构建好
		// p.right 构建好

		// 判断左边是否是表达式 是就加括号
		_, ok := exp.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}

		if err := s.buildExpression(exp.left); err != nil {
			return err
		}

		if ok {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(string(exp.op))
		s.sb.WriteByte(' ')

		// 判断左边是否是表达式 是就加括号
		_, ok = exp.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

	//switch left := p.left.(type) {
	//case Column:
	//	sb.WriteByte('`')
	//	sb.WriteString(left.name)
	//	sb.WriteByte('`')
	//	// 剩下不考虑
	//}
	//sb.WriteString(string(p.op))
	//switch right := p.right.(type) {
	//case value:
	//	sb.WriteByte('?')
	//	args = append(args, right.val)
	//	// 剩下不考虑
	//}

	case Column:

		fd, ok := s.model.fields[exp.name]
		// 字段不对, 或者说列不对
		if !ok {
			return errs.NewErrUnkonwField(exp.name)
		}
		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
		// 剩下不考虑

	case value:
		s.sb.WriteByte('?')
		s.addArg(exp.val)
	// 剩下不考虑

	default:
		//return fmt.Errorf("orm: 不支持的表达式类型 %v", expr)
		return errs.NewErrUnsupportedExpression(expr)
	}

	return nil

}

func (s *Selector[T]) addArg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}

	s.args = append(s.args, val)
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
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
