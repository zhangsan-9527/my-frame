package valuer

import (
	"database/sql"
	"my-frame/orm/internal/errs"
	"my-frame/orm/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model

	// 对应于 T 的指针
	val any
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(model *model.Model, val any) Value {
	return unsafeValue{
		model: model,
		val:   val,
	}
}

func (r unsafeValue) SetColumns(rows *sql.Rows) error {

	// 怎么知道 SELECT 出来了那些列?
	// 拿到了 SELECT的列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	var vals []any
	// 起始地址
	address := reflect.ValueOf(r.val).UnsafePointer()

	for _, c := range cs {
		// c 是列名
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		// 要计算字段的地址
		// 字段的地址 = 起始地址 + 偏移量
		fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)

		// 反射在特定的低智商, 创建一个特定类型的实例
		// 反射创建一个实例, 这里创建的实例是原本类型的指针类型
		// 例如 fd.Type = int . 那么val 是 *int
		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}

	err = rows.Scan(vals...)
	return err
}
