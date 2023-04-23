package goontime

import (
	"container/heap"
)

////////////////////////////////////////////////////////

type BiComparatorInterface[T any] interface {
	Less(lhs, rhs T) bool
}

////////////////////////////////////////////////////////

// A PriorityQueue
type PriorityQueue[T any] struct {
	impl priorityQueueImpl[T]
}

func NewPriorityQueueCap[T any](cmp BiComparatorInterface[T], cap int) *PriorityQueue[T] {
	result := &PriorityQueue[T]{}
	result.impl.pq = make([]T, 0, cap)
	result.impl.cmp = cmp
	heap.Init(&result.impl)
	return result
}

func NewPriorityQueueSrc[T any](cmp BiComparatorInterface[T], src []T) *PriorityQueue[T] {
	result := &PriorityQueue[T]{}
	result.impl.pq = make([]T, 0, len(src))
	copy(result.impl.pq, src)
	result.impl.cmp = cmp
	heap.Init(&result.impl)
	return result
}

func (self *PriorityQueue[T]) Len() int {
	return self.impl.Len()
}

func (self *PriorityQueue[T]) Empty() bool {
	return self.Len() == 0
}

func (self *PriorityQueue[T]) Push(x T) {
	heap.Push(&self.impl, x)
}

func (self *PriorityQueue[T]) Pop() T {
	return heap.Pop(&self.impl).(T)
}

func (self *PriorityQueue[T]) Peek() T {
	return self.impl.Peek()
}

////////////////////////////////////////////////////////

type e2eEntry[T any] struct {
	orig   T
	e2eIdx uint64
}

type e2eBiComparator[T any] struct {
	orig BiComparatorInterface[T]
}

func (self *e2eBiComparator[T]) Less(lhs, rhs *e2eEntry[T]) bool {
	if self.orig.Less(lhs.orig, rhs.orig) {
		return true
	} else if self.orig.Less(rhs.orig, lhs.orig) {
		return false
	} else {
		return lhs.e2eIdx < rhs.e2eIdx // end-to-end - ascending order forever
	}
}

// A PriorityQueueE2E
type PriorityQueueE2E[T any] struct {
	impl       priorityQueueImpl[*e2eEntry[T]]
	e2eIdxBase uint64
}

func NewPriorityQueueE2ECap[T any](cmp BiComparatorInterface[T], cap int) *PriorityQueueE2E[T] {
	result := &PriorityQueueE2E[T]{}
	result.impl.pq = make([]*e2eEntry[T], 0, cap)
	result.impl.cmp = &e2eBiComparator[T]{orig: cmp}
	heap.Init(&result.impl)
	return result
}

func NewPriorityQueueE2ESrc[T any](cmp BiComparatorInterface[T], src []T) *PriorityQueueE2E[T] {
	result := &PriorityQueueE2E[T]{}
	result.impl.pq = make([]*e2eEntry[T], 0, len(src))
	for _, x := range src {
		result.impl.pq = append(result.impl.pq, result.createE2E(x))
	}
	result.impl.cmp = &e2eBiComparator[T]{orig: cmp}
	heap.Init(&result.impl)
	return result
}

func (self *PriorityQueueE2E[T]) createE2E(x T) *e2eEntry[T] {
	self.e2eIdxBase++
	return &e2eEntry[T]{orig: x, e2eIdx: self.e2eIdxBase}
}

func (self *PriorityQueueE2E[T]) Len() int {
	return self.impl.Len()
}

func (self *PriorityQueueE2E[T]) Empty() bool {
	return self.Len() == 0
}

func (self *PriorityQueueE2E[T]) Push(x T) {
	heap.Push(&self.impl, self.createE2E(x))
}

func (self *PriorityQueueE2E[T]) Pop() T {
	return heap.Pop(&self.impl).(*e2eEntry[T]).orig
}

func (self *PriorityQueueE2E[T]) Peek() T {
	return self.impl.Peek().orig
}

////////////////////////////////////////////////////////

// A priorityQueueImpl implements heap.Interface and holds Items.
type priorityQueueImpl[T any] struct {
	pq  []T
	cmp BiComparatorInterface[T]
}

// golang sort.Interface
func (self *priorityQueueImpl[T]) Len() int {
	return len(self.pq)
}

// golang sort.Interface
func (self *priorityQueueImpl[T]) Less(i, j int) bool {
	return self.cmp.Less(self.pq[i], self.pq[j])
}

// golang sort.Interface
func (self *priorityQueueImpl[T]) Swap(i, j int) {
	self.pq[i], self.pq[j] = self.pq[j], self.pq[i]
}

// golang heap.Interface
func (self *priorityQueueImpl[T]) Push(x any) {
	self.pq = append(self.pq, x.(T))
}

// golang heap.Interface
func (self *priorityQueueImpl[T]) Pop() any {
	n := len(self.pq)
	var result T
	result, self.pq[n-1] = self.pq[n-1], result // avoid memory leak
	self.pq = self.pq[0 : n-1]
	return result
}

func (self *priorityQueueImpl[T]) Peek() T {
	return self.pq[0]
}

////////////////////////////////////////////////////////
