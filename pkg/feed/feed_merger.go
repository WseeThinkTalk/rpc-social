package feed

import (
	"sort"
)

type FeedMerger struct{}

func NewFeedMerger() *FeedMerger {
	return &FeedMerger{}
}

// MergeMultiStreams 将用户自己的 Inbox 与多位大 V 的 Outbox 流执行倒序去重归并截断
func (m *FeedMerger) MergeMultiStreams(streams [][]*FeedItem, limit int) []*FeedItem {
	if len(streams) == 0 || limit <= 0 {
		return []*FeedItem{}
	}

	var total []*FeedItem
	seen := make(map[int64]bool)

	for _, stream := range streams {
		for _, item := range stream {
			if item != nil && !seen[item.ArticleID] {
				seen[item.ArticleID] = true
				total = append(total, item)
			}
		}
	}

	sort.SliceStable(total, func(i, j int) bool {
		return total[i].PublishTime > total[j].PublishTime
	})

	if len(total) > limit {
		total = total[:limit]
	}

	return total
}
