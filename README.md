# Priority queue and on time delivery queue.

For Go 1.18+ (with generics). Based on [heap](https://pkg.go.dev/container/heap) and [delayqueue](https://pkg.go.dev/github.com/golearnku/delayqueue).

---

## Priority queue

Go heap does not preserve order of insertion of same priority elements.
PriorityQueue behaves as exactly as heap does.
PriorityQueueE2E preserves order of insertion. By adding end-to-end index based on uint64.
Both PriorityQueue and PriorityQueueE2E are not concurrent programming safe.

### Example


```go

import (
	"fmt"

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

func main() {

	pq := got.NewPriorityQueueCap[*Elem](cmpAscending, 5)

	pq.Push(NewElem(15, "15"))
	pq.Push(NewElem(8, "8"))

	for !pq.Empty() {
		fmt.Printf("%v", pq.Pop())
	}
}

```

---

## On Time delivery queue

OnTimeQueue. Based on PriorityQueueE2E. Ready for concurrent programming.

### Example

```go

import (
	"context"
	"fmt"
	"sync"
	"time"

	got "github.com/xfuflo101/goontime"
)

func main() {

	var wg sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())

	onTimeQueue := got.NewOnTimeQueue[string](10)
	defer onTimeQueue.Shutdown()

	onTimeChan := got.CreateOnTimeQueueChannel(ctx, &wg, onTimeQueue)
	// this channel will close writer by ctx cancel

	start := time.Now()

	var resultStr string

	wg.Add(1)
	go func() {
		// reader
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case str := <-onTimeChan:
				resultStr += fmt.Sprintf("; [%v : %v]", str, time.Now().Sub(start).Milliseconds()/100)
			}
		}
	}()

	onTimeQueue.Add("10", time.Duration(10*time.Second))
	onTimeQueue.Add("8", time.Duration(8*time.Second))
	onTimeQueue.Add("3", time.Duration(3*time.Second))
	onTimeQueue.Add("2", time.Duration(2*time.Second))
	onTimeQueue.Add("8_2", time.Duration(8*time.Second))
	onTimeQueue.Add("1", time.Duration(1*time.Second))

	time.Sleep(11 * time.Second)

	ctxCancel()

	wg.Wait()

	fmt.Printf("Result: [%v]\n", resultStr)
}

```