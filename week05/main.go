package main

import (
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"
)

// Metric 度量值
type Metric struct {

    // 成功请求数
    Success int

    // 失败请求数
    Failure int
}

//
func (m *Metric) String() string {
    return fmt.Sprintf("Success: %d, Failure: %d", m.Success, m.Failure)
}

// SlidingWindow 滑动窗口
type SlidingWindow struct {

    // 度量值总量
    Metric

    // 度量值
    metrics []*Metric

    // 容量
    capacity int

    // 互斥锁
    mutex sync.Mutex

    // 最大失败率
    maxFailureRate float64

    // 最大请求数
    maxRequestCount int

    // 限流回调
    rateLimitCallback func(msg string) error

    // 度量值输入
    metricIn chan *Metric
}

// Add 添加最新的度量
func (sw *SlidingWindow) Add(metric *Metric) {
    log.Println("----------")
    log.Println("new metric:", metric)
    sw.mutex.Lock()
    defer sw.mutex.Unlock()

    // 窗口已满，计算窗口值时减去最后一位，否则减去 0。
    var lastMetric *Metric
    if sw.capacity == len(sw.metrics) {
        lastMetric = sw.metrics[0]
        sw.metrics = sw.metrics[1:]
    } else {
        lastMetric = &Metric{}
    }
    sw.metrics = append(sw.metrics, metric)
    sw.Success += metric.Success - lastMetric.Success
    sw.Failure += metric.Failure - lastMetric.Failure

    log.Println(sw)

    requestCount := sw.Success + sw.Failure
    failureRate := float64(sw.Failure) / float64(sw.Success+sw.Failure)

    // 时间窗口内请求数达到阈值，触发限流。
    if requestCount >= sw.maxRequestCount {
        _ = (sw.rateLimitCallback)(fmt.Sprintf("requestCount: %d", requestCount))
        return
    }

    // 时间窗口内错误率达到阈值，触发限流。
    if failureRate > sw.maxFailureRate {
        _ = (sw.rateLimitCallback)(fmt.Sprintf("failure rate: %.2f", failureRate))
        return
    }
}

// Size 度量值数量
func (sw *SlidingWindow) Size() int {
    return len(sw.metrics)
}

// Listen 从 chan 读取
func (sw *SlidingWindow) Listen() error {
    for metric := range sw.metricIn {
        sw.Add(metric)
    }
    return nil
}

func (sw *SlidingWindow) String() string {
    return fmt.Sprintf("window [Size: %d, Success: %d, Failure: %d]", sw.Size(), sw.Success, sw.Failure)
}

// NewSlidingWindow 创建滑动窗口
func NewSlidingWindow(n int, maxRequestCount int, maxFailureRate float64, metricIn chan *Metric, callback func(msg string) error) (sw *SlidingWindow) {
    return &SlidingWindow{
        metrics:           make([]*Metric, 0, n),
        capacity:          n,
        rateLimitCallback: callback,
        maxRequestCount:   maxRequestCount,
        maxFailureRate:    maxFailureRate,
        metricIn:          metricIn,
    }
}

//
func main() {

    // 采集间隔，滑动窗口大小，最大请求数，最大失败率，度量值输入，限流回调。
    interval := 1 * time.Second
    windowSize := 10
    maxRequestCount := 11000
    maxFailureRate := 0.6
    metricIn := make(chan *Metric, windowSize)
    rateLimitCallback := func(msg string) error {
        log.Printf("rate limiting! %s \n", msg)
        return nil
    }

    sw := NewSlidingWindow(windowSize, maxRequestCount, maxFailureRate, metricIn, rateLimitCallback)

    // 模拟每秒采集一次数据：成功、失败的请求数。
    go func() {
        rand.Seed(time.Now().UnixNano())
        for {
            time.Sleep(interval)
            metricIn <- &Metric{Success: rand.Int() % 1000, Failure: rand.Int() % 1000}
        }
    }()

    sw.Listen()
}
