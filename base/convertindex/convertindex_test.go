package convertindex

import (
	"testing"
	"fmt"
)

func BenchmarkSearch1(b *testing.B) {
	sp := SuperProfile{
		memory: &Memory{},
		invertIndex: &InvertIndex{},
	}
	sp.Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("length of result is ,", len(sp.SearchByInfo(Credit)))
	}
}

func BenchmarkSearch2(b *testing.B) {
	sp := &SuperProfile{
		memory: &Memory{},
		invertIndex: &InvertIndex{},
	}
	sp.Init()
	fmt.Printf("%v\n", sp.invertIndex.ActiveIndex[1])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("length of result is ,", len(sp.SearchByInfoIndex(Credit)))
	}
}