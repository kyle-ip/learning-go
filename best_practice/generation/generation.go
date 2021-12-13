package basic

import "fmt"

// Go 代码生成主要用来解决编程泛型的问题。
// 泛型编程：静态类型语言有类型，相关的算法或是对数据处理的程序会因为类型不同而需要复制一份，这会导致数据类型和算法功能耦合。
// 代码生成：
// - 一个函数模板，在里面设置好相应的占位符。
// - 一个脚本，用于按规则来替换文本并生成新的代码。
// - 一行注释代码。

// 生成包名 gen，类型是 uint32，目标文件名以 container 为后缀。

//go:generate ./gen.sh ./template/container.tmp.go gen uint32 container
func generateUint32Example() {
	var u uint32 = 42
	c := NewUint32Container()
	c.Put(u)
	v := c.Get()
	fmt.Printf("generateExample: %d (%T)\n", v, v)
}

// 生成包名 gen，类型是 string，目标文件名是以 container 为后缀。

//go:generate ./gen.sh ./template/container.tmp.go gen string container
func generateStringExample() {
	var s = "Hello"
	c := NewStringContainer()
	c.Put(s)
	v := c.Get()
	fmt.Printf("generateExample: %s (%T)\n", v, v)
}
