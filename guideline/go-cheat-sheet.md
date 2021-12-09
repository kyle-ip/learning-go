---
title: Go Cheat Sheet
date: 2021-06-13 14:42:39
tags:
  - Go
  - 技术
abstract: 参考：https://github.com/a8m/golang-cheat-sheet
description: 参考：https://github.com/a8m/golang-cheat-sheet
---

>   参考：[An overview of Go syntax and features.](https://github.com/a8m/golang-cheat-sheet)
>
>   文中大量的代码都是取自 [A Tour of Go](http://tour.golang.org/)（非常好的 Go 入门教程）。
>
>   如果是刚入门，非常建议先从这份教程上手。

# Go Cheat Sheet

简而言之，就是 ...：

* 命令式语言

* 静态类型

* 类似 C 语言的语法标记（更少的括号，没有分号）和 Oberon-2 的结构

* 编译为本地代码（区别于 JVM）

* 没有类，但有带方法的结构体

* 接口

* 没有实现继承，但有 [类型嵌入](http://golang.org/doc/effective%5Fgo.html#embedding)

* 函数是一等公民

* 函数可以返回多个值

* 闭包

* 指针（不支持算术运算）

* 内置并发原语：Goroutines 和 Channels 

## Hello World

文件：`hello.go`:

```go
package main
// package 名必须为 main，目录名一般与 package 名一致，文件名不作要求。

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Hello Go")
    os.Exit(0)
    // os.Args
}
```
运行：`go run hello.go`

编译：`go build hello.go`

## 测试

-   源码文件以 `_test` 结尾。

-   测试方法名以 `Test` 开头：`func TestXxx(t *testing.T) {...}`

```go
const (
    Monday = 1 + iota
    Tuesday
    Wednesday
)

const (
    Readable = 1 << iota
    Writable
    Executable
)

func TestConstantTry(t *testing.T) {
    t.Log(Monday, Tuesday)
}

func TestConstantTry1(t *testing.T) {
    a := 1 // 0001
    t.Log(a&Readable == Readable, a&Writable == Writable, a&Executable == Executable)
}
```

测试：`go test test.go`

### 单元测试

Go 内置单元测试框架：

-   Fail, Error: 该测试最终失败，但后续代码继续执行，且其他测试继续执行。

-   FailNow, Fatal: 该测试失败且马上中止，其他测试继续执行。

-   测试覆盖率：`go trst -v -cover`

-   断言：https://github.com/stretchr/testfy

```go
import (
    "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
)

// 待测试函数
func square(op int) int {
    return op * op
}

func TestSquare(t *testing.T) {
    inputs := [...]int{1, 2, 3}
    expected := [...]int{1, 4, 9}
    for i := 0; i < len(inputs); i++ {
        ret := square(inputs[i])
        if ret != expected[i] {
            t.Errorf("input is %d, the expected is %d, the actual %d", inputs[i], expected[i], ret)
        }
    }
}

func TestErrorInCode(t *testing.T) {
    fmt.Println("Start")
    t.Error("Error")
    fmt.Println("End")
}

func TestFailInCode(t *testing.T) {
    fmt.Println("Start")
    t.Fatal("Error")
    fmt.Println("End")
}

func TestSquareWithAssert(t *testing.T) {
    inputs := [...]int{1, 2, 3}
    expected := [...]int{1, 4, 9}
    for i := 0; i < len(inputs); i++ {
        ret := square(inputs[i])
        assert.Equal(t, expected[i], ret)
    }
}
```

### 基准测试

基准测试：`go test -bench=. -benchmem`（`-bench=<相关benchmark测试>`）

Windows 下使⽤ `go test` 命令⾏时应写为 `go test -bench="."`。

```go
func TestConcatStringByAdd(t *testing.T) {
    assert := assert.New(t)
    elems := []string{"1", "2", "3", "4", "5"}
    ret := ""
    for _, elem := range elems {
        ret += elem
    }
    assert.Equal("12345", ret)
}

func TestConcatStringByBytesBuffer(t *testing.T) {
    assert := assert.New(t)
    var buf bytes.Buffer
    elems := []string{"1", "2", "3", "4", "5"}
    for _, elem := range elems {
        buf.WriteString(elem)
    }
    assert.Equal("12345", buf.String())
}

func BenchmarkConcatStringByAdd(b *testing.B) {
    elems := []string{"1", "2", "3", "4", "5"}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ret := ""
        for _, elem := range elems {
            ret += elem
        }
    }
    b.StopTimer()
}

func BenchmarkConcatStringByBytesBuffer(b *testing.B) {
    elems := []string{"1", "2", "3", "4", "5"}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var buf bytes.Buffer
        for _, elem := range elems {
            buf.WriteString(elem)
        }
    }
    b.StopTimer()
}
```

### BDD (Behavior Driven Development)

https://github.com/smartystreets/goconvey

安装：`go get -u github.com/smartystreets/goconvey/convey`

启动：`$GOPATH/bin/goconvey`

## 操作符

不支持前置的自增（++）自减（--）。

### 算术

| 操作符 | 描述 |
|--------|-----------|
|`+`| 加 |
|`-`| 减 |
|`*`| 乘 |
|`/`| 除 |
|`%`| 求余 |
|`&`| 按位与 |
|`\|`| 按位或 |
|`^`| 按位异或 |
|`&^`| 按位清零（且非） |
|`<<`| 左移 |
|`>>`| 右移 |

### 比较
| 操作符  | 描述 |
|--------|-----------|
|`==`| 等于 |
|`!=`| 不等于 |
|`<`| 小于 |
|`<=`| 小于等于 |
|`>`| 大于 |
|`>=`| 大于等于 |

### 逻辑 
| 操作符 | 描述 |
|--------|-----------|
|`&&`| 逻辑且 |
|`||`| 逻辑或 |
|`!`| 逻辑非 |

### 其他

| 操作符 | 描述 |
|--------|-----------|
|`&`| 取址 / 创建指针 |
|`*`| 取值 |
|`<-`| 发送 / 接收 (见下文“Channels”)|

## 声明

类型在标识符后面，且赋值时可自动推断类型。

```go
var foo int                     // 无初始化的声明
var foo int = 42                // 带初始化的声明
var foo, bar int = 42, 1302     // 一次声明并初始化多个变量
var foo = 42                    // 类型省略（推断）
foo := 42                       // 初始化且赋值，只能用在 func 主体中，省略 var 关键字，类型是隐式推断的
const constant = "This is a constant"
```

交换两个变量的值（在同一个语句中可对多个变量同时赋值）。

```go
a := 1
b := 2
a, b = b, a
```

快速设置连续的值。

```go
// iota 可用于递增数字，从 0 开始 
const (
    Monday = 1 + iota
    Tuesday
    Wednesday
)

const (
    Readable = 1 << iota
    Writable
    Executable
)

func TestConstantTry(t *testing.T) {
    t.Log(Monday, Tuesday)
}

func TestConstantTry1(t *testing.T) {
    a := 1 //0001
    t.Log(a&Readable == Readable, a&Writable == Writable, a&Executable == Executable)
}
```

## 函数

-   可带多个返回值。

-   所有参数都是值传递：slice，map，channel 等（复制指针，指向同一块内存）。

-   可作为变量的值。

-   可作为参数和返回值。

```go
// 一个简单的函数
func functionName() {}

// function with parameters (again, types go after identifiers)
func functionName(param1 string, param2 int) {}

// 带参数的函数（类型在标识符之后） 
func functionName(param1, param2 int) {}

// 返回类型声明 
func functionName() int {
    return 42
}

// 可以一次返回多个值 
func returnMulti() (int, string) {
    return 42, "foobar"
}
var x, str = returnMulti()

// 简单通过 return 返回多个命名结果 
func returnMulti2() (n int, s string) {
    n = 42
    s = "foobar"
    // n and s will be returned
    return
}
var x, str = returnMulti2()

```

### 作为值的函数与闭包
```go
func main() {
    // 把一个函数赋给 add
    add := func(a, b int) int {
        return a + b
    }
    // 通过 add 调用函数
    fmt.Println(add(3, 4))
}

// 闭包，词法作用域：函数可以访问定义函数作用域内的值 
func scope() func() int{
    outer_var := 2
    foo := func() int { return outer_var}
    return foo
}

func another_scope() func() int{
    // 无法通过编译，因为在这个作用域内没有定义 outer_var 和 foo 
    outer_var = 444
    return foo
}


// 闭包
func outer() (func() int, int) {
    outer_var := 2
    inner := func() int {
        outer_var += 99             // 外部作用域的 outer_var 被修改
        return outer_var
    }
    inner()
    return inner, outer_var         // 返回内部函数和变异的外部变量
}
```

### 带可变参数的函数

```go
func main() {
    fmt.Println(adder(1, 2, 3))         // 6
    fmt.Println(adder(9, 9))            // 18

    nums := []int{10, 20, 30}
    fmt.Println(adder(nums...))         // 60
}

// 在最后一个参数的类型名称之前加上 ... ，可以指示它接收零个或多个参数。
// 除了可以传递任意数量的参数之外，该函数的调用方式与任何其他函数一样。 
func adder(args ...int) int {
    total := 0
    for _, v := range args {            // 迭代访问参数
        total += v
    }
    return total
}
```

### Defer 函数

延迟执行，且不受 panic 影响。

常用于解锁，释放资源等。

```go
func TestDefer(t *testing.T) {
    defer func() {
        t.log("Clear resources")
    }()
    t.log("Started")
    panic("Fatal error")
}
```

## 数据类型

Go 是静态强类型语言，不支持隐式类型转换，而别名和原有类型不能进⾏隐式类型转换。

所有 Go 预先声明的标识符都在 [builtin](https://golang.org/pkg/builtin/) 包中定义。 

```go
bool        // 默认 false

string      // 默认 ""

int  int8  int16  int32  int64                // 默认 0
uint uint8 uint16 uint32 uint64 uintptr

byte        // uint8 的别名

rune        // alias int32 ~= a character (Unicode code point) - very Viking

float32 float64

complex64 complex128
```

```go
type MyInt int64

func TestImplicit(t *testing.T) {
    var a int32 = 1
    var b int64
    b = int64(a)
    var c MyInt
    c = MyInt(b)
    t.Log(a, b, c)
}

func TestPoint(t *testing.T) {
    a := 1
    aPtr := &a
    //aPtr = aPtr + 1
    t.Log(a, aPtr)
    t.Logf("%T %T", a, aPtr)
}

func TestString(t *testing.T) {
    var s string
    t.Log("*" + s + "*")         // 初始化零值是“”
    t.Log(len(s))
}
```

## 类型转换
```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

// 替代语法
i := 42
f := float64(i)
u := uint(f)
```

## 包

Package 是基本复用模块单元，在每个源文件顶部声明。

- 可执行文件放在 `main` 包中。

- 代码的 package 可以和所在的目录不一致，但同一目录中代码的 package 要保持一致。

- 包名 == 导入路径的最后一层（`import math/rand` => `package rand`），导入时从 src 后开始。

- 大写标识符：导出（从其他包中可见）；小写标识符：私有（在其他包中不可见）。 

- `go get [-u] github.com/stretchr/testify/assert` 可（强制）从远程获取/更新依赖。

```go
import (
    "testing"
    // go get -u github.com/easierway/concurrent_map
    cm "github.com/easierway/concurrent_map"
)

func TestConcurrentMap(t *testing.T) {
    m := cm.CreateConcurrentMap(99)
    m.Set(cm.StrKey("key"), 10)
    t.Log(m.Get(cm.StrKey("key")))
}

```

如果要把自己的代码提交到 GitHub，应注意适应 go get：直接以代码路径开始，不要有 src。

参考：https://github.com/easierway/concurrent_map

### init 方法

-   在 main 被执⾏前，所有依赖的 package 的 init ⽅法都会被执⾏。

-   不同包的 init 函数按照包导⼊的依赖关系决定执⾏顺序。

-   每个包可以有多个 init 函数。

-   包的每个源⽂件也可以有多个 init 方法。

```go
package series

func init() {
    fmt.Println("init1")
}

func init() {
    fmt.Println("init2")
}

// 在别处 import series 时，init 方法即被执行。
```

### 依赖管理

Go 未解决依赖管理问题：

- 同⼀环境下，不同项⽬无法使⽤同⼀包的不同版本。

- ⽆法管理对包的特定版本的依赖。

Vendor 路径：Go 1.5 release 后，vendor ⽬录被添加到除了 GOPATH 和 GOROOT 之外的依赖⽬录查找的解决⽅案（在 Go 1.6 前需要⼿动设置环境变量）。查找依赖包路径的解决⽅案：

-   当前包下的 vendor ⽬录。

-   向上级⽬录查找，直到找到 src 下的 vendor ⽬录。

-   在 GOPATH 下⾯查找依赖包。

-   在 GOROOT ⽬录下查找。

第三方依赖管理工具：

- [godep](https://github.com/tools/godep)

- [glide](https://github.com/Masterminds/glide)

- [dep](https://github.com/golang/dep)

```shell
mkdir module_package && cd !$
glide init
vim glide.yaml
    # ...
glide install 
```

## 控制结构

### 条件

```go
func main() {
    // Basic one
    if x > 10 {
        return x
    } else if x == 10 {
        return 10
    } else {
        return -x
    }

    // 可以在条件前放置一个语句 
    if a := b + c; a < 42 {
        return a
    } else {
        return a - 42
    }

    // 在 if 中输入断言 
    var val interface{} = "foo"
    if str, ok := val.(string); ok {
        fmt.Println(str)
    }
}
```

### 循环
```go
// 只有 `for`，没有 `while`，没有 `until` 
for i := 1; i < 10; i++ {
}
for ; i < 10;  {        // while 循环
}
for i < 10  {           // 如果只有一个条件，可以省略分号 
}
for {                   // 省略条件 ~ while (true)
}

// 在当前循环中使用 break/continue
// 在外循环上使用带标签的 break/continue 

for i := 0; i < 2; i++ {
    for j := i + 1; j < 3; j++ {
        if i == 0 {
            continue here
        }
        fmt.Println(j)
        if j == 2 {
            break
        }
    }
}

for i := 0; i < 2; i++ {
    for j := i + 1; j < 3; j++ {
        if j == 1 {
            continue
        }
        fmt.Println(j)
        if j == 2 {
            break there
        }
    }
}
```

### 分支
```go
// switch 声明
switch operatingSystem {
    case "darwin":
        fmt.Println("Mac OS Hipster")
        // case 自动中断，默认情况下没有失败（不需要 break）
    case "linux":
        fmt.Println("Linux Geek")
    default:
        // Windows, BSD, ...
        fmt.Println("Other")
}

// 与 for 和 if 一样，可以在 switch 值之前加一个赋值语句 
switch os := runtime.GOOS; os {
    case "darwin": ...
}

// 还可以在 switch case 中进行比较 
number := 42
switch {
    case number < 42:
        fmt.Println("Smaller")
    case number == 42:
        fmt.Println("Equal")
    case number > 42:
        fmt.Println("Greater")
}

// case 可以用逗号分隔的列表显示 
var char byte = '?'
switch char {
    case ' ', '?', '&', '=', '#', '+', '%':
    fmt.Println("Should escape")
}
```

## 字符串

-   string 是基本数据类型（不是引⽤或指针）。

-   string 是只读的 byte slice，`len` 函数可以求它所包含的 byte 数。

-   string 的 byte 数组可以存放任何数据。

-   常用的字符串包：[strings](https://golang.org/pkg/strings), [strcov](https://golang.org/pkg/strconv)

```go
func TestString(t *testing.T) {
    var s string
    t.Log(s) //初始化为默认零值“”
    s = "hello"
    t.Log(len(s))
    // s[1] = '3' // string是不可变的byte slice
    // s = "\xE4\xB8\xA5" // 可以存储任何二进制数据
    s = "\xE4\xBA\xBB\xFF"
    t.Log(s)
    t.Log(len(s))
    s = "中"
    t.Log(len(s)) // 是 byte 数

    c := []rune(s)
    t.Log(len(c))
    // t.Log("rune size:", unsafe.Sizeof(c[0]))
    t.Logf("中 unicode %x", c[0])
    t.Logf("中 UTF8 %x", s)
}

func TestStringToRune(t *testing.T) {
    s := "中华人民共和国"
    for _, c := range s {
        t.Logf("%[1]c %[1]x", c)
    }
}

func TestStringFn(t *testing.T) {
    s := "A,B,C"
    parts := strings.Split(s, ",")
    for _, part := range parts {
        t.Log(part)
    }
    t.Log(strings.Join(parts, "-"))
}

func TestConv(t *testing.T) {
    s := strconv.Itoa(10)
    t.Log("str" + s)
    if i, err := strconv.Atoi("10"); err == nil {
        t.Log(10 + i)
    }
}
```



## 数组, 切片, 范围

### 数组

```go
var a [10]int           // 声明一个长度为 10 的 int 数组。数组长度是类型的一部分
a[3] = 42               // 设值
i := a[3]               // 取值

// 声明和初始化
var a = [2]int{1, 2}
a := [2]int{1, 2}       // 简写
a := [...]int{1, 2}     // elipsis -> 编译器自动计算数组长度 
```

只有相同维数且含有相同元素个数的数组才可以进行比较：所有元素相等才相等。

### 切片

```go
var a []int                              // 声明一个切片 - 类似于数组，但长度未指定 
var a = []int{1,2,3,4}                   // 声明并初始化一个切片（由隐式给出的数组支持） 
a := []int{1,2,3,4}                      // shorthand
chars := []string{0:"a", 2:"c", 1: "b"}  // ["a", "b", "c"]

var b = a[lo:hi]                         // 创建一个从索引 lo 到 hi-1 的切片（数组的视图）
var b = a[1:4]                           // 从索引 1 到 3 的切片
var b = a[:3]                            // 缺少低索引意味着 0
var b = a[3:]                            // 缺少高索引意味着 len(a)
a = append(a,17,3)                       // 把值得附加到切片 a
c := append(a,b...)                      // 连接切片 a 和 b 

// 用 make 创建一个切片 
a = make([]byte, 5, 5)                   // 第一个参数是长度，第二个是容量 
a = make([]byte, 5)                      // 容量参数可选

// 从数组创建切片 
x := [3]string{"Лайка", "Белка", "Стрелка"}
s := x[:]                                // 一个引用 x 存储的切片
```

从数组中创建的切片的容量，等于数组的容量 - 切片的起始下标。

### 数组与切片操作

`len(a)` 提供数组/切片的长度。它是一个内置函数，而不是数组的属性或方法。 

```go
// 循环遍历数组/切片 
for i, e := range a {
    // i 是索引, e 是元素
}

// 如果只需要元素 e:
for _, e := range a {
    // e 是元素
}

// 如果只需要索引：
for i := range a {
}

// 在 Go 1.4 之前的版本中，如果不使用 i 和 e，则会出现编译错误。
//  Go 1.4 引入了无变量形式：
for range time.Tick(time.Second) {
    // 每秒执行一次 
}

```

## 映射

```go
m := make(map[string]int, 10)    // cap 为 10，不需要初始化 len
m["key"] = 42
fmt.Println(m["key"])

delete(m, "key")

if elem, ok := m["key"]; ok {    // 判断键“key”是否存在并取值
    // ...
}

// map 表述
var m = map[string]Vertex{
    "Bell Labs": {40.68433, -74.39967},
    "Google":    {37.42202, -122.08408},
}

// 遍历 map 的内容
for key, value := range m {
}

```

实现工厂模式：

```go
func TestMapWithFunValue(t *testing.T) {
    m := map[int]func(op int) int{}
    m[1] = func(op int) int { return op }
    m[2] = func(op int) int { return op * op }
    m[3] = func(op int) int { return op * op * op }
    t.Log(m[1](2), m[2](2), m[3](2))
}

func TestMapForSet(t *testing.T) {
    mySet := map[int]bool{}
    mySet[1] = true
    n := 3
    if mySet[n] {
        t.Logf("%d is existing", n)
    } else {
        t.Logf("%d is not existing", n)
    }
    mySet[3] = true
    t.Log(len(mySet))
    delete(mySet, 1)
    n = 1
    if mySet[n] {
        t.Logf("%d is existing", n)
    } else {
        t.Logf("%d is not existing", n)
    }
}
```

Go 的内置集合中没有 Set，但可以利用 `map[type]bool` 实现。

## 结构

可附带方法，不支持继承。
```go
// 结构是一种类型，也是字段的集合。

// 声明
type Vertex struct {
    X, Y int
}

// 创建
var v = Vertex{1, 2}
var v = Vertex{X: 1, Y: 2}              // 通过使用键定义值来创建结构 
var v = []Vertex{{1,2},{5,2},{5,5}}     // 初始化一个结构切片

// 访问成员
v.X = 4

// 声明结构的方法
func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 调用方法
v.Abs()

// 对于修改方法，参数类型为指向 Struct 的指针（见下文，这样就不会为方法调用复制整个结构的值）。 
func (v *Vertex) add(n float64) {
    v.X += n
    v.Y += n
}

```
**匿名结构**：比使用 `map[string]interface{}` 开销更低、更安全.

```go
point := struct {
    X, Y int
}{1, 2}
```

扩展（不支持访问子类的字段、方法，无法重载等）：

```go
type Pet struct {}

func (p *Pet) Speak() {
    fmt.Print("...")
}

func (p *Pet) SpeakTo(host string) {
    p.Speak()
    fmt.Println(" ", host)
}

type Dog struct {
    Pet
}

func (d *Dog) Speak() {
    fmt.Print("Wang!")
}

func TestDog(t *testing.T) {
    dog := new(Dog)
    dog.SpeakTo("Chao")
}
```

## 指针

不支持指针运算。

```go
p := Vertex{1, 2}       // p 是一个 Vertex
q := &p                 // q 是一个指向 Vertex 的指针
r := &Vertex{1, 2}      // r 是一个指向 Vertex 的指针

// 指向顶点的指针的类型：*Vertex 

var s *Vertex = new(Vertex)     // new 创建一个指向新结构体实例的指针 
```

## 接口

Go 的接口是非侵入性的，实现不依赖于接口定义，因此接口的定义可以包含在接口使用者包内。

```go
type Programmer interface {
    WriteHelloWorld() string
}

type GoProgrammer struct {
    name string
}

// 如果类型实现了所有必需的方法，则类型隐式地满足接口
func (g *GoProgrammer) WriteHelloWorld() string {
    return "fmt.Println(\"Hello World\")"
}

// 强制实现接口的方法
// var _ Programmer = (*GoProgrammer)(nil)

func TestClient(t *testing.T) {
    var p Programmer
    p = new(GoProgrammer)
    t.Log(p.WriteHelloWorld())
}
```

多态：

```go
type Code string

type Programmer interface {
    WriteHelloWorld() Code
}

type GoProgrammer struct {
}

func (p *GoProgrammer) WriteHelloWorld() Code {
    return "fmt.Println(\"Hello World!\")"
}

type JavaProgrammer struct {
}

func (p *JavaProgrammer) WriteHelloWorld() Code {
    return "System.out.Println(\"Hello World!\")"
}

func writeFirstProgram(p Programmer) {
    fmt.Printf("%T %v\n", p, p.WriteHelloWorld())
}

func TestPolymorphism(t *testing.T) {
    goProg := &GoProgrammer{}
    javaProg := new(JavaProgrammer)
    writeFirstProgram(goProg)
    writeFirstProgram(javaProg)
}

```

其中空接口可以表示任何类型，通过断言可将空接口转换为指定类型：`v, ok := p.(int)`。

（可类比 Java 的 `Object`、C/C++ 的 `void *` 指针）

```go
func DoSomething(p interface{}) {
    switch v := p.(type) {
        case int:
            fmt.Println("Integer", v)
        case string:
            fmt.Println("String", v)
        default:
            fmt.Println("Unknow Type")
        }
}

func TestEmptyInterfaceAssertion(t *testing.T) {
    DoSomething(10)
    DoSomething("10")
}
```

## 嵌入

Go 中没有子类化，但有接口和结构嵌入。 

```go
// ReadWriter 的实现必须同时满足 Reader 和 Writer 
type ReadWriter interface {
    Reader
    Writer
}

// Server 暴露了 Logger 拥有的所有方法 
type Server struct {
    Host string
    Port int
    *log.Logger
}

// 初始化嵌入类型 
server := &Server{"localhost", 80, log.New(...)}

// 在嵌入式结构上实现的方法被传递 
server.Log(...) // calls server.Logger.Log(...)

// 嵌入类型的字段名称是其类型名称（在本例中为 Logger） 
var logger *log.Logger = server.Logger
```

## 异常

Go 没有异常处理机制，`error` 是可恢复的错误，`panic` 为不可恢复错误（退出前会执行 defer 的内容）。

可能产生错误的函数会声明一个额外的 `error` 类型返回值（https://golang.org/pkg/builtin/#error）。

其中 `error` 类型实现了 `error` 接口，有一个 `Error` 方法，并可以通过 `errors.New("xxx")` 创建错误实例。

```go
// 内置接口类型 error 是表示错误条件的常规接口，nil 值表示没有错误。 
type error interface {
    Error() string
}
```

示例：
```go
var LessThanTwoError = errors.New("n should be not less than 2")

var LargerThenHundredError = errors.New("n should be not larger than 100")

func GetFibonacci(n int) ([]int, error) {
    if n < 2 {
        return nil, LessThanTwoError
    }
    if n > 100 {
        return nil, LargerThenHundredError
    }
    fibList := []int{1, 1}

    for i := 2; /*短变量声明 := */ i < n; i++ {
        fibList = append(fibList, fibList[i-2]+fibList[i-1])
    }
    return fibList, nil
}

func GetFibonacci1(str string) {
    var (
        i    int
        err  error
        list []int
    )
    if i, err = strconv.Atoi(str); err != nil {
        fmt.Println("Error", err)
        return
    }
    if list, err = GetFibonacci(i); err != nil {
        fmt.Println("Error", err)
        return
    }
    fmt.Println(list)
}

func TestGetFibonacci(t *testing.T) {
    if v, err := GetFibonacci(1); err != nil {
        if err == LessThanTwoError {
            fmt.Println("It is less.")
        }
        t.Error(err)
    } else {
        t.Log(v)
    }

}
```

异常处理（类似 Java、C++ 的 `try ... catch ...`）：

```go
func TestPanicVxExit(t *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println("recovered from ", err)
        }
    }()
    fmt.Println("Start")
    panic(errors.New("Something wrong!"))
    // os.Exit(-1)
    // fmt.Println("End")
}
```

recover 要慎用：

-   避免形成僵尸进程，导致 health check 失效。

-   恢复不确定性错误时，让其 crash 是最好的办法。

第三方库 error：https://github.com/pkg/errors

## 协程：Goroutines

协程是由 Go 管理的轻量级线程（stack 初始化为 2K，远远小于 Java Thread Stack 的 1M）。
`go f(a, b)` 启动一个运行 `f` 函数的 goroutine。 

```go
// 可以作为 goroutine 启动的函数 
func doStuff(s string) {
}

func main() {
    // 在 goroutine 中使用命名函数 
    go doStuff("foobar")

    // 在 goroutine 中使用匿名内部函数 
    go func (x int) {
    }(42)
}
```

### 共享内存并发机制

```go
// 无同步，线程不安全
func TestCounterThreadUnsafe(t *testing.T) {
    counter := 0
    for i := 0; i < 5000; i++ {
        go func() {
            counter++
        }()
    }
    time.Sleep(1 * time.Second)
    t.Logf("counter = %d", counter)
}

// 线程安全
func TestCounterThreadSafe(t *testing.T) {
    var mut sync.Mutex
    counter := 0
    for i := 0; i < 5000; i++ {
        go func() {
            // 解锁
            defer func() {
                mut.Unlock()
            }()
            // 加锁
            mut.Lock()
            counter++
        }()
    }
    time.Sleep(1 * time.Second)
    t.Logf("counter = %d", counter)

}

// 等待组
func TestCounterWaitGroup(t *testing.T) {
    var mut sync.Mutex
    var wg sync.WaitGroup
    counter := 0
    for i := 0; i < 5000; i++ {
        // 计数器 +1
        wg.Add(1)
        go func() {
            defer func() {
                mut.Unlock()
            }()
            mut.Lock()
            counter++
            // 计数器 -1
            wg.Done()
        }()
    }
    // 等待计数器清零
    wg.Wait()
    t.Logf("counter = %d", counter)
}
```

## 通道：Channels

一般情况下通过 Channel 通信，发送者和接收者都是阻塞、相互等待的。如果使用 Buffer Channel，则 Channel 中会缓存发送者数条消息（不超过容量），接收者随时都可以从中获取。

```go
ch := make(chan int)    // 创建一个 int 类型的 channel
ch <- 42                // 向 ch 发送一个值。 
v := <-ch               // 从 ch 接收一个值 

// 非缓冲 channel 块，当没有可用值时读取块，写入块直到可读取

// 创建一个缓冲 channel。如果写入的值小于 <buffer size>，则写入缓冲 channel 不会阻塞。 
ch := make(chan int, 100)

close(ch)               // 关闭 channel（只有发送者可关闭） 

// 从 channel 读取并测试它是否已关闭（ok 为 false 则通道已关闭 ）
v, ok := <-ch

// 从 channel 读取，直到它关闭 
for i := range ch {
    fmt.Println(i)
}

// 多 channel 阻塞分支，如果其中有一个就绪，则执行相应的 case 
func doStuff(channelOut, channelIn chan int) {
    select {
        case channelOut <- 42:
            fmt.Println("We could write to channelOut!")
        case x := <- channelIn:
            fmt.Println("We could read from channelIn")
        case <-time.After(time.Second * 1):
            fmt.Println("timeout")
    }
}
```

### 异步任务

```go
func service() string {
    time.Sleep(time.Millisecond * 50)
    return "Done"
}

func otherTask() {
    fmt.Println("working on something else")
    time.Sleep(time.Millisecond * 100)
    fmt.Println("Task is done.")
}

func TestService(t *testing.T) {
    fmt.Println(service())
    otherTask()
}

func AsyncService() chan string {
    // 创建一个 string 类型的 channel。
    retCh := make(chan string, 1)
    // 创建协程，执行异步任务（非阻塞）。
    go func() {
        // 协程阻塞，直到执行完成、把结果写入 channel。
        ret := service()
        fmt.Println("returned result.")
        retCh <- ret
        fmt.Println("service exited.")
    }()
    // 返回 channel。
    return retCh
}

func TestAsyncService(t *testing.T) {
    // 执行异步任务，返回一个 channel（类似 Java 的 Future）。
    retCh := AsyncService()
    
    // 执行其他任务。
    otherTask()
    
    // 从 channel 读取结果。
    fmt.Println(<-retCh)
    time.Sleep(time.Second * 1)
}

```

### 多通道选择

只要任何一个通道消息就绪，即执行对应的 case。

```go
select {
    case ret := <-retCh1:
        t.Logf("result %s", ret)
    case ret :=<-retCh2:
        t.Logf("result %s", ret)
    default:
        t.Error(“No one returned”)
}
```

超时控制：

```go
select {
    case ret := <-retCh:
        t.Logf("result %s", ret)
    // 阻塞等待 time.Second * 1。 
    case <-time.After(time.Second * 1):
        t.Error("time out")
}
```

### 任务取消

通过关闭 channel 取消：

```
func cancel_2(cancelChan chan struct{}) {
    close(cancelChan)
}
```

发送取消消息：

```go
func cancel_1(cancelChan chan struct{}) {
    cancelChan <- struct{}{}
}
```

获取取消通知：

```go
func isCancelled(cancelChan chan struct{}) bool {
    select {
        case <-cancelChan:
            return true
        default:
            return false
    } 
}
```

```go
func TestCancel(t *testing.T) {
    cancelChan := make(chan struct{}, 0)
    for i := 0; i < 5; i++ {
        go func(i int, cancelCh chan struct{}) {
            for {
                if isCancelled(cancelCh) {
                    break
                }
                time.Sleep(time.Millisecond * 5)
            }
            fmt.Println(i, "Cancelled")
        }(i, cancelChan)
    }
    cancel_2(cancelChan)
    time.Sleep(time.Second * 1)
}
```



### Channel 规则

发送到 nil channel 会一直阻塞：

```go
var c chan string
c <- "Hello, World!"
// fatal error: 所有 goroutine 都处于睡眠状态 - 死锁！ 
```

从 nil channel 读取会一直阻塞：

```go
var c chan string
fmt.Println(<-c)
// fatal error: 所有 goroutine 都处于睡眠状态 - 死锁！ 
```

发送数据到已关闭的 channel 发生错误 ：

```go
var c = make(chan string, 1)
c <- "Hello, World!"
close(c)
c <- "Hello, Panic!"
// panic: 在关闭的 channel 上发送 
```

在关闭的 channel 上接收会立即返回零值：

```go
var c = make(chan int, 2)
c <- 1
c <- 2
close(c)
for i := 0; i < 3; i++ {
    fmt.Printf("%d ", <-c)
}
// 1 2 0

```

在 channel 关闭时，所有的接收者都会立即从阻塞等待中返回，且 `v, ok <-ch` 的 ok 值为 false。

```go
func dataProducer(ch chan int, wg *sync.WaitGroup) {
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        close(ch)
        wg.Done()
    }()
}

func dataReceiver(ch chan int, wg *sync.WaitGroup) {
    go func() {
        // 循环读取，直到 ok 为 false（通道关闭）
        for {
            if data, ok := <-ch; ok {
                fmt.Println(data)
            } else {
                break
            }
        }
        wg.Done()
    }()
}

func TestCloseChannel(t *testing.T) {
    var wg sync.WaitGroup
    ch := make(chan int)
    wg.Add(1)
    dataProducer(ch, &wg)
    wg.Add(1)
    dataReceiver(ch, &wg)
    // wg.Add(1)
    // dataReceiver(ch, &wg)
    wg.Wait()
}
```

这个⼴播机制常被利⽤，进⾏向多个订阅者同时发送信号，如退出信号。

## 上下文

可用于关联任务取消：

-   根 Context：通过 `context.Background ()` 创建。

-   子 Context：``context.WithCancel(parentContext)` 创建。

-   `ctx, cancel := context.WithCancel(context.Background())`。

-   当前 Context 被取消时，基于他的⼦ context 都会被取消。

-   接收取消通知 `<-ctx.Done()`。

```go
func isCancelled(ctx context.Context) bool {
    select {
        case <-ctx.Done():
            return true
        default:
            return false
    }
}

func TestCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    for i := 0; i < 5; i++ {
        go func(i int, ctx context.Context) {
            for {
                if isCancelled(ctx) {
                    break
                }
                time.Sleep(time.Millisecond * 5)
            }
            fmt.Println(i, "Cancelled")
        }(i, ctx)
    }
    cancel()
    time.Sleep(time.Second * 1)
}
```

## 常见并发执行模式

### 单例模式（懒汉式）

```go
type Singleton struct {
    data string
}

var singleInstance *Singleton

var once sync.Once

func GetSingletonObj() *Singleton {
    // 确保只执行一次。
    once.Do(func() {
        fmt.Println("Create Obj")
        singleInstance = new(Singleton)
    })
    return singleInstance
}

func TestGetSingletonObj(t *testing.T) {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            obj := GetSingletonObj()
            fmt.Printf("%X\n", unsafe.Pointer(obj))
            wg.Done()
        }()
    }
    wg.Wait()
}

```

### 任意任务完成

```go
func runTask(id int) string {
    time.Sleep(10 * time.Millisecond)
    return fmt.Sprintf("The result is from %d", id)
}

func FirstResponse() string {
    numOfRunner := 10
    // 容量为 10 的 buffer channel。
    ch := make(chan string, numOfRunner)
    // 10 个协程执行任务。
    for i := 0; i < numOfRunner; i++ {
        go func(i int) {
            ret := runTask(i)
            ch <- ret
        }(i)
    }
    // 当 channel 中读取到数据即返回。
    return <-ch
}

func TestFirstResponse(t *testing.T) {
    t.Log("Before:", runtime.NumGoroutine())
    t.Log(FirstResponse())
    time.Sleep(time.Second * 1)
    t.Log("After:", runtime.NumGoroutine())
}
```

### 所有任务完成

```go
func runTask(id int) string {
    time.Sleep(10 * time.Millisecond)
    return fmt.Sprintf("The result is from %d", id)
}

func FirstResponse() string {
    numOfRunner := 10
    ch := make(chan string, numOfRunner)
    for i := 0; i < numOfRunner; i++ {
        go func(i int) {
            ret := runTask(i)
            ch <- ret
        }(i)
    }
    return <-ch
}

func AllResponse() string {
    // 创建 10 个协程执行任务
    numOfRunner := 10
    ch := make(chan string, numOfRunner)
    for i := 0; i < numOfRunner; i++ {
        go func(i int) {
            ret := runTask(i)
            ch <- ret
        }(i)
    }
    // 等待 10 个结果（10 个协程都执行完成返回）
    finalRet := ""
    for j := 0; j < numOfRunner; j++ {
        finalRet += <-ch + "\n"
    }
    return finalRet
}

func TestFirstResponse(t *testing.T) {
    t.Log("Before:", runtime.NumGoroutine())
    t.Log(AllResponse())
    time.Sleep(time.Second * 1)
    t.Log("After:", runtime.NumGoroutine())

}
```

### 对象池

```go
type ReusableObj struct {
}

type ObjPool struct {
    // 用于缓冲可重用对象
    bufChan chan *ReusableObj 
}

func NewObjPool(numOfObj int) *ObjPool {
    objPool := ObjPool{}
    objPool.bufChan = make(chan *ReusableObj, numOfObj)
    for i := 0; i < numOfObj; i++ {
        objPool.bufChan <- &ReusableObj{}
    }
    return &objPool
}

func (p *ObjPool) GetObj(timeout time.Duration) (*ReusableObj, error) {
    select {
    case ret := <-p.bufChan:
        return ret, nil
    // 超时控制
    case <-time.After(timeout): 
        return nil, errors.New("time out")
    }

}

func (p *ObjPool) ReleaseObj(obj *ReusableObj) error {
    select {
    case p.bufChan <- obj:
        return nil
    default:
        return errors.New("overflow")
    }
}

// 测试
func TestObjPool(t *testing.T) {
    pool := NewObjPool(10)
    // if err := pool.ReleaseObj(&ReusableObj{}); err != nil { //尝试放置超出池大小的对象
    //     t.Error(err)
    // }
    for i := 0; i < 11; i++ {
        if v, err := pool.GetObj(time.Second * 1); err != nil {
            t.Error(err)
        } else {
            fmt.Printf("%T\n", v)
            if err := pool.ReleaseObj(v); err != nil {
                t.Error(err)
            }
        }
    }
    fmt.Println("Done")
}
```

### 对象缓存（sync.Pool）

适用于通过复用降低复杂对象创建和 GC 代价。

协程安全，但会有锁的开销，因此需要权衡维护对象的开销与锁的开销。

取用：

-   尝试从私有对象获取，不存在则尝试从当前 Processor 的共享池获取。

-   如果当前 Processor 共享池也为空，则尝试从其他 Processor 的共享池获取。

-   如果所有⼦池都是空的，则由⽤户指定的 New 函数产⽣⼀个新的对象返回。

放回：如果私有对象不存在则保存为私有对象，否则放⼊当前 Processor ⼦池的共享池中。

Go GC 会清除 sync.pool 缓存的对象，因此对象的缓存有效期为下⼀次 GC 之前。由于生命周期受 GC 影响，不适用于连接池等，需要自己管理生命周期的资源池化。
```go
func TestSyncPool(t *testing.T) {
    pool := &sync.Pool{
        New: func() interface{} {
            fmt.Println("Create a new object.")
            return 100
        },
    }
    v := pool.Get().(int)
    fmt.Println(v)
    pool.Put(3)
    // GC 会清除 sync.pool 中缓存的对象
    runtime.GC() 
    v1, _ := pool.Get().(int)
    fmt.Println(v1)
}

func TestSyncPoolInMultiGroutine(t *testing.T) {
    pool := &sync.Pool{
        New: func() interface{} {
            fmt.Println("Create a new object.")
            return 10
        },
    }
    pool.Put(100)
    pool.Put(100)
    pool.Put(100)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            fmt.Println(pool.Get())
            wg.Done()
        }(i)
    }
    wg.Wait()
}
```

## 反射

-   `reflect.TypeOf` 返回类型 `reflct.Type`
-   `reflect.ValueOf` 返回值 `reflect.Value`（可以从中获取类型）

### 类型分支
类型 switch 就像一个常规的 switch 语句，但 case 指定类型（而不是值），并且这些值与给定接口值持有的值的类型进行比较。 
```go
func do(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Twice %v is %v\n", v, v*2)
    case string:
        fmt.Printf("%q is %v bytes long\n", v, len(v))
    default:
        fmt.Printf("I don't know about type %T!\n", v)
    }
}

func main() {
    do(21)
    do("hello")
    do(true)
}
```

### 深度比较

```go
func TestDeepEqual(t *testing.T) {
    a := map[int]string{1: "one", 2: "two", 3: "three"}
    b := map[int]string{1: "one", 2: "two", 3: "three"}
    //    t.Log(a == b)
    t.Log("a==b?", reflect.DeepEqual(a, b))

    s1 := []int{1, 2, 3}
    s2 := []int{1, 2, 3}
    s3 := []int{2, 3, 1}

    t.Log("s1 == s2?", reflect.DeepEqual(s1, s2))
    t.Log("s1 == s3?", reflect.DeepEqual(s1, s3))

    c1 := Customer{"1", "Mike", 40}
    c2 := Customer{"1", "Mike", 40}
    fmt.Println(c1 == c2)
    fmt.Println(reflect.DeepEqual(c1, c2))
}
```

### kind 判断类型

```go
// const ( 
//     Invalid Kind = iota 
//     Bool
//     Int 
//     Int8 
//     Int16 
//     Int32 
//     Int64 
//     Uint 
//     Uint8 
//     Uint16 
//     Uint32 
//     Uint64 
//     // ...
// )

func TestBasicType(t *testing.T) {
    var f float64 = 12
    t := reflect.TypeOf(&f)
    switch t.Kind() {
    case reflect.Float32, reflect.Float64:
        fmt.Println("Float")
    case reflect.Int, reflect.Int32, reflect.Int64:
        fmt.Println("Integer")
    default:
        fmt.Println("Unknown", t)
    }
}
```

### 根据名称调用

```go
func TestInvokeByName(t *testing.T) {
    e := &Employee{"1", "Mike", 30}
    // 按名称获取成员
    t.Logf("Name: value(%[1]v), Type(%[1]T) ", reflect.ValueOf(*e).FieldByName("Name"))
    if nameField, ok := reflect.TypeOf(*e).FieldByName("Name"); !ok {
        t.Error("Failed to get 'Name' field.")
    } else {
        t.Log("Tag:format", nameField.Tag.Get("format"))
    }
    // 按名称调用方法
    reflect.ValueOf(e).MethodByName("UpdateAge").Call([]reflect.Value{reflect.ValueOf(1)})
    t.Log("Updated Age:", e)
}
```

### Struct Tag

```go
type BasicInfo struct {
    Name string `json:"name"`
    Age int `json:"age"`
}

if nameField, ok := reflect.TypeOf(*e).FieldByName("Name"); !ok {
    t.Error("Failed to get 'Name' field.")
} else {
    t.Log("Tag:format", nameField.Tag.Get("format")) 
}
```

### 万能设值

```go
type Employee struct {
    EmployeeID string
    Name       string `format:"normal"`
    Age        int
}

func (e *Employee) UpdateAge(newVal int) {
    e.Age = newVal
}

type Customer struct {
    CookieID string
    Name     string
    Age      int
}

func fillBySettings(st interface{}, settings map[string]interface{}) error {
    if reflect.TypeOf(st).Kind() != reflect.Ptr {
        return errors.New("the first param should be a pointer to the struct type.")
    }
    // Elem() 获取指针指向的值
    if reflect.TypeOf(st).Elem().Kind() != reflect.Struct {
        return errors.New("the first param should be a pointer to the struct type.")
    }
    if settings == nil {
        return errors.New("settings is nil.")
    }
    var (
        field reflect.StructField
        ok    bool
    )
    for k, v := range settings {
        if field, ok = (reflect.ValueOf(st)).Elem().Type().FieldByName(k); !ok {
            continue
        }
        if field.Type == reflect.TypeOf(v) {
            vstr := reflect.ValueOf(st)
            vstr = vstr.Elem()
            vstr.FieldByName(k).Set(reflect.ValueOf(v))
        }
    }
    return nil
}

func TestFillNameAndAge(t *testing.T) {
    settings := map[string]interface{}{"Name": "Mike", "Age": 30}
    e := Employee{}
    if err := fillBySettings(&e, settings); err != nil {
        t.Fatal(err)
    }
    t.Log(e)
    c := new(Customer)
    if err := fillBySettings(c, settings); err != nil {
        t.Fatal(err)
    }
    t.Log(*c)
}
```

## Unsafe

### 类型转换

Go 本身不支持强制转换，利用 Unsafe 可以实现。

```go
type Customer struct {
    Name string
    Age  int
}

func TestUnsafe(t *testing.T) {
    i := 10
    f := *(*float64)(unsafe.Pointer(&i))
    t.Log(unsafe.Pointer(&i))
    t.Log(f)
}

type MyInt int

// 合理的类型转换
func TestConvert(t *testing.T) {
    a := []int{1, 2, 3, 4}
    b := *(*[]MyInt)(unsafe.Pointer(&a))
    t.Log(b)
}
```

### 原子操作

```go
func TestAtomic(t *testing.T) {
    var shareBufPtr unsafe.Pointer
    writeDataFn := func() {
        data := []int{}
        for i := 0; i < 100; i++ {
            data = append(data, i)
        }
        // sync/atomic
        atomic.StorePointer(&shareBufPtr, unsafe.Pointer(&data))
    }
    readDataFn := func() {
        data := atomic.LoadPointer(&shareBufPtr)
        fmt.Println(data, *(*[]int)(data))
    }
    var wg sync.WaitGroup
    writeDataFn()
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            for i := 0; i < 10; i++ {
                writeDataFn()
                time.Sleep(time.Microsecond * 100)
            }
            wg.Done()
        }()
        wg.Add(1)
        go func() {
            for i := 0; i < 10; i++ {
                readDataFn()
                time.Sleep(time.Microsecond * 100)
            }
            wg.Done()
        }()
    }
    wg.Wait()
}
```

## 性能分析

下载并安装：

-   [Graphviz](https://www.graphviz.org/)

-   [go-torch](https://github.com/orktes/go-torch)

-   [FlameGraph](https://github.com/brendangregg/FlameGraph)

通过⽂件⽅式输出 Profile：

-   灵活性⾼，适⽤于特定代码段的分析。

-   通过⼿动调⽤ runtime/pprof 的 API，相关⽂档 https://studygolang.com/static/pkgdoc/pkg/runtime_pprof.htm

-   `go tool pprof [binary] [binary.prof]`

通过 HTTP ⽅式输出 Profile：

-   简单，适合于持续性运⾏的应⽤：在应⽤程序中导⼊ `import _ "net/http/pprof"`，并启动 http server 即可（http://\<host\>:\<port\>/debug/pprof/）。

-   `go tool pprof http://<host>:<port>/debug/pprof/profile?seconds=10`

