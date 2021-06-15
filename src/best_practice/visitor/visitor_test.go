package visitor

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "testing"
)

// Kubernetes 的 kubectl 命令中用到 Builder 和 Visitor 模式。
// 其中 Visitor 模式是将算法与操作对象的结构分离的方法，能在不修改结构的情况下向现有对象结构添加新操作，是遵循开放 / 封闭原则的一种方法。

type Visitor func(shape Shape)

type Shape interface {
    accept(Visitor)
}

type Circle struct {
    Radius int
}

func (c Circle) accept(v Visitor) {
    v(c)
}

type Rectangle struct {
    Width, Heigh int
}

func (r Rectangle) accept(v Visitor) {
    v(r)
}

// 两个 Visitor，表示以不同方式进行序列化。

func JsonVisitor(shape Shape) {
    bytes, err := json.Marshal(shape)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(bytes))
}

func XmlVisitor(shape Shape) {
    bytes, err := xml.Marshal(shape)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(bytes))
}

func Test30(t *testing.T) {
    // 这段代码的目的是解耦数据结构和算法：
    // 虽然使用 Strategy 模式也可以完成而且会比较干净，
    // 但是在有些情况下多个 Visitor 是来访问一个数据结构的不同部分，
    // 此时数据结构比较像一个数据库，而各个 Visitor 会成为一个个小应用。

    // 对于不同的 Shape，可以调用不同的 Visitor，实现多种方式的序列化。
    c := Circle{10}
    r := Rectangle{100, 200}
    shapes := []Shape{c, r}

    for _, s := range shapes {
        s.accept(JsonVisitor)
        s.accept(XmlVisitor)
    }

}
