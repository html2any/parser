package html2json

import (
	"encoding/json"
	"parser"
)

func Convert(data []byte) ([]byte, error) {
	var root Tag
	if err := parser.ParseHTML(data, &root); err != nil {
		return []byte(""), err
	}
	return json.Marshal(root)
}
