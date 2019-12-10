package main

import (
	"fmt"
	"testing"
)

func BenchmarkGoogleStarlark(b *testing.B) {
	g := -1
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		g = i
		googleStarlark()
	}
	b.StopTimer()
	fmt.Println(g)
}

func BenchmarkStarlight(b *testing.B) {
	g := -1
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		g = i
		light()
	}
	b.StopTimer()
	fmt.Println(g)
}
