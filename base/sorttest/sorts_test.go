package sorttest

import (
	"testing"
	"sort"
	"time"
	"fmt"
)

const (
	length = 1000000
	randSeed = 47
	sortIndex = Active
)

func BenchmarkSort(b *testing.B) {
	profileSlice := New(length, randSeed, sortIndex)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		sort.Sort(profileSlice)
	}
	b.StopTimer()
	for _, m := range profileSlice.Profiles {
		fmt.Println(m.Active)
	}
}

func TestSorts(t *testing.T)  {
	profileSlice := New(length, randSeed, sortIndex)
	start := time.Now().UnixNano()
	for i := 0; i < 100; i++ {
		sort.Sort(profileSlice)
		end := time.Now().UnixNano()
		fmt.Printf("testsorts time cost is %d\n", end - start)
		start = end
	}
}

func TestSortsEach(t *testing.T)  {
	for i := 0; i < 100; i++ {
		profileSlice := New(length, randSeed, sortIndex)
		start := time.Now().UnixNano()
		sort.Sort(profileSlice)
		end := time.Now().UnixNano()
		fmt.Printf("testsortseach time cost is %d\n", end - start)
	}
}