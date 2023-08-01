package orm

import (
	"context"
	"reflect"
)

// goland重命名 shift + F6

type Deleter[T any] struct {
	builder
	table string
	where []Predicate
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sb.WriteString("DELETE * FROM ")
	// 我怎么把表名拿到
	var t T
	if d.table == "" {
		d.sb.WriteByte('`')
		d.sb.WriteString(reflect.TypeOf(t).Name())
		d.sb.WriteByte('`')
	} else {
		//segs := strings.Split(s.table, ".")
		//d.sb.WriteByte('`')
		//d.sb.WriteString(segs[0])
		//d.sb.WriteByte('`')
		//d.sb.WriteByte('`')
		//d.sb.WriteByte('.')
		//d.sb.WriteByte('`')
		//d.sb.WriteByte('`')
		//d.sb.WriteString(segs[1])
		//d.sb.WriteByte('`')
		d.sb.WriteString(d.table)
	}

	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}

		if err := d.buildExpression(p); err != nil {
			return nil, err
		}

	}

	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = predicates
	return d
}

func (d *Deleter[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Deleter[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
