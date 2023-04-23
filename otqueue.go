package goontime

import (
	"context"
	"sync"
	"time"
)

////////////////////////////////////////////////////////

// A OnTimeQueue
type OnTimeQueue[T any] struct {
	guard   sync.Mutex
	active  bool
	pq      *PriorityQueueE2E[*onTimeElement[T]]
	updChan *oneSignalChannel
}

func NewOnTimeQueue[T any](cap int) *OnTimeQueue[T] {
	cmp := &onTimeBiComparator[T]{}
	result := &OnTimeQueue[T]{
		active:  true,
		pq:      NewPriorityQueueE2ECap[*onTimeElement[T]](cmp, cap),
		updChan: newOneSignalChannel(),
	}
	return result
}

func (self *OnTimeQueue[T]) Shutdown() {
	self.guard.Lock()
	defer self.guard.Unlock()

	if self.active {
		self.active = false
		self.updChan.Close()
	}
}

func (self *OnTimeQueue[T]) Add(data T, timeout time.Duration) {

	self.guard.Lock()
	defer self.guard.Unlock()

	if self.active {
		elem := &onTimeElement[T]{
			time: time.Now().Add(timeout),
			data: data,
		}
		var updated bool = true
		if !self.pq.Empty() {
			updated = elem.time.Before(self.pq.Peek().time)
		}
		self.pq.Push(elem)
		if updated {
			self.updChan.TrySignal()
		}
	}
}

var timeoutDefault time.Duration = time.Duration(10 * time.Millisecond)

func (self *OnTimeQueue[T]) Get() (data T, timeout time.Duration, ok bool) {
	self.guard.Lock()
	defer self.guard.Unlock()

	if self.active {
		if !self.pq.Empty() {
			dur := time.Until(self.pq.Peek().time)
			if dur > 0 {
				timeout = dur
				return
			} else {
				data = self.pq.Pop().data
				timeout = time.Duration(0)
				ok = true
				return
			}
		}
	}

	timeout = timeoutDefault
	return
}

////////////////////////////////////////////////////////

func CreateOnTimeQueueChannel[T any](ctx context.Context, wg *sync.WaitGroup, tq *OnTimeQueue[T]) chan T {
	inputCh := make(chan T)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inputCh)

		timer := time.NewTimer(time.Duration(0))

		for {

			data, dur, ok := tq.Get()
			if ok {
				inputCh <- data
				continue
			}

			timer.Reset(dur)

			select {
			case <-ctx.Done():
				return
			case sig := <-tq.updChan.ch:
				if !sig {
					return
				}
			case <-timer.C:
			}
		}
	}()

	return inputCh
}

////////////////////////////////////////////////////////

type onTimeElement[T any] struct {
	time time.Time
	data T
}

type onTimeBiComparator[T any] struct{}

func (*onTimeBiComparator[T]) Less(lhs, rhs *onTimeElement[T]) bool {
	return lhs.time.Before(rhs.time)
}

/////////////////////////////////////////////////////////////

type oneSignalChannel struct {
	ch chan bool
}

func newOneSignalChannel() *oneSignalChannel {
	return &oneSignalChannel{
		ch: make(chan bool, 1),
	}
}

func (self *oneSignalChannel) TrySignal() {
	select {
	case self.ch <- true:
	default: // channel is full
	}
}

func (self *oneSignalChannel) Close() {
	close(self.ch)
}

/////////////////////////////////////////////////////////////