-   `go-torch -seconds 10 http://<host>:<port>/debug/pprof/profifile`

Go 支持的多种 Profile：https://golang.org/src/runtime/pprof/pprof.go

## 打印文本

```go
fmt.Println("Hello, 你好, नमस्ते, Привет, ᎣᏏᏲ")
p := struct { X, Y int }{ 17, 2 }
fmt.Println( "My point:", p, "x coord=", p.X )                  // 打印结构, ints, 等等
s := fmt.Sprintln( "My point:", p, "x coord=", p.X )            // 打印字符串

fmt.Printf("%d hex:%x bin:%b fp:%f sci:%e",17,17,17,17.0,17.0)  // c-ish 格式
s2 := fmt.Sprintf( "%d %f", 17, 17.0 )                          // 格式化打印到字符串变量 

hellomsg := `
 "Hello" in Chinese is 你好 ('Ni Hao')
 "Hello" in Hindi is नमस्ते ('Namaste')
` 
// 多行字符串文字，在开头和结尾使用反引号 
```

## 文件嵌入

Go 程序可以使用 `"embed"` 包嵌入静态文件，如下所示： 

```go
package main

import (
    "embed"
    "log"
    "net/http"
)

// content 包含静态内容（2 个文件）或 Web 服务器。 
// go:embed a.txt b.txt
var content embed.FS

func main() {
    http.Handle("/", http.FileServer(http.FS(content)))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

[完整示例](https://play.golang.org/p/pwWxdrQSrYv)

## HTTP 服务器

```go
package main

