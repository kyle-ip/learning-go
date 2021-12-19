package visitor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

// Kubernetes 的 kubectl 命令中用到 Builder 和 Visitor 模式。
// 其中 Visitor 模式是将算法与操作对象的结构分离的方法，不修改结构、向现有对象结构添加新操作，遵循开放封闭原则。

type Visitor func(shape Shape)

// Circle 和 Rectangle 都属于 Shape、都实现其 accept 方法。

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
	Width, Height int
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

func TestVisitor(t *testing.T) {
	// 解耦数据结构和算法：
	// 虽然使用 Strategy 模式也可以完成而且会比较干净，
	// 但是在有些情况下多个 Visitor 是来访问一个数据结构的不同部分。数据结构像一个数据库，而各个 Visitor 会成为一个个小应用。

	// 对于不同的 Shape，可以调用不同的 Visitor，实现多种方式的序列化。
	for _, s := range []Shape{
		Circle{10},
		Rectangle{100, 200},
	} {
		s.accept(JsonVisitor)
		s.accept(XmlVisitor)
	}

}
