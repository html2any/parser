package parser

import (
	"bytes"
	"errors"
	"fmt"
)

var ErrHtmlClose error = errors.New("invalid Html Close Error")
var ErrHtmlAttrError error = errors.New("invalid Html Attr Error")
var ErrHtmlTagMismatch error = errors.New("invalid Html Tag Mismatch")

type IHtmlTag interface {
	NewTag() IHtmlTag
	SetTagName(name string) IHtmlTag
	GetTagName() string
	SetContent(content string) IHtmlTag
	AddChild(child IHtmlTag) IHtmlTag
	SetAttr(attr string, value string) IHtmlTag
}

func findAllAttrs(data *[]byte, attrs *map[string]string, cur *int) (self_close bool, err error) {
	// attrs = make(map[string]string)
	last_key_start := *cur
	last_key_end := -1
	last_val_start := -1
	for *cur < len(*data) {
		switch (*data)[*cur] {
		case '=':
			last_key_end = *cur
		case '"':
			if last_val_start == -1 {
				last_val_start = *cur + 1
			} else {
				last_val_end := *cur
				key := string((*data)[last_key_start:last_key_end])
				val := string((*data)[last_val_start:last_val_end])
				last_key_start = *cur + 1
				last_val_start = -1
				last_key_end = -1
				(*attrs)[key] = val
			}
		case ' ', '\t', '\n', '\r':
			if last_key_end == -1 {
				last_key_start = *cur + 1
			}
		case '>':
			if (*data)[*cur-1] == '/' {
				self_close = true
			}
			return
		}
		*cur++
	}
	err = ErrHtmlAttrError
	return
}

func sprintLineAndPos(head string, data []byte, pos int) string {
	tt := 0
	for line, d := range bytes.Split(data, []byte("\n")) {
		c := len(d) + 1
		if tt+c > pos {
			return fmt.Sprint("|", head, "@[BPOS:", pos, ",LINE:", line+1, ",POS:", pos-tt+1, "]")
		}
		tt = tt + c
	}
	return ""
}

func findTag(data *[]byte, pos *int) (tag string, err error) {
	tag_start := -1
	for *pos < len(*data) {
		c := (*data)[*pos]
		if tag_start == -1 {
			if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
				tag_start = *pos
			}
		} else if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			tag = string((*data)[tag_start:*pos])
			return
		}
		*pos++
	}
	return "", ErrHtmlTagMismatch
}

func getContent(data *[]byte, start, end int) []byte {
	content_start := -1
	content_end := -1
	if start == -1 || end == -1 || start > end {
		return nil
	}
	if start == end {
		ct := (*data)[start : end+1]
		if ct[0] == ' ' || ct[0] == '\t' || ct[0] == '\n' || ct[0] == '\r' {
			return nil
		} else {
			return ct
		}
	}
	for i := start; i < end; i++ {
		c := (*data)[i]
		if content_start == -1 {
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				content_start = i
				break
			}
		}
	}
	left := start
	if content_start != -1 {
		left = content_start
	}
	for i := end; i >= left; i-- {
		c := (*data)[i]
		if content_end == -1 {
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			} else {
				content_end = i + 1
				break
			}
		}
	}
	if content_start != -1 && content_end != -1 {
		return (*data)[content_start:content_end]
	} else {
		return nil
	}
}

