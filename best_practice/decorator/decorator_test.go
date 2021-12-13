package decorator

// Python 装饰器的函数式编程     https://coolshell.cn/articles/11265.html
// 函数式编程                  https://coolshell.cn/articles/10822.html
// 相比之下 Python 是动态语言、Java 是虚拟机语言，Go 作为静态语言在编译期就需要确定好类型（否则无法通过）。
// The Laws of Reflection    https://blog.golang.org/laws-of-reflection

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// 高阶函数 decorator() 在调用时，先把 Hello() 函数传进去，然后会返回一个匿名函数。
// 匿名函数中除了运行自己的代码，也调用了被传入的 Hello() 函数。

func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("Started")
		f(s)
		fmt.Println("Done")
	}
}

func Hello(s string) {
	fmt.Println(s)
}

func TestDecorator(t *testing.T) {
	hello := decorator(Hello)
	hello("Hello")
	// decorator(Hello)("Hello, World!")
}

// ========== 时间统计 ==========

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {
		// 最后输出运行时间。
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n", getFunctionName(f), time.Since(t))
		}(time.Now())
		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		sum += i
	}
	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

func TestTimeSum(t *testing.T) {
	// 封装时间统计的 decorator。
	sum1 := timedSumFunc(Sum1)
	sum2 := timedSumFunc(Sum2)
	fmt.Printf("%d, %d\n", sum1(-10000, 10000000), sum2(-10000, 10000000))
}

// ========== 请求日志 ==========

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader()")
		w.Header().Set("Server", "HelloServer v0.0.1")
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

func TestServerHeader(t *testing.T) {
	http.HandleFunc("/v1/hello", WithServerHeader(hello))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// ========== 多装饰器 Pipeline ==========

type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}

// http.HandleFunc("/v4/hello", Handler(hello, WithServerHeader, WithBasicAuth, WithDebugLog))

// ========== 泛型装饰器 ==========

func Decorator(decoPtr, fn interface{}) (err error) {
	// 两个参数：出参 decoPtr（装饰后的函数），入参 fn（待装饰的函数）。
	var decoratedFunc, targetFunc reflect.Value
	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	v := reflect.MakeFunc(targetFunc.Type(), func(in []reflect.Value) (out []reflect.Value) {
		fmt.Println("before")
		out = targetFunc.Call(in)
		fmt.Println("after")
		return
	})
	decoratedFunc.Set(v)
	return
}

func foo(a, b, c int) int {
	fmt.Printf("%d, %d, %d \n", a, b, c)
	return a + b + c
}

func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

func TestGenericsDecorator(t *testing.T) {
	mybar := bar
	err := Decorator(&mybar, bar)
	if err != nil {
		return
	}
	mybar("hello,", "world!")
}
