package CustomUtils

import (
	"fmt"
	"testing"
)

func TestFixedDequeInsert(t *testing.T) {
	q := New(100)
	var i = 0
	for i = 0; i < 90; i++ {
		q.Insert(i)
	}
	for i = 90; i < 120; i++ {
		q.Insert(i)
		fmt.Printf("begining %v \n", q.Back())
	}
}

func TestStringInsert(t *testing.T) {
	q := New(100)
	s := "Hello world"
	bits := []byte(s)
	for _, b := range bits {
		q.Insert(b)
	}

	s2 := string(q.Bytes())

	fmt.Println(s2)
}
