package main

import (
	"context"
	"fmt"
	"galactus/common/middleware/concurrent"
	"time"
)

var tl = concurrent.NewThreadLocal()

func demo1() {
	tl.Set("main test")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 使用上下文
	go doWorkWithContext(ctx)

	// 使用普通参数
	go doWorkWithParams("Hello, World!")
	// 等待
	time.Sleep(4 * time.Second)
}

func doWorkWithContext(ctx context.Context) {
	select {
	case <-time.After(1 * time.Second):
		tl.Set("test")
		fmt.Println("Work completed with context", getTest())
	case <-ctx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("Work canceled with context:", ctx.Err())
	}
}

func getTest() any {
	v, _ := tl.Get()
	return v
}

func doWorkWithParams(msg string) {
	time.Sleep(1 * time.Second)
	fmt.Println("Work completed with message:", msg, getTest())
}
