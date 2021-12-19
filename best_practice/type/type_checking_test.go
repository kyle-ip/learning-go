package _type

import (
	"fmt"
	"reflect"
	"testing"
)

// 类型检查：Type Assert 和 Reflection。

// ========== Type Assert ==========

// Container is a generic container, accepting anything.
type Container []interface{}

// Put adds an element to the container.
func (c *Container) Put(elem interface{}) {
	*c = append(*c, elem)
}

// Get gets an element from the container.
func (c *Container) Get() interface{} {
	elem := (*c)[0]
	*c = (*c)[1:]
	return elem
}

func TestTypeAssert(t *testing.T) {
	intContainer := &Container{}
	intContainer.Put(7)
	intContainer.Put(42)

	// 因为类型是 interface{} ，所以还要做转型，只有转型成功，才能进行后续操作
	// assert that the actual type is int
	elem, ok := intContainer.Get().(int)
	if !ok {
		fmt.Println("Unable to read an int from intContainer")
	}

	fmt.Printf("assertExample: %d (%T)\n", elem, elem)
}

// ========== Reflection ==========

type ReflectContainer struct {
	s reflect.Value
}

func NewContainer(t reflect.Type, size int) *ReflectContainer {
	if size <= 0 {
		size = 64
	}
	// 根据参数的类型初始化一个 Slice。
	return &ReflectContainer{s: reflect.MakeSlice(reflect.SliceOf(t), 0, size)}
}

func (c *ReflectContainer) Put(val interface{}) error {
	// 检查 val 是否和 Slice 的类型一致。
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		return fmt.Errorf("cannot put a %T into a slice of %s", val, c.s.Type().Elem())
	}
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
	return nil
}

func (c *ReflectContainer) Get(refval interface{}) error {
	// 需要用入参的方式，因为无法返回 reflect.Value 或 interface{}，否则还要做 Type Assert。
	if reflect.ValueOf(refval).Kind() != reflect.Ptr || reflect.ValueOf(refval).Elem().Type() != c.s.Type().Elem() {
		return fmt.Errorf("needs *%s but got %T", c.s.Type().Elem(), refval)
	}
	reflect.ValueOf(refval).Elem().Set(c.s.Index(0))
	c.s = c.s.Slice(1, c.s.Len())
	return nil
}

func TestTypeCheckByReflection(t *testing.T) {
	f1, f2 := 3.1415926, 1.41421356237

	c := NewContainer(reflect.TypeOf(f1), 16)

	if err := c.Put(f1); err != nil {
		panic(err)
	}
	if err := c.Put(f2); err != nil {
		panic(err)
	}
	g := 0.0
	if err := c.Get(&g); err != nil {
		panic(err)
	}
	fmt.Printf("%v (%T)\n", g, g) //3.1415926 (float64)
	fmt.Println(c.s.Index(0))     //1.4142135623
}
