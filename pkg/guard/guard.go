package guard

import (
	"sync"
	"time"
)

// MemoryGuard 轻量级操作防重与频控限制器
type MemoryGuard struct {
	mu    sync.Mutex
	locks map[string]int64
	ttl   time.Duration
}

// NewMemoryGuard 创建防刷限制器实例
func NewMemoryGuard(ttl time.Duration) *MemoryGuard {
	g := &MemoryGuard{
		locks: make(map[string]int64),
		ttl:   ttl,
	}
	return g
}

// Acquire 尝试获取操作锁，成功返回 true，在锁定期内返回 false
func (g *MemoryGuard) Acquire(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()
	if expireAt, exists := g.locks[key]; exists {
		if now < expireAt {
			return false
		}
	}

	g.locks[key] = now + g.ttl.Milliseconds()
	return true
}

// Release 手动释放锁
func (g *MemoryGuard) Release(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.locks, key)
}