import (
    "fmt"
    "net/http"
)

// 定义响应类型
type Hello struct{}

// 该类型实现 ServeHTTP 方法（在接口 http.Handler 中定义） 
func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello!")
}

func main() {
    var h Hello
    http.ListenAndServe("localhost:4000", h)
}

// http.ServeHTTP 的方法签名：
// type Handler interface {
//     ServeHTTP(w http.ResponseWriter, r *http.Request)
// }
```

### 路由规则

- URL 分为两种，末尾是 /：表示⼀个⼦树，后⾯可以跟其他⼦路径； 末尾不

是 /，表示⼀个叶⼦，固定的路径。

- 以 / 结尾的 URL 可以匹配它的任何⼦路径，⽐如 /images 会匹配 /images/

cute-cat.jpg。

- 它采⽤最⻓匹配原则，如果有多个匹配，⼀定采⽤匹配路径最⻓的那个进⾏处

理。

- 如果没有找到任何匹配项，会返回 404 错误。

### Default Router

```go
func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    handler := sh.srv.Handler
    if handler == nil {
        handler = DefaultServeMux //使⽤缺省的Router
    }
    if req.RequestURI == "*" && req.Method == "OPTIONS" {
        handler = globalOptionsHandler{}
    }
    handler.ServeHTTP(rw, req)
}
```

### 更好的 Router

HttpRouter: https://github.com/julienschmidt/httprouter

```go
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

