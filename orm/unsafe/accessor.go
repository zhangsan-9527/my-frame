package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()

	fields := make(map[string]FieldMeta, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{
			Offset: fd.Offset,
			typ:    fd.Type,
		}
	}
	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields: fields,
		//address: val.UnsafeAddr(), // 为什么 不用 因为这个位置不是稳定的
		address: val.UnsafePointer(),
	}
}

/*
	注意点: unsafe 操作的是内存，本质上是对象的起始地址。
		读:*(*T)(ptr)，T 是目标类型，如果类型不知道，只能拿到反射的 Type，那么可以用reflect.NewAt(typ, ptr).Elem().

		写: *(*T)(ptr) = T，T 是目标类型

		ptr 是字段偏移量:
			ptr = 结构体起始地址 +字段偏移量

*/

func (a *UnsafeAccessor) Field(field string) (any, error) {
	// a.address 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}

	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 如果知道类型, 就这么读
	//return *(*int)(fdAddress), nil

	// 不知道确切类型 NewAt 将已知地址转成对象 New出来是对象指针
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil

}

func (a *UnsafeAccessor) SetField(field string, val any) error {
	// a.address 起始地址 + 字段偏移量
	fd, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}

	// 字段起始地址
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 知道确切类型
	//*(*int)(fdAddress) = val.(int)

	// 不知道确切类型
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

type FieldMeta struct {
	Offset uintptr
	typ    reflect.Type
}
