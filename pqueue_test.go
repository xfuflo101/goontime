package goontime_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	got "github.com/xfuflo101/goontime"
)

type Elem struct {
	i int64
	s string
}

func NewElem(i int64, s string) *Elem {
	return &Elem{i: i, s: s}
}

func (self *Elem) String() string {
	return fmt.Sprintf("[%v : %v]", self.i, self.s)
}

type BiComparatorAscending struct{}

func (BiComparatorAscending) Less(lhs, rhs *Elem) bool {
	return lhs.i < rhs.i
}

type BiComparatorDescending struct{}

func (BiComparatorDescending) Less(lhs, rhs *Elem) bool {
	return rhs.i < lhs.i
}

var cmpAscending BiComparatorAscending
var cmpDescending BiComparatorDescending

func dumpSlice[T any](pq *got.PriorityQueue[T]) (slc []T) {
	for pq != nil && !pq.Empty() {
		slc = append(slc, pq.Pop())
	}
	return
}

func dumpE2ESlice[T any](pq *got.PriorityQueueE2E[T]) (slc []T) {
	for pq != nil && !pq.Empty() {
		slc = append(slc, pq.Pop())
	}
	return
}

func toStrSlice[T any](slc []T) string {
	return fmt.Sprintf("%v", slc)
}

func TestPQueue01_Ascending(t *testing.T) {

	t.Logf("TestPQueue01 - started")

	pq := got.NewPriorityQueueCap[*Elem](cmpAscending, 5)

	var result string

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[]", result, "TestPQueue01")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(8, "8"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(6, "6"))
	pq.Push(NewElem(10, "10"))
	pq.Push(NewElem(6, "6_2"))
	pq.Push(NewElem(3, "3"))
	pq.Push(NewElem(8, "8_2"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(6, "6_3"))

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_2] [3 : 3] [5 : 5] [6 : 6] [6 : 6_3] [6 : 6_2] [8 : 8] [8 : 8_2] [10 : 10] [15 : 15]]", result, "TestPQueue01")

	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(1, "1_3"))
	pq.Push(NewElem(1, "1_4"))
	pq.Push(NewElem(1, "1_5"))
	pq.Push(NewElem(1, "1_6"))
	pq.Push(NewElem(1, "1_7"))
	pq.Push(NewElem(1, "1_8"))
	pq.Push(NewElem(1, "1_9"))
	pq.Push(NewElem(1, "1_10"))

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_10] [1 : 1_9] [1 : 1_8] [1 : 1_7] [1 : 1_6] [1 : 1_5] [1 : 1_4] [1 : 1_3] [1 : 1_2]]", result, "TestPQueue01")

	t.Logf("TestPQueue01 - OK")
}

func TestPQueue02_Descending(t *testing.T) {

	t.Logf("TestPQueue02 - started")

	pq := got.NewPriorityQueueCap[*Elem](cmpDescending, 5)

	var result string

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[]", result, "TestPQueue02")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(8, "8"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(6, "6"))
	pq.Push(NewElem(10, "10"))
	pq.Push(NewElem(6, "6_2"))
	pq.Push(NewElem(3, "3"))
	pq.Push(NewElem(8, "8_2"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(6, "6_3"))

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[[15 : 15] [10 : 10] [8 : 8] [8 : 8_2] [6 : 6] [6 : 6_2] [6 : 6_3] [5 : 5] [3 : 3] [1 : 1_2] [1 : 1]]", result, "TestPQueue02")

	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(1, "1_3"))
	pq.Push(NewElem(1, "1_4"))
	pq.Push(NewElem(1, "1_5"))
	pq.Push(NewElem(1, "1_6"))
	pq.Push(NewElem(1, "1_7"))
	pq.Push(NewElem(1, "1_8"))
	pq.Push(NewElem(1, "1_9"))
	pq.Push(NewElem(1, "1_10"))

	result = toStrSlice(dumpSlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_10] [1 : 1_9] [1 : 1_8] [1 : 1_7] [1 : 1_6] [1 : 1_5] [1 : 1_4] [1 : 1_3] [1 : 1_2]]", result, "TestPQueue01")

	t.Logf("TestPQueue02 - OK")
}

func TestPQueue03_Ascending_End2End(t *testing.T) {

	t.Logf("TestPQueue03 - started")

	pq := got.NewPriorityQueueE2ECap[*Elem](cmpAscending, 5)

	var result string

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[]", result, "TestPQueue03")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(10, "10"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[5 : 5] [10 : 10] [15 : 15]]", result, "TestPQueue03")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(8, "8"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(6, "6"))
	pq.Push(NewElem(10, "10"))
	pq.Push(NewElem(6, "6_2"))
	pq.Push(NewElem(3, "3"))
	pq.Push(NewElem(8, "8_2"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(6, "6_3"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_2] [3 : 3] [5 : 5] [6 : 6] [6 : 6_2] [6 : 6_3] [8 : 8] [8 : 8_2] [10 : 10] [15 : 15]]", result, "TestPQueue03")

	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(1, "1_3"))
	pq.Push(NewElem(1, "1_4"))
	pq.Push(NewElem(1, "1_5"))
	pq.Push(NewElem(1, "1_6"))
	pq.Push(NewElem(1, "1_7"))
	pq.Push(NewElem(1, "1_8"))
	pq.Push(NewElem(1, "1_9"))
	pq.Push(NewElem(1, "1_10"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_2] [1 : 1_3] [1 : 1_4] [1 : 1_5] [1 : 1_6] [1 : 1_7] [1 : 1_8] [1 : 1_9] [1 : 1_10]]", result, "TestPQueue03")

	t.Logf("TestPQueue03 - OK")
}

func TestPQueue04_Descending_End2End(t *testing.T) {

	t.Logf("TestPQueue04 - started")

	pq := got.NewPriorityQueueE2ECap[*Elem](cmpDescending, 5)

	var result string

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[]", result, "TestPQueue04")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(10, "10"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[15 : 15] [10 : 10] [5 : 5]]", result, "TestPQueue04")

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(8, "8"))
	pq.Push(NewElem(5, "5"))
	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(6, "6"))
	pq.Push(NewElem(10, "10"))
	pq.Push(NewElem(6, "6_2"))
	pq.Push(NewElem(3, "3"))
	pq.Push(NewElem(8, "8_2"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(6, "6_3"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[15 : 15] [10 : 10] [8 : 8] [8 : 8_2] [6 : 6] [6 : 6_2] [6 : 6_3] [5 : 5] [3 : 3] [1 : 1] [1 : 1_2]]", result, "TestPQueue04")

	pq.Push(NewElem(1, "1"))
	pq.Push(NewElem(1, "1_2"))
	pq.Push(NewElem(1, "1_3"))
	pq.Push(NewElem(1, "1_4"))
	pq.Push(NewElem(1, "1_5"))
	pq.Push(NewElem(1, "1_6"))
	pq.Push(NewElem(1, "1_7"))
	pq.Push(NewElem(1, "1_8"))
	pq.Push(NewElem(1, "1_9"))
	pq.Push(NewElem(1, "1_10"))

	result = toStrSlice(dumpE2ESlice(pq))
	assert.Equal(t, "[[1 : 1] [1 : 1_2] [1 : 1_3] [1 : 1_4] [1 : 1_5] [1 : 1_6] [1 : 1_7] [1 : 1_8] [1 : 1_9] [1 : 1_10]]", result, "TestPQueue04")

	t.Logf("TestPQueue04 - OK")
}
