package gmh

import (
	"container/list"
)

type stack struct {
	list *list.List
}

func newStack() *stack {
	return &stack{
		list: list.New(),
	}
}

func (this *stack) push(val interface{}) *stack {
	this.list.PushFront(val)
	return this
}

func (this *stack) pop() interface{} {
	ele := this.list.Front()
	if ele == nil {
		return nil
	}

	val := ele.Value
	this.list.Remove(ele)

	return val
}

func (this *stack) clear() *stack {
	this.list.Init()
	return this
}
