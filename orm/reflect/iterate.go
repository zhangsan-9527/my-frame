package reflect

import "reflect"

func IterateArrayOrSlice(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		ele := val.Index(i)
		res = append(res, ele.Interface())
	}
	return res, nil
}

// keys, values, error
func IterateMap(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	resKeys := make([]any, 0, val.Len())
	resValues := make([]any, 0, val.Len())

	itr := val.MapRange()
	for itr.Next() {
		resKeys = append(resKeys, itr.Key().Interface())
		resValues = append(resValues, itr.Value().Interface())
	}
	//keys := val.MapKeys()
	//for _, key := range keys {
	//	v := val.MapIndex(key)
	//	resKeys = append(resKeys, key.Interface())
	//	resValues = append(resValues, v.Interface())
	//}
	return resKeys, resValues, nil
}

/*
	Go反射编程小技巧
		读写值，使用 reflect.Value
		读取类型信息，使用 reflect.Type
		时刻注意你现在操作的类型是不是指针。指针和指针指向的对象在反射层面上是两个东西
		大多数情况，指针类型对应的 reflect.Type 毫无用处。我们操作的都是指针指向的那个类型
		没有足够的测试就不要用反射，因为反射 API充斥着 panic
		切片和数组在反射上也是两个东西
		方法类型的字段和方法，在反射上也是两个不同的东西
*/

/*
	Go 反射面试要点
		Go反射面得很少，因为 Go 反射本身是一个写代码用的，理论上的东西不太多
		什么是反射?
			反射可以看做是对对象和对类型的描述，而我们可以通过反射来间接操作对象
		反射的使用场景?
			一大堆，基本上任何高级框架都会用反射，ORM 是一个典型例子Beego 的 controller 模式的 Web 框架也利用了反射。
		能不能通过反射修改方法?
			不能。为什么不能? Go runtime 没暴露接口。
		什么样的字段可以被反射修改?
			有一个方法 CanSet 可以判断，简单来说就是addressable.

*/
