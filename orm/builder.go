package orm

import (
	"my-frame/orm/internal/errs"
	"strings"
)

type builder struct {
	sb    *strings.Builder
	args  []any
	model *Model
	db    *DB
}

//type Predicates []Predicate
//
//func (p *Predicates) build(s *strings.Builder) error {
//	// 写在这里
//}

//type predicates struct {
//	// WHERE 或者 HAVING
//	prefix string
//	ps     []Predicate
//}
//func (p *predicates) build(s *strings.Builder) error {
//	// 包含拼接WHERE 或者 HAVING部分
//	// 写在这里
//}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}

	return b.buildExpression(p)
}

func (b *builder) buildExpression(expr Expression) error {
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
			b.sb.WriteByte('(')
		}

		if err := b.buildExpression(exp.left); err != nil {
			return err
		}

		if ok {
			b.sb.WriteByte(')')
		}

		b.sb.WriteByte(' ')
		b.sb.WriteString(string(exp.op))
		b.sb.WriteByte(' ')

		// 判断左边是否是表达式 是就加括号
		_, ok = exp.right.(Predicate)
		if ok {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			b.sb.WriteByte(')')
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

		fd, ok := b.model.FieldMap[exp.name]
		// 字段不对, 或者说列不对
		if !ok {
			return errs.NewErrUnknownField(exp.name)
		}
		b.sb.WriteByte('`')
		b.sb.WriteString(fd.ColName)
		b.sb.WriteByte('`')
		// 剩下不考虑

	case value:
		b.sb.WriteByte('?')
		b.addArg(exp.val)
	// 剩下不考虑

	default:
		//return fmt.Errorf("orm: 不支持的表达式类型 %v", expr)
		return errs.NewErrUnsupportedExpression(expr)
	}

	return nil

}

func (b *builder) addArg(val any) {
	if b.args == nil {
		b.args = make([]any, 0, 4)
	}

	b.args = append(b.args, val)
}
