package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() {
	fmt.Printf("Hi! My name is %s", p.Name)
}

func NewPerson(name string, age int) *Person {
	// 通过传入参数，创建目标对象并返回。
	// 如果要增加一种产品，就要在工厂中修改创建产品的函数，导致耦合性过高。
	return &Person{
		Name: name,
		Age:  age,
	}
}
