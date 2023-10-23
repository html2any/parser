package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/html2any/parser"
	"github.com/html2any/parser/html2json"
)

func BenchmarkHTMLParser(b *testing.B) {
	data, err := os.ReadFile("complicated.html")
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		var root html2json.Tag
		if err := parser.ParseHTML(data, &root); err != nil {
			panic(err)
		}
	}
}

// func BenchmarkHTMLParserByNetHTML(b *testing.B) {
// 	data, err := os.ReadFile("complicated.html")
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		var root html2json.Tag
// 		if err := parser.ParseHTML2(data, &root); err != nil {
// 			panic(err)
// 		}
// 	}
// }

func BenchmarkJSONParser(b *testing.B) {
	data, err := os.ReadFile("complicated.json")
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		var root html2json.Tag
		if err := json.Unmarshal(data, &root); err != nil {
			panic(err)
		}
	}
}

// func BenchmarkSonicParser(b *testing.B) {
// 	data, err := os.ReadFile("complicated.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		var root html2json.Tag
// 		if err := sonic.Unmarshal(data, &root); err != nil {
// 			panic(err)
// 		}
// 	}
// }
