package main

import (
	"os"

	"github.com/html2any/parser/html2json"
)

func main() {
	if data, err := os.ReadFile("simple.html"); err != nil {
		panic(err)
	} else {
		if data, err := html2json.Convert(data); err != nil {
			panic(err)
		} else {
			os.WriteFile("simple.json", data, 0644)
		}
	}

	if data, err := os.ReadFile("complicated.html"); err != nil {
		panic(err)
	} else {
		if data, err := html2json.Convert(data); err != nil {
			panic(err)
		} else {
			os.WriteFile("complicated.json", data, 0644)
		}
	}
}
