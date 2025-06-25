package main

import "testing"

// 寻找最长的不含有重复字符的子串
func searchOnce(s string) string {
	pre := 0
	res := ""
	firIdxMap := make(map[int32]int)
	for i, e := range s {
		if _, ok := firIdxMap[e]; !ok {
			firIdxMap[e] = i
		}
		firIdx := firIdxMap[e]
		if firIdx < i {
			if len(res) < i-pre {
				res = s[pre:i]
			}
			firIdxMap[e] = i
			pre = firIdx + 1
		}
	}
	if len(res) < len(s)-pre {
		res = s[pre:]
	}
	return res
}

func TestSearchOnce(t *testing.T) {
	tests := []struct {
		s   string
		ans string
	}{
		// Normal cases
		{"abcabcbb", "abc"},
		{"pwwkew", "wke"},
		// Edge cases
		{"", ""},
		{"b", "b"},
		{"bbbbbbbb", "b"},
		{"abcabcabcd", "abcd"},
	}

	for _, tt := range tests {
		actual := searchOnce(tt.s)
		if actual != tt.ans {
			t.Errorf("searchOnce(%q)=%q, want %q", tt.s, actual, tt.ans)
		}
	}
}

func BenchmarkSearchOnce(b *testing.B) {
	s := "abcabcabcd"
	ans := "abcd"

	for i := 0; i < b.N; i++ {
		actual := searchOnce(s)
		if actual != ans {
			b.Errorf("searchOnce(%q)=%q, want %q", s, actual, ans)
		}
	}
}
