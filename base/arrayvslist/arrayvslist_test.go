package arrayvslist

import (
	"testing"
)

var lenght = 100000

func BenchmarkDigArray1(b *testing.B) {
	arr := MakeArray(lenght, 47)
	//fmt.Printf("%v", arr)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		DigArray1(arr)
	}

}

func BenchmarkDigArray2(b *testing.B) {
	arr := MakeArray(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		DigArray2(arr)
	}
}

func BenchmarkDigList(b *testing.B) {
	l := MakeList(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		DigList(l)
	}
}

func BenchmarkDigMap1(b *testing.B) {
	m := MakeMap1(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		DigMap1(m)
	}
}

func BenchmarkDigMap2(b *testing.B) {
	m := MakeMap2(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		DigMap2(m)
	}
}