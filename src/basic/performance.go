package basic

// 无锁队列的实现          https://coolshell.cn/articles/8239.html
// 无锁 Hashmap 实现      https://coolshell.cn/articles/9703.html

// Go 语言仍需要关心性能问题：
// - 如果需要把数字转换成字符串，使用 strconv.Itoa() 比 fmt.Sprintf() 要快一倍左右。
// - 尽可能避免把 String 转成 []Byte，会导致性能下降。
// - 在 for-loop 里对 Slice 使用 append()，先把 Slice 的容量扩充到位，可以避免内存重新分配以及系统自动按 2 的 N 次方幂扩展但又用不到，避免浪费内存。
// - 使用 StringBuffer 或是 StringBuild 来拼接字符串，性能会比使用 + 或 += 高三到四个数量级。
// - 尽可能使用并发的 goroutine，然后使用 sync.WaitGroup 来同步分片操作。
// - 避免在热代码中分配内存，会导致 gc 繁忙。尽可能使用 sync.Pool 来重用对象。
// - 使用 lock-free 的操作，避免使用 mutex，尽可能使用 sync/Atomic包。
// - I/O 是个非常非常慢的操作，使用 bufio.NewWrite() 和 bufio.NewReader() 可以带来更高的性能。
// - 在 for-loop 里的固定的正则表达式，一定要使用 regexp.Compile() 编译正则表达式，性能会提升两个数量级。
// - 如果需要更高性能的协议，就要考虑使用 protobuf 或 msgp 而不是 JSON（序列化和反序列化中使用了反射）。
// - 在使用 Map 时，使用整型的 key 会比字符串的要快。
