package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestThreadPool(t *testing.T) {
	a := func() {
		fmt.Println(time.Now())
	}
	pool := NewPool(1)
	go func() {
		for i := 0; i < 2; i++ {
			time.Sleep(1 * time.Second)
			pool.AddTask(a)
		}
		time.Sleep(2 * time.Second)
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("pool is closed")
			}
		}()
		pool.AddTask(a)
	}()
	//启动协程池p
	pool.Run()
	time.Sleep(3 * time.Second)
	pool.Stop()
	time.Sleep(5 * time.Second)
}
