package feed

import (
	"testing"
)

func TestFeedMerger_Merge(t *testing.T) {
	// 模拟普通博主推入的 Inbox 数据
	inbox := []*FeedItem{
		{ArticleID: 103, PublishTime: 1003},
		{ArticleID: 101, PublishTime: 1001},
	}

	// 模拟所关注大 V 的 Outbox 数据
	vOutbox1 := []*FeedItem{
		{ArticleID: 104, PublishTime: 1004},
		{ArticleID: 102, PublishTime: 1002},
	}
	vOutbox2 := []*FeedItem{
		{ArticleID: 105, PublishTime: 1005},
	}

	merger := NewFeedMerger()
	merged := merger.MergeMultiStreams([][]*FeedItem{inbox, vOutbox1, vOutbox2}, 4)

	// 预期取最新 4 条，按发布时间倒序排列
	if len(merged) != 4 {
		t.Fatalf("expected 4 merged items, got %d", len(merged))
	}
	if merged[0].ArticleID != 105 || merged[1].ArticleID != 104 || merged[2].ArticleID != 103 || merged[3].ArticleID != 102 {
		t.Fatalf("incorrect merge order: %+v", merged)
	}
}
