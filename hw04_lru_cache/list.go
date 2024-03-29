package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v, Prev: nil, Next: l.front}
	if l.length != 0 {
		l.front.Prev = &item
		l.front = &item
	} else {
		l.front = &item
		l.back = &item
	}

	l.length++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v, Prev: l.back, Next: nil}
	if l.length != 0 {
		l.back.Next = &item
		l.back = &item
	} else {
		l.front = &item
		l.back = &item
	}

	l.length++
	return &item
}

func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
		if l.length > 1 {
			l.back.Next = nil
		}
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
		if l.length > 1 {
			l.front.Prev = nil
		}
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.front {
		l.Remove(i)
		item := l.PushFront(i.Value)
		*i = *item
	}
}
