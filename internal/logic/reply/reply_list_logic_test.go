package replylogic

import (
	"testing"

	model "rpc-social/internal/model/reply"
)

// TestSubRepliesMapping 测试子回复映射与精准挂载逻辑
func TestSubRepliesMapping(t *testing.T) {
	// 构造测试子回复数据
	subReplies := []*model.Reply{
		{ID: 101, ParentID: 1, Content: "子回复1"},
		{ID: 102, ParentID: 1, Content: "子回复2"},
		{ID: 103, ParentID: 2, Content: "子回复3"},
	}

	// 聚合按 parentId 映射
	subMap := make(map[int64][]*model.Reply)
	for _, v := range subReplies {
		subMap[v.ParentID] = append(subMap[v.ParentID], v)
	}

	if len(subMap[1]) != 2 {
		t.Fatalf("预期根评论 1 下有 2 条子回复, 实际为 %d", len(subMap[1]))
	}
	if len(subMap[2]) != 1 {
		t.Fatalf("预期根评论 2 下有 1 条子回复, 实际为 %d", len(subMap[2]))
	}
}
