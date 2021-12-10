package template

import (
	"fmt"
	"testing"
)

type Cooker interface {
	fire()
	cooke()
	outfire()
}

// CookMenu 类似于一个抽象类
type CookMenu struct {
}

func (CookMenu) fire() {
	fmt.Println("开火")
}

// 做菜，交给具体的子类实现
func (CookMenu) cooke() {
}

func (CookMenu) outfire() {
	fmt.Println("关火")
}

// doCook 封装具体步骤，即对外提供的模板
func doCook(cook Cooker) {
	cook.fire()
	cook.cooke()
	cook.outfire()
}

type Tomato struct {
	CookMenu
}

func (*Tomato) cooke() {
	fmt.Println("做西红柿")
}

type Egg struct {
	CookMenu
}

func (Egg) cooke() {
	fmt.Println("做炒鸡蛋")
}

func TestTemplate(t *testing.T) {
	// 做西红柿
	tomato := &Tomato{}
	doCook(tomato)

	fmt.Println("\n=====> 做另外一道菜")
	// 做炒鸡蛋
	egg := &Egg{}
	doCook(egg)

}
