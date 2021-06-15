package pipeline

// [Go Concurrency Patterns: Pipelines and cancellation](https://blog.golang.org/pipelines)
// [Google I/O 2012 - Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
// [Advanced Go Concurrency Patterns](https://blog.golang.org/advanced-go-concurrency-patterns)
// [Squinting at Power Series](https://swtch.com/~rsc/thread/squint.pdf)

import (
    "fmt"
    "math"
    "net/http"
    "sync"
    "testing"
)

type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

type PipelineHttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func PipelineHandler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
    for i := range decors {
        d := decors[len(decors)-1-i] // iterate in reverse
        h = d(h)
    }
    return h
}

// http.HandleFunc("/v3/hello", WithServerHeader(WithBasicAuth(WithDebugLog(hello))))

// 使用 Pipeline 代替多层嵌套：
// http.PipelineHandler("/v4/hello", Handler(hello, WithServerHeader, WithBasicAuth, WithDebugLog))

// ========== Channel 转发函数 ==========

func echo(nums []int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func odd(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 != 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func sum(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        var sum = 0
        for n := range in {
            sum += n
        }
        out <- sum
        close(out)
    }()
    return out
}

type EchoFunc func([]int) <-chan int

type PipeFunc func(<-chan int) <-chan int

func pipeline(nums []int, echo EchoFunc, pipeFns ...PipeFunc) <-chan int {
    ch := echo(nums)
    for i := range pipeFns {
        ch = pipeFns[i](ch)
    }
    return ch
}

func TestChannelForward(t *testing.T) {
    // 效果：echo $nums | square | sum
    var nums1 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    for n := range sum(square(odd(echo(nums1)))) {
        fmt.Println(n)
    }
    // 或使用代理函数完成
    // var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    // for n := range pipeline(nums, gen, odd, sq, sum) {
    // 	 fmt.Println(n)
    // }
}

// ========== Fan in/Out ==========
// Goroutine 和 Channel 可以写出 1 对多或多对 1 的 Pipeline，即 Fan In/ Fan Out。

func makeRange(min, max int) []int {
    a := make([]int, max-min+1)
    for i := range a {
        a[i] = min + i
    }
    return a
}

func isPrime(value int) bool {
    for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
        if value%i == 0 {
            return false
        }
    }
    return value > 1
}

func prime(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if isPrime(n) {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func merge(cs []<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    wg.Add(len(cs))
    for _, c := range cs {
        go func(c <-chan int) {
            for n := range c {
                out <- n
            }
            wg.Done()
        }(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func TestFanInFanOut(t *testing.T) {
    // 构造数组 [1, 10000]。
    nums := makeRange(1, 10000)

    // fan out: 把数组 echo 到一个 Channel in 中。
    in := echo(nums)

    // 生成 5 个 Channel，都调用 sum(prime(in)) ，于是每个 Sum 的 Goroutine 都会开始计算和；
    const nProcess = 5
    var chans [nProcess]<-chan int

    for i := range chans {
        chans[i] = sum(prime(in))
    }

    // fan in: 最后再把所有的结果再求和拼起来得到最终的结果。
    for n := range sum(merge(chans[:])) {
        fmt.Println(n)
    }
}
