package CustomUtils

import (
	"sssh_server/CustomUtils/DequeCopy"
)

// Basically what this will do is,
// given a fixed size when you insert if the insert goes past the size
// remove elements from te beginning
// I could use a linked list but that will required space x2
type FixedDeque struct {
	DequeCopy.Queue
	maxSize int
}

func New(size int) *FixedDeque {
	f := new(FixedDeque)
	f.maxSize = size
	f.Init()
	return f
}

func (f *FixedDeque) Insert(i interface{}) {
	if f.Len() > f.maxSize {
		f.PopBack()
	}
	f.PushFront(i)
}

func (f *FixedDeque) InsertMultiple(is []interface{}) {
	for _, interfaceElem := range is {
		f.Insert(interfaceElem)
	}
}
