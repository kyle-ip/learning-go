package std_lib

import (
    "container/heap"
    "container/list"
    "container/ring"
    "fmt"
    "testing"
)

// ========== Heap ==========

type IntHeap []int

func (h IntHeap) Len() int { return len(h) }

func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }

func (h IntHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func TestHeap(t *testing.T) {
    h := &IntHeap{2, 1, 5}
    heap.Init(h)
    heap.Push(h, 3)
    heap.Pop(h)
    //                   0
    //            1            2
    //         3    4       5      6
    //        7 8  9 10   11
}

// ========== Linked List ==========

func TestLinkedList(t *testing.T) {
    l := list.New()
    l.PushBack(1)
    l.PushBack(2)

    fmt.Printf("len: %v\n", l.Len())
    fmt.Printf("first: %#v\n", l.Front())
    fmt.Printf("second: %#v\n", l.Front().Next())
}

// ========== Ring ==========

func TestRing(t *testing.T) {
    ring := ring.New(3)

    for i := 1; i <= 3; i++ {
        ring.Value = i
        ring = ring.Next()
    }

    // 计算 1+2+3
    s := 0
    ring.Do(func(p interface{}) {
        s += p.(int)
    })
    fmt.Println("sum is", s)
}
