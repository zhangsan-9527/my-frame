package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

/*
前面我们使用了 unsafe.Pointer 和 uintptr这两者都代表指针，那么有什么区别?

	unsafe.Pointer: 是 Go 层面的指针，GC会维护 unsafe.Pointer 的值

	uintptr:直接就是一个数字，代表的是一个内存地址

*/

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

type FieldMeta struct {
	Offset uintptr // GC前后会变化的  表示偏移量(相对的量)
	typ    reflect.Type
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

/*
	unsafe 面试要点
		uintptr 和 unsafe.Pointer 的区别:
			前者代表的是一个具体的地址，后者代表的是一个逻辑上的指针。后者在 GC 等情况下，go runtime 会帮你调整，使其永远指向真实存放对象的地址。

		Go 对象是怎么对齐的?
			按照字长。有些比较恶心的面试官可能要你手动演示如何对齐，或者写一个对象问你怎么计算对象的大小。

		怎么计算对象地址?
			对象的起始地址是通过反射来获取，对象内部字段的地址是通过起始地址 + 字段偏移量来计算。

		unsafe 为什么比反射高效?
			可以简单认为反射帮我们封装了很多 unsafe 的操作，所以我们直接使用unsafe 绕开了这种封装的开销。有点像是我们不用ORM 框架，而是直接自己写 SQL 执行查询。
*/
