package feed

import (
	"container/heap"
)

type FeedMerger struct{}

func NewFeedMerger() *FeedMerger {
	return &FeedMerger{}
}

// heapNode 存储大顶堆节点的流索引、当前流中元素索引及元素指针
type heapNode struct {
	streamIdx int
	itemIdx   int
	item      *FeedItem
}

// maxHeap 基于 container/heap 实现的大顶堆（严格全序：时间倒序优先，时间相同时按 ArticleID 倒序）
type maxHeap []*heapNode

func (h maxHeap) Len() int { return len(h) }
func (h maxHeap) Less(i, j int) bool {
	if h[i].item.PublishTime == h[j].item.PublishTime {
		return h[i].item.ArticleID > h[j].item.ArticleID
	}
	return h[i].item.PublishTime > h[j].item.PublishTime
}
func (h maxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(*heapNode))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MergeMultiStreams 使用 K 路大顶堆算法 (K-Way Max-Heap Merge) 将多个已倒序排列的动态流归并、去重并截断
// 时间复杂度：O(limit * log(N))，其中 N 为并发流数；空间复杂度：O(N) 辅助堆空间，零冗余全量数据拷贝。
func (m *FeedMerger) MergeMultiStreams(streams [][]*FeedItem, limit int) []*FeedItem {
	if len(streams) == 0 || limit <= 0 {
		return []*FeedItem{}
	}

	h := &maxHeap{}
	heap.Init(h)

	// 初始化：将每个非空流的第一个元素压入大顶堆
	for sIdx, stream := range streams {
		if len(stream) > 0 && stream[0] != nil {
			heap.Push(h, &heapNode{
				streamIdx: sIdx,
				itemIdx:   0,
				item:      stream[0],
			})
		}
	}

	result := make([]*FeedItem, 0, limit)
	seen := make(map[int64]bool)

	// 迭代弹出堆顶最大元素，直到满足 limit 或所有流消耗完毕
	for h.Len() > 0 && len(result) < limit {
		top := heap.Pop(h).(*heapNode)
		item := top.item

		// 去重保护
		if !seen[item.ArticleID] {
			seen[item.ArticleID] = true
			result = append(result, item)
		}

		// 将该流的下一条记录压入堆中
		nextIdx := top.itemIdx + 1
		stream := streams[top.streamIdx]
		if nextIdx < len(stream) && stream[nextIdx] != nil {
			heap.Push(h, &heapNode{
				streamIdx: top.streamIdx,
				itemIdx:   nextIdx,
				item:      stream[nextIdx],
			})
		}
	}

	return result
}
