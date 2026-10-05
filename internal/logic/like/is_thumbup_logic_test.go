package likelogic

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLikeStatusCacheKey(t *testing.T) {
	bizId := "article"
	targetId := int64(1001)
	userId := int64(2002)

	expected := "biz#like#status:article:1001:2002"
	key := fmt.Sprintf("biz#like#status:%s:%d:%d", bizId, targetId, userId)

	if key != expected {
		t.Fatalf("expected key %s, got %s", expected, key)
	}
}

func TestLikeStatusParse(t *testing.T) {
	val := "1"
	likeType, err := strconv.Atoi(val)
	if err != nil {
		t.Fatalf("failed to parse likeType: %v", err)
	}
	if likeType != 1 {
		t.Fatalf("expected likeType 1, got %d", likeType)
	}

	emptyVal := "0"
	emptyType, _ := strconv.Atoi(emptyVal)
	if emptyType != 0 {
		t.Fatalf("expected emptyType 0, got %d", emptyType)
	}
}
