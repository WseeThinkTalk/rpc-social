package feed

// FeedItem 单条动态元数据
type FeedItem struct {
	ArticleID   int64 `json:"articleId"`
	AuthorID    int64 `json:"authorId"`
	PublishTime int64 `json:"publishTime"`
}

const (
	BigVFollowerThreshold = 5000 // 大 V 粉丝数分水岭
	MaxInboxCapacity      = 800  // 粉丝收件箱最大容量
)