### Restful 服务

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/julienschmidt/httprouter"
)

type Employee struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var employeeDB map[string]*Employee

func init() {
    employeeDB = map[string]*Employee{}
    employeeDB["Mike"] = &Employee{"e-1", "Mike", 35}
    employeeDB["Rose"] = &Employee{"e-2", "Rose", 45}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func GetEmployeeByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    qName := ps.ByName("name")
    var (
        ok       bool
        info     *Employee
        infoJson []byte
        err      error
    )
    if info, ok = employeeDB[qName]; !ok {
        w.Write([]byte("{\"error\":\"Not Found\"}"))
        return
    }
    if infoJson, err = json.Marshal(info); err != nil {
        w.Write([]byte(fmt.Sprintf("{\"error\":,\"%s\"}", err)))
        return
    }

    w.Write(infoJson)
}

func main() {
    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/employee/:name", GetEmployeeByName)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## JSON 解析

### Struct Tag 解析

```go
type BasicInfo struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
type JobInfo struct {
    Skills []string `json:"skills"`
}
type Employee struct {
    BasicInfo BasicInfo `json:"basic_info"`
    JobInfo   JobInfo   `json:"job_info"`
}

var jsonStr = `{
    "basic_info":{
          "name":"Mike",
        "age":30
    },
    "job_info":{
        "skills":["Java","Go","C"]
    }
}    `

func TestEmbeddedJson(t *testing.T) {
    e := new(Employee)
    err := json.Unmarshal([]byte(jsonStr), e)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(*e)
    if v, err := json.Marshal(e); err == nil {
        fmt.Println(string(v))
    } else {
        t.Error(err)
    }

}
```

### 第三方库

EasyJSON: `go get -u github.com/mailru/easyjson/...`