package exception

// [由苹果的低级 BUG 想到的](https://coolshell.cn/articles/11112.html)
// [Golang Error Handling lesson by Rob Pike](http://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike)
// [Errors are values](https://blog.golang.org/errors-are-values)

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"testing"
)

// ========== 返回值检查 ==========
// Go 函数支持多返回值，可在返回接口把业务语义（业务返回值）和控制语义（出错返回值）区分开。
// Go 的很多函数都会返回 result、err 两个值：
// - 参数上基本上就是入参，而返回接口把结果和错误分离，使得函数的接口语义清晰；
// - Go 中的错误参数如果要忽略，需要用 _ 变量来显式地忽略；
// - 因为返回的 error 是接口（其中只有一个方法 Error()，返回一个 string），所以可以扩展自定义的错误处理。

// 如果一个函数返回了多个不同类型的 error：
// if err != nil {
//     switch err.(type) {
//     case *json.SyntaxError:
//         ...
//     case *ZeroDivisionError:
//         ...
//     case *NullPointerError:
//         ...
//     default:
//         ...
//     }
// }
// 因此错误处理的方式本质上是返回值检查，也兼顾了异常的一些好处：对错误的扩展。

// ========== 资源清理 ==========
// Go 使用 defer 关键词进行资源清理（类似 Java 的 finally）。

func Open(s string) (r io.Closer, err error) {
	return nil, nil
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestDefer(t *testing.T) {
	r, err := Open("a")
	if err != nil {
		log.Fatalf("error opening 'a'\n")
	}
	defer func() {
		Close(r)
		fmt.Println("close a")
	}()

	r, err = Open("b")
	if err != nil {
		log.Fatalf("error opening 'b'\n")
	}
	defer func() {
		Close(r)
		fmt.Println("close b")
	}()
}

// ========== 函数式编程 ==========
// 要避免大量的 if err != nil 这种判断代码，可以使用函数式编程（类似 Java 的 Optional）。

type Point struct {
	Longitude     int
	Latitude      int
	Distance      int
	ElevationGain int
	ElevationLoss int
}

func Read(r io.Reader) (*Point, error) {
	var p Point
	var err error
	// 使用 Closure 的方式把对每个字段的 if xxx != nil 抽出来定义一个函数 read。
	// 但是会带来一个问题：有一个 err 变量和一个内部的函数，不是很干净。
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(r, binary.BigEndian, data)
	}
	read(&p.Longitude)
	read(&p.Latitude)
	read(&p.Distance)
	read(&p.ElevationGain)
	read(&p.ElevationLoss)

	if err != nil {
		return &p, err
	}

	return &p, nil
}

type Reader struct {
	r   io.Reader
	err error
}

func (r *Reader) read(data interface{}) {
	if r.err == nil {
		r.err = binary.Read(r.r, binary.BigEndian, data)
	}
}

func parse(input io.Reader) (*Point, error) {
	var p Point
	r := Reader{r: input}

	r.read(&p.Longitude)
	r.read(&p.Latitude)
	r.read(&p.Distance)
	r.read(&p.ElevationGain)
	r.read(&p.ElevationLoss)

	if r.err != nil {
		return nil, r.err
	}

	return &p, nil
}

var b = []byte{0x48, 0x61, 0x6f, 0x20, 0x43, 0x68, 0x65, 0x6e, 0x00, 0x00, 0x2c}

var r = bytes.NewReader(b)

// ========== 流式接口（Fluent Interface）：在结构中记录错误 ==========

type User struct {
	Name   [10]byte
	Age    uint8
	Weight uint8
	err    error
}

func (p *User) read(data interface{}) {
	if p.err == nil {
		p.err = binary.Read(r, binary.BigEndian, data)
	}
}

func (p *User) ReadName() *User {
	p.read(&p.Name)
	return p
}
func (p *User) ReadAge() *User {
	p.read(&p.Age)
	return p
}
func (p *User) ReadWeight() *User {
	p.read(&p.Weight)
	return p
}
func (p *User) Print() *User {
	if p.err == nil {
		fmt.Printf("Name=%s, Age=%d, Weight=%d\n", p.Name, p.Age, p.Weight)
	}
	return p
}

func TestFluentInterface(t *testing.T) {
	p := User{}
	p.
		ReadName().
		ReadAge().
		ReadWeight().
		Print()
	fmt.Println(p.err) // EOF 错误
}

// ========== 封装错误 ==========

// 在 Go 的开发中，更为普遍的做法是将错误包装在另一个错误中，同时保留原始内容：

type authorizationError struct {
	operation string
	err       error // original error
}

func (e *authorizationError) Error() string {
	return fmt.Sprintf("authorization failed during %s: %v", e.operation, e.err)
}

// 还可以通过标准的访问方法，使用一个接口比如 causer接口中实现 Cause() 方法来暴露原始错误，以供进一步检查：

type causer interface {
	Cause() error
}

func (e *authorizationError) Cause() error {
	return e.err
}

// 或利用第三方库 error：https://github.com/pkg/errors
func TestPkgErrors(t *testing.T) {

	// 错误包装
	// if err != nil {
	//     return errors.Wrap(err, "read failed")
	// }

	// Cause 接口
	// switch err := errors.Cause(err).(type) {
	//     case *MyError:
	//     default:
	// }
}
