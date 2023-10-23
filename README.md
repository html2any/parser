# Parse HTML to Go Struct or Dump to JSON

The **parser** project provides a method to convert HTML bytes to Golang Object or JSON bytes.

## Installation

To install and set up the project, follow these steps:

1. Clone the repository: `git clone https://github.com/html2any/parser.git`
2. Navigate to the project directory: `cd parser/example`
3. Install the required dependencies: `go build && ./example`

## Usage

Follow these instructions to use the project:
### HTML TO JSON
```Go

    if data, err := os.ReadFile("complicated.html"); err != nil {
		panic(err)
	} else {
		if data, err := html2json.Convert(data); err != nil {
			panic(err)
		} else {
			os.WriteFile("complicated.json", data, 0644)
		}
	}
```

### HTML TO Go Object
1. You should define a struct with methods
2. Implement Methods
3. Examples could find in html2json
```Go
type Tag struct {
	TagName  string            `json:"tagName"`
	Attrs    map[string]string `json:"attrs"`
	Children []*Tag            `json:"children"`
	Content  string            `json:"content"`
}

// When Find A New Tag, The Call Step is:
// 0. Check Content Before Tag, If It Is Not Empty, Add To Parent.
// 1. Create A NewTag.
// 2. Set The TagName.
// 3. Set Attrs.
// 4. Add Child To Parent.

// 0. Check Content Before Tag, If It Is Not Empty, Add To Parent.
func (t *Tag) SetContent(content string) parser.IHtmlTag {
	t.Content = content
	return t
}

// 1. Create A NewTag.
func (t *Tag) NewTag() parser.IHtmlTag {
	tag := new(Tag)
	return tag
}

// 2. Set The TagName.
func (t *Tag) SetTagName(tagname string) parser.IHtmlTag {
	t.TagName = tagname
	t.Children = make([]*Tag, 0)
	t.Attrs = make(map[string]string)
	t.TagName = tagname
	return t
}

// 3. Set Attrs
func (t *Tag) SetAttr(attr string, value string) parser.IHtmlTag {
	t.Attrs[attr] = value
	return t
}

// 4. Add Child To Parent
func (t *Tag) AddChild(child parser.IHtmlTag) parser.IHtmlTag {
	t.Children = append(t.Children, child.(*Tag))
	return t
}

// GetTagName Help To Check Tag Close is Correct?
func (t *Tag) GetTagName() string {
	return t.TagName
}
```
4. Use method just like `json.Unmarshal`
```Go
    var root Tag
	if err := parser.ParseHTML(data, &root); err != nil {
		panic(err)
	}
```

## Performance
```
goos: linux
goarch: amd64
pkg: parser/example
cpu: Intel(R) Xeon(R) CPU E5-2620 0 @ 2.00GHz
BenchmarkHTMLParser-24                      6390            988853 ns/op          197689 B/op       4407 allocs/op
BenchmarkHTMLParserByNetHTML-24             4892           1192781 ns/op          200719 B/op       4020 allocs/op
BenchmarkJSONParser-24                      1194           5007065 ns/op          234991 B/op       5561 allocs/op
BenchmarkSonicParser-24                     4796           1198912 ns/op          235447 B/op       2232 allocs/op
PASS
ok      parser/example  19.045s
```

## Contributing

Contributions to the project are welcome! If you would like to contribute, please follow these guidelines:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/your-feature`.
3. Make your changes and commit them: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin feature/your-feature`.
5. Open a pull request.

Please ensure that your code adheres to the project's coding style and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
