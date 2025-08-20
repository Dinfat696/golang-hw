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
	if l == nil {
		panic("list is nil")
	}
	litem := &ListItem{Value: v, Prev: nil, Next: l.front}
	if l.front != nil {
		l.front.Prev = litem
	}
	if l.length == 0 {
		l.back = litem
	}
	l.length++
	l.front = litem
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l == nil {
		panic("list is nil")
	}
	litem := &ListItem{Value: v, Prev: l.back, Next: nil}
	if l.back != nil {
		l.back.Next = litem
	}
	if l.length == 0 {
		l.front = litem
	}
	l.length++
	l.back = litem
	return l.back
}

func (l *list) Remove(item *ListItem) {
	if l == nil {
		panic("list is nil")
	}
	if item != nil && l.length > 0 {
		switch {
		case l.length == 1:
			l.front = nil
			l.back = nil
		case item == l.front:
			item.Next.Prev = nil
			l.front = item.Next
		case item == l.back:
			item.Prev.Next = nil
			l.back = item.Prev
		default:
			item.Prev.Next = item.Next
			item.Next.Prev = item.Prev
		}
		l.length--
		item = nil
	}
}

func (l *list) MoveToFront(item *ListItem) {
	if l == nil {
		panic("list is nil")
	}
	if item != nil && l.length > 1 && item != l.front {
		if item == l.back {
			l.back = item.Prev
			item.Prev.Next = nil
		} else {
			item.Prev.Next = item.Next
			item.Next.Prev = item.Prev
		}
		l.front.Prev = item
		item.Prev = nil
		item.Next = l.front
		l.front = item
	}
}

func NewList() List {
	return new(list)
}
