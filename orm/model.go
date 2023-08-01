package orm

import (
	"my-frame/orm/internal/errs"
	"reflect"
	"unicode"
)

type model struct {
	tableName string
	// 表名
	fields map[string]*field
}

type field struct {
	// 列名
	colName string
}

func parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	// 限制只能用一级指针
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fields:    fieldMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) { // 判断是否是大写
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
}

/*

	元数据:
		ORM 框架需要解析模型以获得模型的元数据，这些元数据将被用于构建 SQL、执行校验，以及用于处理结果集。
		模型:一般是指对应到数据库表的 Go结构体定义，也被称为 Schema、Table 等
*/