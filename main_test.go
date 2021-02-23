package main

import (
	"embed"
	"io/ioutil"
	"testing"
)

//go:embed dist/index.html
var big []byte

//go:embed dist
var fs embed.FS

func Benchmark_embed(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_ = big
	}
}

func Benchmark_fs_dir(t *testing.B) {
	for i := 0; i < t.N; i++ {
		fs.ReadFile("dist/index.html")
	}
}

func Benchmark_fs_file(t *testing.B) {
	for i := 0; i < t.N; i++ {
		fs.ReadFile("data/index.html")
	}
}

func Benchmark_read_file(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ioutil.ReadFile("data/index.html")
	}
}
