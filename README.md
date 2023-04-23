# Priority queue and on time delivery queue.

For Go 1.18+ (with generics)
Based on [heap](https://pkg.go.dev/container/heap) and [delayqueue](https://pkg.go.dev/github.com/golearnku/delayqueue).

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
