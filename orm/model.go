package orm

import (
	"my-frame/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagKeyColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Registry(val any, opts ...ModelOpt) (*Model, error)
}

type Model struct {
	// 表名
	TableName string

	// 字段名到字段的映射
	FieldMap map[string]*Field

	// 列名到字段定义的映射
	ColumnMap map[string]*Field
}

// ModelOpt option模式(变种)
type ModelOpt func(m *Model) error

//var models = map[reflect.Type]*Model{}

// defultRegistry 全局默认的registry
//var defultRegistry = &registry{
//	models: map[reflect.Type]*Model{},
//}

// registry 代表的是元数据的注册中心
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	//models map[reflect.Type]*Model

	models sync.Map // 性能好一点但是可能会有覆盖
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val) // 会有重复解析的问题  但是i很轻微  刚启动的时候 可能会重复解析model 一次两次三次
	if err != nil {
		return nil, err
	}
	//r.models.Store(Typ, m)
	return m.(*Model), nil

}

// double check 写法
//	func (r *registry) get1(val any) (*Model, error) {
//		Typ := reflect.TypeOf(val)
//
//		r.lock.RLock()
//		m, ok := r.models[Typ]
//		r.lock.RUnlock()
//		if ok {
//			return m, nil
//		}
//
//		r.lock.Lock()
//		defer r.lock.Unlock()
//		m, ok = r.models[Typ]
//		if ok {
//			return m, nil
//		}
//
//		m, err := r.Register(val)
//		if err != nil {
//			return nil, err
//		}
//		r.models[Typ] = m
//		return m, nil
//	}

type Field struct {
	// 字段名
	GoName string

	// 列名
	ColName string

	// 代表的是字段的类型
	Typ reflect.Type

	// 字段相对于结构体本身的偏移量
	Offset uintptr
}

// Register 限制只能用一级指针
func (r *registry) Register(entity any, opts ...ModelOpt) (*Model, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	eleTyp := typ.Elem()
	numField := eleTyp.NumField()
	fieldMap := make(map[string]*Field, numField)
	columnMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := eleTyp.Field(i)
		pair, err := r.parseTab(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagKeyColumn]
		if colName == "" {
			// 用户没有设置
			colName = underscoreName(fd.Name)

		}
		fdMeta := &Field{

			GoName: fd.Name,

			ColName: colName,
			// 字段类型
			Typ:    fd.Type,
			Offset: fd.Offset,
		}

		fieldMap[fd.Name] = fdMeta
		columnMap[colName] = fdMeta
	}

	// 自定义表名
	var tableName string
	if tbl, ok := entity.(TableNane); ok {
		tableName = tbl.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(eleTyp.Name())
	}

	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)
	return res, nil
}

func ModelWithTableName(tableName string) ModelOpt {
	return func(m *Model) error {
		m.TableName = tableName
		//if TableName == "" {
		//	return err
		//}
		return nil
	}
}

func ModelWithColumnName(field, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = colName
		return nil
	}
}

type User struct {
	ID uint64 `orm:"column=id,xx=bb"`
}

func (r *registry) parseTab(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil

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

/*
	registry-并发安全思路
		并发问题解决的思路有两种想办法去除掉并发读写的场景，但是可以保留并发读，因为只读永远都是并发安全的

		用并发工具保护起来

		其实在Web 框架里面我们已经用过第一种思路了:
			要求服务器在启动之前一定要先注册好路由。

			我们可以考虑要求用户一定要提前注册好 Model。

			性能苛刻的场景下，第一种做法是比较好的选择，相当于牺牲了开发体验来换取高性能。

*/
