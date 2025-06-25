package main

import (
	"fmt"
	"strings"
)

// 寻找最长的不含有重复字符的子串 兼容中文
func maxLen(s string) int {
	pre := 0
	maxLen := 0
	firIdxMap := make(map[rune]int)
	for i, ch := range []rune(s) {
		if _, ok := firIdxMap[ch]; !ok {
			firIdxMap[ch] = i
		}
		firIdx := firIdxMap[ch]
		if firIdx < i {
			maxLen = max(maxLen, i-pre)
			firIdxMap[ch] = i
			pre = firIdx + 1
		}
	}

	return max(maxLen, len(s)-pre)
}

func main() {
	m1 := map[string]string{
		"name":    "weipan4",
		"company": "iflytek",
		"job":     "java+LLM+golang",
	}
	fmt.Println(m1)

	for k, v := range m1 {
		fmt.Println(k, v)
	}
	for _, v := range m1 {
		fmt.Println(v)
	}

	m2 := make(map[string]string) // m2 = empty map
	var m3 map[string]int         // m3 = nil 注意Go语言中的nil是可以参与运算的
	fmt.Println(m2, m3)

	fmt.Println("getting values")
	name := m1["name"]
	fmt.Println(name)
	company := m1["command"] //如果key不存在的话 value
	fmt.Println(company)

	if job, ok := m1["job"]; ok {
		fmt.Println(job)
	} else {
		fmt.Println("job not found")
	}

	fmt.Println("delete values")
	name, ok := m1["name"]
	fmt.Println(name, ok)
	delete(m1, "name")
	name, ok = m1["name"]
	fmt.Println(name, ok)

	fmt.Println(maxLen("abcabcabcd"))
	fmt.Println(strings.Fields("abc abc bb"))

	// map赋值传递的是源map的引用
	tag1 := map[string]int{
		"name":    1,
		"job":     2,
		"company": 3,
	}
	fmt.Println(tag1)
	tag2 := tag1
	tag2["name"] = 2
	tag2["job"] = 3
	tag2["company"] = 4
	fmt.Println(tag1)

	// 切片赋值会创建一片新的内存存储源切片的副本
	l1 := []int{1, 2, 3, 4, 5}
	l2 := l1
	l2 = append(l2, 6)
	fmt.Println(l1)
}
