package best_practise

import (
	"fmt"
	"testing"
)

// IOC/DIP 其实是一种管理思想        https://coolshell.cn/articles/9949.html

// ========== 结构体嵌入 ==========

type Widget struct {
	X, Y int
}

type Label struct {
	Widget        // Embedding (delegation)
	Text   string // Aggregation
}

type Button struct {
	Label // Embedding (delegation)
}

type ListBox struct {
	Widget          // Embedding (delegation)
	Texts  []string // Aggregation
	Index  int      // Aggregation
}

func Test15(t *testing.T) {
	label := Label{Widget{10, 10}, "State: "}
	label.X = 11
	label.Y = 12
}

// ========== 方法重写 ==========

// 两个接口：Painter 画出组件；Clicker 表明点击事件。
// 对于 Lable 而言只有 Painter 没有 Clicker；对于 Button 和 ListBox 而言 Painter 和 Clicker 都有。

type Painter interface {
	Paint()
}

type Clicker interface {
	Click()
}

func (label Label) Paint() {
	fmt.Printf("%p:Label.Paint(%q)\n", &label, label.Text)
}

// 这个接口可以通过 Label 的嵌入带到新的结构体，所以可以在 Button 中重载这个接口方法。

func (button Button) Paint() { // Override
	fmt.Printf("Button.Paint(%s)\n", button.Text)
}

func (button Button) Click() {
	fmt.Printf("Button.Click(%s)\n", button.Text)
}

func (listBox ListBox) Paint() {
	fmt.Printf("ListBox.Paint(%q)\n", listBox.Texts)
}

func (listBox ListBox) Click() {
	fmt.Printf("ListBox.Click(%q)\n", listBox.Texts)
}

// ========== 嵌入结构多态 ==========

func NewButton(X, Y int, Text string) Button {
	return Button{Label{Widget{X, Y}, Text}}
}

func Test16(t *testing.T) {
	button1 := Button{Label{Widget{10, 70}, "OK"}}
	button2 := NewButton(50, 70, "Cancel")
	listBox := ListBox{Widget{10, 40}, []string{"AL", "AK", "AZ", "AR"}, 0}

	label := Label{Widget{5, 10}, "ABC"}

	for _, painter := range []Painter{label, listBox, button1, button2} {
		painter.Paint()
	}

	for _, widget := range []interface{}{label, listBox, button1, button2} {
		widget.(Painter).Paint()
		if clicker, ok := widget.(Clicker); ok {
			clicker.Click()
		}
		fmt.Println() // print a empty line
	}

	// 可以使用接口或泛型的 interface{} 来实现多态（但是需要有一个类型转换）。
}
