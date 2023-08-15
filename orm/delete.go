package orm

import (
	"context"
	"strings"
)

// goland重命名 shift + F6

type Deleter[T any] struct {
	builder
	table string
	where []Predicate
	r     *registry
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.sb = &strings.Builder{}
	var err error
	d.model, err = d.r.Register(new(T))
	if err != nil {
		return nil, err
	}
	sb := d.sb
	sb.WriteString("DELETE * FROM ")
	// 我怎么把表名拿到
	if d.table == "" {
		sb.WriteByte('`')
		sb.WriteString(d.model.TableName)
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
		sb.WriteString(d.table)
	}

	if len(d.where) > 0 {
		sb.WriteString(" WHERE ")
		if err = d.buildPredicates(d.where); err != nil {
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
