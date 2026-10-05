package guard

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryGuard_AcquireAndRelease(t *testing.T) {
	g := NewMemoryGuard(100 * time.Millisecond)

	// 第一次获取成功
	if !g.Acquire("user:101:thumbup") {
		t.Fatalf("expected first acquire to succeed")
	}

	// 锁定期内重复获取失败
	if g.Acquire("user:101:thumbup") {
		t.Fatalf("expected concurrent acquire to be blocked")
	}

	// 等待过期后再次获取成功
	time.Sleep(120 * time.Millisecond)
	if !g.Acquire("user:101:thumbup") {
		t.Fatalf("expected acquire after ttl to succeed")
	}
}

func TestMemoryGuard_Concurrent(t *testing.T) {
	g := NewMemoryGuard(500 * time.Millisecond)
	key := "user:999:reply"

	var successCount int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if g.Acquire(key) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 success in concurrent acquire, got %d", successCount)
	}
}
