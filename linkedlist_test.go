package linkedlist

import (
	"container/list"
	"fmt"
	"testing"

	intlist "github.com/itsmontoya/linkedlist/typed/int"
	"time"
)

var (
	testFilterVal    []GenericVal
	testFilterIntVal []int
)

func TestLinkedList(t *testing.T) {
	var (
		l   LinkedList
		err error
	)

	l.Append(0, 1, 2, 3, 4, 5, 6)
	if l.Len() != 7 {
		t.Fatalf("invalid length, expected %v and received %v", 7, l.Len())
	}

	if err = testIteration(&l, 0); err != nil {
		t.Fatal(err)
	}

	if err = testMap(&l, 0); err != nil {
		t.Fatal(err)
	}

	if err = testFilter(&l, 0, true); err != nil {
		t.Fatal(err)
	}

	if err = testReduce(&l, 0); err != nil {
		t.Fatal(err)
	}

	l.ForEach(nil, func(n *Node, _ GenericVal) bool {
		// Call a new goroutine to remove Node
		// Node: If this is not a goroutine, it will be a deadlock
		go l.Remove(n)
		return false
	})

	// Give time for goroutines to execute
	time.Sleep(time.Millisecond * 10)

	// Ensure that all the nodes were properly removed
	if l.Len() != 0 {
		t.Fatalf("invalid length, expected %v and received %v", 0, l.Len())
	}

	return
}

func TestMapFilterReduce(t *testing.T) {
	var l LinkedList
	l.Append(0, 1, 2, 3, 4, 5, 6)

	val := l.Map(testAddOne).Filter(testIsEven).Reduce(testAddInts)
	if val != 12 {
		t.Fatalf("expected %v and received %v", 12, val)
	}
}

func testIteration(l *LinkedList, start int) (err error) {
	cnt := start

	l.ForEach(nil, func(_ *Node, val GenericVal) bool {
		if val.(int) != cnt {
			err = fmt.Errorf("invalid value, expected %d and received %d", cnt, val)
			return true
		}

		cnt++
		return false
	})

	cnt--

	l.ForEachRev(nil, func(_ *Node, val GenericVal) bool {
		if val.(int) != cnt {
			err = fmt.Errorf("invalid value, expected %d and received %d", cnt, val)
			return true
		}

		cnt--
		return false
	})

	return
}

func testMap(l *LinkedList, start int) (err error) {
	list := l.Map(func(val GenericVal) (nval GenericVal) {
		nval = val.(int) * 2
		return
	}).Slice()

	for i := 0; i < len(list); i++ {
		v := list[i]
		ev := (i + start) * 2
		if v != ev {
			return fmt.Errorf("invalid value, expected %d and received %d", ev, v)
		}
	}

	return
}

func testFilter(l *LinkedList, tgt int, expected bool) (err error) {
	list := l.Filter(func(val GenericVal) (ok bool) {
		return val.(int) == tgt
	}).Slice()

	expectedLen := 1
	if !expected {
		expectedLen = 0
	}

	if ll := len(list); ll != expectedLen {
		err = fmt.Errorf("invalid list length, expected %d and received %d", expectedLen, ll)
	}

	return
}

func testReduce(l *LinkedList, start int) (err error) {
	var cv int
	len := int(l.Len())
	val := l.Reduce(func(acc, val GenericVal) (sum GenericSum) {
		accV, _ := acc.(int)
		sum = accV + val.(int)
		return
	}).(int)

	for i := start; i < len+start; i++ {
		cv += i
	}

	if val != cv {
		err = fmt.Errorf("invalid value, expected %d and received %d", cv, val)
	}

	return
}

func testAddOne(val GenericVal) (nval GenericVal) {
	nval = val.(int) + 1
	return
}

func testIsEven(val GenericVal) (ok bool) {
	return val.(int)%2 == 0
}

func testAddInts(acc, val GenericVal) (sum GenericSum) {
	accV, _ := acc.(int)
	sum = accV + val.(int)
	return
}

func BenchmarkListAppend(b *testing.B) {
	var l LinkedList
	for i := 0; i < b.N; i++ {
		l.Append(i)
	}

	b.ReportAllocs()
}

func BenchmarkListFilter(b *testing.B) {
	var l LinkedList
	for i := 0; i < b.N; i++ {
		l.Append(i)
	}
	b.ResetTimer()

	testFilterVal = l.Filter(func(val GenericVal) bool {
		return val.(int)%2 == 0
	}).Slice()

	b.ReportAllocs()
}

func BenchmarkIntListAppend(b *testing.B) {
	var l intlist.LinkedList
	for i := 0; i < b.N; i++ {
		l.Append(i)
	}

	b.ReportAllocs()
}

func BenchmarkIntListFilter(b *testing.B) {
	var l intlist.LinkedList
	for i := 0; i < b.N; i++ {
		l.Append(i)
	}
	b.ResetTimer()

	testFilterIntVal = l.Filter(func(val int) bool {
		return val%2 == 0
	}).Slice()

	b.ReportAllocs()
}

func BenchmarkStdListAppend(b *testing.B) {
	var l list.List
	for i := 0; i < b.N; i++ {
		l.PushBack(i)
	}

	b.ReportAllocs()
}

func BenchmarkSliceAppend(b *testing.B) {
	s := make([]GenericVal, 0, 32)
	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}

	b.ReportAllocs()
}

func BenchmarkMapAppend(b *testing.B) {
	s := make(map[int]GenericVal, 32)
	for i := 0; i < b.N; i++ {
		s[i] = i
	}

	b.ReportAllocs()
}

func BenchmarkListPrepend(b *testing.B) {
	var l LinkedList
	for i := 0; i < b.N; i++ {
		l.Prepend(i)
	}

	b.ReportAllocs()
}

func BenchmarkIntListPrepend(b *testing.B) {
	var l intlist.LinkedList
	for i := 0; i < b.N; i++ {
		l.Prepend(i)
	}

	b.ReportAllocs()
}

func BenchmarkStdListPrepend(b *testing.B) {
	var l list.List
	for i := 0; i < b.N; i++ {
		l.PushFront(i)
	}

	b.ReportAllocs()
}

func BenchmarkSlicePrepend(b *testing.B) {
	s := make([]GenericVal, 0, 32)
	for i := 0; i < b.N; i++ {
		s = append([]GenericVal{i}, s...)
	}

	b.ReportAllocs()
}

func BenchmarkSliceFilter(b *testing.B) {
	s := make([]GenericVal, 0, b.N)
	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}
	b.ResetTimer()

	var ns []GenericVal
	for _, val := range s {
		if val.(int)%2 == 0 {
			ns = append(ns, val)
		}
	}

	testFilterVal = ns
	b.ReportAllocs()
}

func BenchmarkMapPrepend(b *testing.B) {
	s := make(map[int]GenericVal, 32)
	for i := 0; i < b.N; i++ {
		s[i] = i
	}

	b.ReportAllocs()
}

func BenchmarkMapFilter(b *testing.B) {
	m := make(map[int]GenericVal, b.N)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
	b.ResetTimer()

	var ns []GenericVal
	for _, val := range m {
		if val.(int)%2 == 0 {
			ns = append(ns, val)
		}
	}

	testFilterVal = ns
	b.ReportAllocs()
}
