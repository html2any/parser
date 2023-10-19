package html2json

import "parser"

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
