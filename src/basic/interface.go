package basic

import "fmt"

// ========== 使用成员方法封装 ==========

type Person struct {
	Name   string
	Sexual string
	Age    int
}

// Print 成员方法（Receiver，其中 p 相当于其他语言的 this）
func (p *Person) Print() {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n", p.Name, p.Sexual, p.Age)
}

// Print 函数
func Print(p *Person) {
	fmt.Printf("Name=%s, Sexual=%s, Age=%d\n", p.Name, p.Sexual, p.Age)
}

func Example4() {
	var p = Person{
		Name:   "Kylo Yip",
		Sexual: "Male",
		Age:    26,
	}

	Print(&p)
	p.Print()
}

// ========== 面向接口而非实现编程 ==========

type Country struct {
	Name string
}

type City struct {
	Name string
}

type Stringable interface {
	ToString() string
}

func (c Country) ToString() string {
	return "Country = " + c.Name
}
func (c City) ToString() string {
	return "City = " + c.Name
}

func PrintStr(p Stringable) {
	fmt.Println(p.ToString())
}

func Example5() {
	// 使用 Stringable 接口，Country 和 City 分别实现其 ToString 方法。
	// 用接口把业务类型 Country、City 和控制逻辑 Print() 解耦。
	// 实现 Stringable 接口的结构，都可以传递到 PrintStr() 使用。
	// 类似的做法在 io.Read 和 ioutil.ReadAll 都可以看到：
	// io.Read 是接口，实现它的 Read(p []byte) (n int, err error) 接口方法，就可以被 ioutil.ReadAll 方法使用。
	d1 := Country{"USA"}
	d2 := City{"Los Angeles"}
	PrintStr(d1)
	PrintStr(d2)
}

// ========== 接口完整性检查 ==========

type Shape interface {
	Sides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Sides() int {
	return 4
}

// Square 没有实现 Shape 接口的 Area 方法，但程序仍然能跑通。
// 如果希望强制实现接口的所有方法，可以声明一个 _ 变量把 nil 的空指针从 Square 转成 Shape，使得如果没有实现完相关的接口方法，编译器就会报错：
// cannot use (*Square)(nil) (type *Square) as type Shape in assignment: *Square does not implement Shape (missing Area method)

//var _ Shape = (*Square)(nil)
//
//func (s *Square) Area() int {
//    panic("implement me")
//}

func Example6() {
	// Go 语言的编译器没有严格检查一个结构是否实现了某接口所有的方法。
	s := Square{len: 5}
	fmt.Printf("%d\n", s.Sides())
}