func ParseHTML(data []byte, root IHtmlTag) error {
	stack := NewStack()
	tt_len := len(data)
	content_start := -1
	for cur := 0; cur < tt_len; cur++ {
		if data[cur] == '<' { //look for <
			content_end := cur - 1
			cur++
			if cur < tt_len && data[cur] == '/' { //look for </
				cur++
				if tag_name, err := findTag(&data, &cur); err != nil { //look for tag
					return err
				} else { //Find TagClose
					//ReadContent
					ct := getContent(&data, content_start, content_end)
					content_start = cur + 1
					ptop := stack.Pop()
					top := *ptop
					if len(ct) > 0 {
						text_tag := top.NewTag().SetTagName("span")
						text_tag.SetContent(string(ct))
						top.AddChild(text_tag)
					}
					if top.GetTagName() != tag_name {
						return ErrHtmlTagMismatch
					} else if parent := stack.Peek(); parent != nil {
						(*parent).AddChild(top)
					} else {
						return nil
					}
				}
			} else if tag_name, err := findTag(&data, &cur); err != nil { //look for tag
				return err
			} else {
				attrs := make(map[string]string)
				if self_close, err := findAllAttrs(&data, &attrs, &cur); err != nil { //look for attr
					return err
				} else if ptop := stack.Peek(); ptop == nil { //first tag is root
					root.SetTagName(tag_name)
					for k, v := range attrs {
						root.SetAttr(k, v)
					}
					stack.Push(&root)
					ct := getContent(&data, content_start, content_end)
					content_start = cur + 1
					if len(ct) > 0 {
						text_tag := (root).NewTag().SetTagName("span")
						text_tag.SetContent(string(ct))
						(root).AddChild(text_tag)
					}
				} else {
					ct := getContent(&data, content_start, content_end)
					content_start = cur + 1
					if len(ct) > 0 {
						text_tag := (*ptop).NewTag().SetTagName("span")
						text_tag.SetContent(string(ct))
						if parent := stack.Peek(); parent != nil {
							(*parent).AddChild(text_tag)
						}
					}

					sub_tag := (*ptop).NewTag().SetTagName(tag_name)
					for k, v := range attrs {
						sub_tag.SetAttr(k, v)
					}

					if self_close {
						if parent := stack.Peek(); parent != nil {
							(*parent).AddChild(sub_tag)
						}
					} else {
						stack.Push(&sub_tag)
					}
				}
			}
		}
	}
	return ErrHtmlClose
}

// func ParseHTML2(data []byte, root IHtmlTag) error {
// 	stack := NewStack()
// 	z := html.NewTokenizer(bytes.NewReader(data))
// 	for {
// 		tt := z.Next()
// 		switch tt {
// 		case html.ErrorToken:
// 			return z.Err()
// 		case html.TextToken:
// 			d := z.Text()
// 			ct := getContent(&d, 0, len(d)-1)
// 			ptop := stack.Peek()
// 			top := *ptop
// 			if len(ct) > 0 {
// 				text_tag := top.NewTag().SetTagName("span")
// 				text_tag.SetContent(string(ct))
// 				top.AddChild(text_tag)
// 			}
// 		case html.StartTagToken:
// 			tname, hasAttr := z.TagName()
// 			tag_name := string(tname)
// 			if ptop := stack.Peek(); ptop == nil { //first tag is root
// 				root.SetTagName(tag_name)
// 				if hasAttr {
// 					for {
// 						k, v, more := z.TagAttr()
// 						root.SetAttr(string(k), string(v))
// 						if !more {
// 							break
// 						}
// 					}
// 				}
// 				stack.Push(&root)
// 			} else {
// 				sub_tag := (*ptop).NewTag().SetTagName(tag_name)
// 				if hasAttr {
// 					for {
// 						k, v, more := z.TagAttr()
// 						sub_tag.SetAttr(string(k), string(v))
// 						if !more {
// 							break
// 						}
// 					}
// 				}
// 				stack.Push(&sub_tag)
// 			}
// 		case html.EndTagToken:
// 			ptop := stack.Pop()
// 			top := *ptop
// 			tag_name, _ := z.TagName()
// 			if top.GetTagName() != string(tag_name) {
// 				return ErrHtmlTagMismatch
// 			} else if parent := stack.Peek(); parent != nil {
// 				(*parent).AddChild(top)
// 			} else {
// 				return nil
// 			}
// 		case html.SelfClosingTagToken:
// 			ptop := stack.Peek()
// 			tname, hasAttr := z.TagName()
// 			tag_name := string(tname)
// 			sub_tag := (*ptop).NewTag().SetTagName(tag_name)
// 			if hasAttr {
// 				for {
// 					k, v, more := z.TagAttr()
// 					sub_tag.SetAttr(string(k), string(v))
// 					if !more {
// 						break
// 					}
// 				}
// 			}

// 			if parent := stack.Peek(); parent != nil {
// 				(*parent).AddChild(sub_tag)
// 			}
// 		}
// 	}
// }
