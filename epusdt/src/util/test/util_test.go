package test

import (
	"fmt"
	"sort"
	"testing"
)

func Test_hash_count(t *testing.T) {
	// 测试用例
	hash := "5aa74d3c8044858d439c07003788529f1bb6e9b7ace101cb588f8d2a3e50df3b"

	// 创建一个map来存储字符出现次数
	charCount := make(map[rune]int)

	// 统计每个字符出现的次数
	for _, char := range hash {
		charCount[char]++
	}

	// 创建一个结构体切片来存储字符和其出现次数
	type CharFreq struct {
		Char  rune
		Count int
	}

	// 将map转换为结构体切片
	var freqList []CharFreq
	for char, count := range charCount {
		freqList = append(freqList, CharFreq{char, count})
	}

	// 按照出现次数从大到小排序
	sort.Slice(freqList, func(i, j int) bool {
		if freqList[i].Count == freqList[j].Count {
			// 如果出现次数相同，按字符顺序排序
			return freqList[i].Char < freqList[j].Char
		}
		return freqList[i].Count > freqList[j].Count
	})

	// 打印结果
	fmt.Printf("Hash string: %s\n\n", hash)
	fmt.Println("Character frequency (from high to low):")
	fmt.Println("Character | Count")
	fmt.Println("----------|-------")

	for _, cf := range freqList {
		fmt.Printf("    %c    |   %d\n", cf.Char, cf.Count)
	}
}
