package orm

import (
	"my-frame/orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type model struct {
	tableName string
	// 表名
	fields map[string]*field
}

//var models = map[reflect.Type]*model{}

// defultRegistry 全局默认的registry
//var defultRegistry = &registry{
//	models: map[reflect.Type]*model{},
//}

// registry 代表的是元数据的注册中心
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	//models map[reflect.Type]*model

	models sync.Map // 性能好一点但是可能会有覆盖
}

func newRegistry() *registry {
	return &registry{}
}

func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.parseModel(val) // 会有重复解析的问题  但是i很轻微  刚启动的时候 可能会重复解析model 一次两次三次
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), nil

}

// double check 写法
//
//	func (r *registry) get1(val any) (*model, error) {
//		typ := reflect.TypeOf(val)
//
//		r.lock.RLock()
//		m, ok := r.models[typ]
//		r.lock.RUnlock()
//		if ok {
//			return m, nil
//		}
//
//		r.lock.Lock()
//		defer r.lock.Unlock()
//		m, ok = r.models[typ]
//		if ok {
//			return m, nil
//		}
//
//		m, err := r.parseModel(val)
//		if err != nil {
//			return nil, err
//		}
//		r.models[typ] = m
//		return m, nil
//	}
type field struct {
	// 列名
	colName string
}

func (r *registry) parseModel(entity any) (*model, error) {
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

/*
	registry-并发安全思路
		并发问题解决的思路有两种想办法去除掉并发读写的场景，但是可以保留并发读，因为只读永远都是并发安全的

		用并发工具保护起来

		其实在Web 框架里面我们已经用过第一种思路了:
			要求服务器在启动之前一定要先注册好路由。

			我们可以考虑要求用户一定要提前注册好 model。

			性能苛刻的场景下，第一种做法是比较好的选择，相当于牺牲了开发体验来换取高性能。

*/
