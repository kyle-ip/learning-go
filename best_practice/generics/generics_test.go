package generics

// Go 泛型在 1.18 落地前只能用 interface{} + reflect 来完成。
// interface{} 可以理解为 C 中的 void*、Java 中的 Object ，reflect 是 Go 的反射机制包，作用是在运行时检查类型。
// Go 反射：https://golang.org/pkg/reflect

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// ========== 泛型的 Map Reduce ==========

func Map(data interface{}, fn interface{}) []interface{} {
	// 通过 reflect.ValueOf() 获得 interface{} 的值，其中一个是数据 vdata，另一个是函数 vfn。
	vfn := reflect.ValueOf(fn)
	vdata := reflect.ValueOf(data)

	// 通过 vfn.Call() 方法调用函数，通过 []refelct.Value{vdata.Index(i)} 获得数据。
	result := make([]interface{}, vdata.Len())
	for i := 0; i < vdata.Len(); i++ {
		result[i] = vfn.Call([]reflect.Value{vdata.Index(i)})[0].Interface()
	}
	return result
}

func TestGenerics(t *testing.T) {
	// 由于是泛型，Map 可以接收不同类型参数的 data 和 func。

	square := func(x int) int {
		return x * x
	}
	nums := []int{1, 2, 3, 4}

	squared_arr := Map(nums, square)
	fmt.Println(squared_arr)
	// [1 4 9 16]

	upcase := func(s string) string {
		return strings.ToUpper(s)
	}
	strs := []string{"Kylo", "Yip", ""}
	upstrs := Map(strs, upcase)

	fmt.Println(upstrs)
	// [HAO CHEN MEGAEASE]

	// 但由于反射是运行时完成，如果类型出问题的话，就会有运行时的错误。
	// 比如以下代码可以编译通过，但运行会报错（panic）。
	x := Map(5, 5)
	fmt.Println(x)
}

// ========== 带类型检查的泛型 ==========

// 直接使用 interface{} 属于过度泛型，因此需要自己实现类型检查。
// Transorm 即 Map 函数（避免与数据结构 Map 混淆），会返回全新的数组。而 TransformInPlace 则是就地完成。

func Transform(slice, function interface{}) interface{} {
	return transform(slice, function, false)
}

func TransformInPlace(slice, function interface{}) interface{} {
	return transform(slice, function, true)
}

func transform(slice, function interface{}, inPlace bool) interface{} {

	// check the `slice` type is Slice
	sliceInType := reflect.ValueOf(slice)

	// Kind 方法判断通过反射取得的类型是否匹配。
	if sliceInType.Kind() != reflect.Slice {
		panic("transform: not slice")
	}

	// check the function signature
	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()

	// 检查函数的参数和返回类型，NumIn() 检查入参，NumOut() 检查返回值。
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("transform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
	}

	sliceOutType := sliceInType
	if !inPlace {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}
	return sliceOutType.Interface()

}

func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {

	// Check it is a function
	// Kind 方法判断通过反射取得的类型是否匹配。
	if fn.Kind() != reflect.Func {
		return false
	}
	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}
	// In() - returns the type of a function type's i'th input parameter.
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - returns the type of a function type's i'th output parameter.
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}
