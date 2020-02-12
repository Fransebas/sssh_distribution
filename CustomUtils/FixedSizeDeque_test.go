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
		fmt.Printf("begining %v and last %v \n", q.Back(), q.Front())

	}
}

func TestTerminalSim(t *testing.T) {
	var offset = 0
	q := New(10)
	frase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var i = 0
	for i = 0; i < 10; i++ {
		q.Insert(frase[i])
	}
	for i = 10; i < 25; i++ {
		var b []byte
		b, offset = q.BytesFrom(offset)
		q.Insert(frase[i])
		//t.Log()
		fmt.Printf("b %v offset %v \n", string(b), offset)
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
