package snowflake

import (
	"sync"
	"testing"
)

func TestSnowflake_GenerateID_Unique(t *testing.T) {
	node, err := NewNode(2)
	if err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	const count = 5000
	var wg sync.WaitGroup
	var idMap sync.Map

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := node.Generate()
			if _, loaded := idMap.LoadOrStore(id, true); loaded {
				t.Errorf("duplicate id generated: %d", id)
			}
		}()
	}
	wg.Wait()
}
