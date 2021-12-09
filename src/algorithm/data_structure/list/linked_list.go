package list

type ListNode struct {
	Key, Val   int
	Next, Prev *ListNode
}

func (this *ListNode) Add(l *ListNode) {
	this.Next = l
	l.Prev = this
}

func (this *ListNode) find(val int) *ListNode {
	for p := this; p != nil; p = p.Next {
		if p.Val == val {
			return p
		}
	}
	return nil
}

func (this *ListNode) remove(val int) {
	dummy := ListNode{Next: this}
	for p := &dummy; p.Next != nil; {
		if p.Next.Val == val {
			p.Next = p.Next.Next
			if p.Next.Next != nil {
				p.Next.Next.Prev = p
			}
		}
	}
}
