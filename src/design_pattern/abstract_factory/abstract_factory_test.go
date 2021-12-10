package abstract_factory

import "fmt"

type Person interface {
	Greet()
}

type person struct {
	name string
	age  int
}

func (p person) Greet() {
	fmt.Printf("Hi! My name is %s", p.name)
}

// NewPerson returns an interface, and not the person struct itself
func NewPerson(name string, age int) Person {
	// 由于返回接口，可以实现多个工厂函数，来返回不同的接口实现。
	// 比如 Mock 一个简单的客户端实现，可避免调用真实外部接口可能带来的失败。

	// 在实际开发中建议返回非指针的实例，因为目的主要是想通过创建实例，调用其提供的方法，
	// 如果需要对实例做更改，可以实现 SetXXX 方法。
	// 通过返回非指针的实例可确保实例的属性避免属性被意外/任意修改。
	return person{
		name: name,
		age:  age,
	}
}
