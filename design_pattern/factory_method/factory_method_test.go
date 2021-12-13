package factory_method

type Person struct {
	name string
	age  int
}

// NewPersonFactory 带默认年龄的工厂
func NewPersonFactory(age int) func(name string) Person {
	// 通过实现工厂函数来创建多种工厂，将对象创建从由一个对象负责所有具体类的实例化，
	// 变成由一群子类来负责对具体类的实例化，从而将过程解耦。
	return func(name string) Person {
		return Person{
			name: name,
			age:  age,
		}
	}
}

func NewPerson() {
	newBaby := NewPersonFactory(1)
	_ = newBaby("john")

	newTeenager := NewPersonFactory(16)
	_ = newTeenager("jill")
}
