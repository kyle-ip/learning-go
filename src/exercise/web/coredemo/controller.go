package main

import (
	"context"
	"coredemo/framework"
	"fmt"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	// 为单个请求设置超时：
	// 1. 继承 request 的 Context，创建出一个设置超时时间的 Context。

	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	// 2. 创建一个新的 Goroutine 来处理具体的业务逻辑。
	// mu := sync.Mutex{}
	go func() {
		// 处理异常。
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		time.Sleep(10 * time.Second)
		c.SetOkStatus().Json("ok")

		finish <- struct{}{}
	}()

	select {
	// 主 Goroutine 中读取到异常，返回 500。
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.SetStatus(500).Json("panic")
	// 主 Goroutine 中读取到完成，返回 200。
	case <-finish:
		fmt.Println("finish")

	// 3. 设计事件处理顺序，当前 Goroutine 监听超时时间 Context 的 Done() 事件，和具体的业务处理结束事件，先到先处理。
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.SetStatus(500).Json("time out")
		c.SetHasTimeout()
	}
	return nil
}

//func Foo(request *http.Request, response http.ResponseWriter) {
//    obj := map[string]interface{}{
//        "errno":  50001,
//        "errmsg": "inner error",
//        "data":   nil,
//    }
//
//    response.Header().Set("Content-Type", "application/json")
//
//    foo := request.PostFormValue("foo")
//    if foo == "" {
//        foo = "10"
//    }
//    fooInt, err := strconv.Atoi(foo)
//    if err != nil {
//        response.WriteHeader(500)
//        return
//    }
//    obj["data"] = fooInt
//    byt, err := json.Marshal(obj)
//    if err != nil {
//        response.WriteHeader(500)
//        return
//    }
//    response.WriteHeader(200)
//    response.Write(byt)
//    return
//}
//
//func Foo2(ctx *framework.Context) error {
//    obj := map[string]interface{}{
//        "errno":  50001,
//        "errmsg": "inner error",
//        "data":   nil,
//    }
//
//    fooInt := ctx.FormInt("foo", 10)
//    obj["data"] = fooInt
//    return ctx.Json(http.StatusOK, obj)
//}
//
//func Foo3(ctx *framework.Context) error {
//    rdb := redis.NewClient(&redis.Options{
//        Addr:     "localhost:6379",
//        Password: "", // no password set
//        DB:       0,  // use default DB
//    })
//
//    return rdb.Set(ctx, "key", "value", 0).Err()
//}
