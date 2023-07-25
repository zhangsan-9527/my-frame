package orm

// 衍生类型
type op string

const (
	opEq  op = "="
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (o op) string() string {
	return string(o)
}

// 这种叫做别名
//type op = string

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

// Eq("id", "12")
// Eq(sub.id, "id", 12)
// Eq("sub.id", 12)
//func Eq(column string, arg any) Predicate {
//	return Predicate{
//		Conlum: column,
//		Op:     "=",
//		Arg:    arg,
//	}
//}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{
		name: name,
	}

}

// C("id").Eq(12)
// 子查询 sub.C("id").Eq(12)
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		right: value{
			val: arg,
		},
	}
}

func (c Column) expr() {}

// Not(C("name").Eq("Tom"))
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}

}

// C("id").Eq(12).And(C("name").Eq("Tom"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}

}

// C("id").Eq(12).Or(C("name").Eq("Tom"))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}

}

func (Predicate) expr() {}

// Expression 是一个标记接口, 代表表达式
type Expression interface {
	expr()
}

type value struct {
	val any
}

func (value) expr() {

}
