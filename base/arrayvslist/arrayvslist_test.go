package arrayvslist

import (
	"github.com/arl/assertgo"
	"testing"
)

var lenght = 100000

func BenchmarkDigArray1(b *testing.B) {
	arr := MakeArray(lenght, 47)
	//fmt.Printf("%v", arr)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		assert.True(DigArray1(arr) == 499928920144)
	}

}

func BenchmarkDigArray2(b *testing.B) {
	arr := MakeArray(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		assert.True(DigArray2(arr) == 499928920144)
	}
}

func BenchmarkDigList(b *testing.B) {
	l := MakeList(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		assert.True(DigList(l) == 499928920144)
	}
}

func BenchmarkDigMap1(b *testing.B) {
	m := MakeMap1(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		assert.True(DigMap1(m) == 499928920144)
	}
}

func BenchmarkDigMap2(b *testing.B) {
	m := MakeMap2(lenght, 47)
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		assert.True(DigMap2(m) == 499928920144)
	}
}