package valuer

import "database/sql"

// Value 这样设计主要适用于INSERT
type Value interface {
	SetColumns(rows *sql.Rows) error
}

// Creator 函数式工厂接口
type Creator func(entity any) Value

//type ValuerV1 interface {
//	SetColumns(entity any, rows sql.Rows) error
//}

// 包方法扩展性太差(能不用就不用)
//func UnSafeSetColumns(entity any, rows sql.Rows) error {
//
//}
//
//func ReflectSetColumns(entity any, rows sql.Rows) error {
//
//}
