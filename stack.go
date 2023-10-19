package parser

type Stack []*IHtmlTag

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(node *IHtmlTag) {
	*s = append(*s, node)
}

func (s *Stack) Len() int {
	return len(*s)
}
func (s *Stack) Next(cur int) *IHtmlTag {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[cur]
}

func (s *Stack) Back(step int) *IHtmlTag {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1-step]
}

func (s *Stack) Peek() *IHtmlTag {
	if len(*s) == 0 {
		return nil
	}
	return (*s)[len(*s)-1]
}

func (s *Stack) Copy() *Stack {
	ret := NewStack()
	*ret = append(*ret, *s...)
	return ret
}

func (s *Stack) Pop() *IHtmlTag {
	if len(*s) == 0 {
		return nil
	}
	index := len(*s) - 1
	node := (*s)[index]
	*s = (*s)[:index]
	return node
}
