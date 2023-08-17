package valuer

import (
	"database/sql"
	"my-frame/orm/internal/errs"
	"my-frame/orm/model"
	"reflect"
)

type reflectValue struct {
	model *model.Model

	// 对应于 T 的指针
	val any
}

var _ Creator = NewReflectValue

func NewReflectValue(model *model.Model, val any) Value {
	return reflectValue{
		model: model,
		val:   val,
	}
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	// 在这里, 继续处理结果集

	// 怎么知道 SELECT 出来了那些列?
	// 拿到了 SELECT的列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	// 第一个问题: 类型要匹配
	// 第二个问题: 顺序要匹配
	// 怎么处理 cs? 怎么利用cs解决顺序问题和类型问题?

	// 通过 cs 来构造 vals
	vals := make([]any, 0, len(cs))
	valElems := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		// c 是列名
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 反射创建一个实例, 这里创建的实例是原本类型的指针类型
		// 例如 fd.Type = int . 那么val 是 *int
		val := reflect.New(fd.Typ)
		vals = append(vals, val.Interface())
		// 记得要调用 Elem, 因为 fd.Type = int . 那么val 是 *int
		valElems = append(valElems, val.Elem())
		// 通过 columnMap 优化运行速度 (牺牲空间)
		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		// 反射创建一个实例
		//		// 这里创建的实例是原本类型的指针类型
		//		// 例如 fd.Type = int . 那么val 是 *int
		//		val := reflect.New(fd.typ)
		//		vals = append(vals, val.Interface())
		//	}
		//}
	}
	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	// 想办法把 vals 塞进去 结果 tp 里面
	tpValueElem := reflect.ValueOf(r.val).Elem()
	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		tpValueElem.FieldByName(fd.GoName).Set(valElems[i])
		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
		//	}
		//}
	}

	return err
}
