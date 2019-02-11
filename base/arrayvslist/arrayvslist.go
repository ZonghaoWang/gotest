package arrayvslist

import (
	"math/rand"
	"strconv"
)

type ListNode struct {
	value int
	next  *ListNode
}

func MakeMap1(length int, seed int64) (result map[int]int)  {
	result = make(map[int]int, length)
	rand.Seed(seed)
	for index :=0; index < length; index++ {
		result[index] = rand.Intn(10000)
	}
	return
}

func MakeMap2(length int, seed int64) (result map[string]int)  {
	result = make(map[string]int, length)
	rand.Seed(seed)
	for index :=0; index < length; index++ {
		result[strconv.Itoa(index)] = rand.Intn(10000)
	}
	return
}

func MakeList(length int, seed int64) *ListNode {
	rand.Seed(seed)
	head := ListNode{
		value: 0,
	}
	now := &head
	for index := 0; index < length; index++ {
		tmp := &ListNode{
			value: rand.Intn(10000),
		}
		now.next = tmp
		now = tmp
	}
	return &head
}

func MakeArray(length int, seed int64) []int {
	rand.Seed(seed)
	result := make([]int, length, length)
	for index, _ := range result {
		result[index] = rand.Intn(10000)
	}
	return result
}

func DigArray1(arr []int) int {
	var result int
	for _, value := range arr {
		result += value
	}
	return result
}


func DigArray2(arr []int) int {
	var result int
	length := len(arr)
	for index := 0; index < length; index++ {
		result += arr[index]
	}
	return result
}

func DigList(l *ListNode) int {
	now := l
	var result int
	for {
		if now.next == nil {
			break
		}
		result += now.value
		now = now.next
	}
	return result
}

func DigMap1(m map[int]int) int {
	var result int
	for _, v := range m {
		result += v
	}
	return result
}

func DigMap2(m map[string]int) int {
	var result int
	for _, v := range m {
		result += v
	}
	return result
}