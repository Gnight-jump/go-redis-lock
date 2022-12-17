package alock

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLockTransaction(t *testing.T) {
	err := InitRedisLock("localhost:6379", "", 0, 10)
	if err != nil {
		t.Fatalf("init redis failed, err=%v", err)
	}
	wait := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go SayHello(t, fmt.Sprint(i), wait)
		time.Sleep(5 * time.Second)
	}
	wait.Wait()
}

func SayHello(t *testing.T, reqID string, wait *sync.WaitGroup) {
	defer func() {
		wait.Done()
	}()

	result, err := LockTransaction("say_hello", reqID, func() error {
		for i := 0; i < 3; i++ {
			time.Sleep(5 * time.Second)
			t.Logf("hello reqID=%v", reqID)
		}
		return nil
	})
	if !result {
		t.Logf("say hello fail, lock exist, reqID=%v", reqID)
		return
	}
	if err != nil {
		t.Logf("say hello fail, process err=%v, reqID=%v", err, reqID)
		return
	}
}
