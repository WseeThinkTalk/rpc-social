package sensitive

import (
	"testing"
)

// TestSensitiveFilter 验证敏感词检测准确性
func TestSensitiveFilter(t *testing.T) {
	words := []string{"暴力", "涉黄", "赌博"}
	f := NewFilter(words)

	cases := []struct {
		input     string
		sensitive bool
	}{
		{"这是一条友善正常的评论", false},
		{"评论包含涉黄词汇", true},
		{"涉*黄测试", true},
		{"赌  博 网站", true},
	}

	for _, v := range cases {
		got := f.IsSensitive(v.input)
		if got != v.sensitive {
			t.Errorf("输入 [%s] 预期为 %v, 实际为 %v", v.input, v.sensitive, got)
		}
	}
}
