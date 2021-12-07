package framework

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const shortDuration = 1 * time.Millisecond

func TestContext(t *testing.T) {
	// 创建有截止时间的 Context，定时结束时主线程通过 Done() 函数收到事件结束通知。
	// 然后主动调用函数句柄 CancelFunc 来通知所有子 Context 结束。
	// CancelFunc 是主动让下游结束，而 Done 是被上游通知结束。

	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// 最终調用 cancel。
	defer cancel()

	// 使用 select 监听 1s 和有截止时间的 Context 哪个先结束
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

}
