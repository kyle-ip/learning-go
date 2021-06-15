package stack

type Stack struct {
    array []int
}

func (this *Stack) init(size int) {
    this.array = make([]int, size)
}

func (this *Stack) push(val int) (err error) {
    if cap(this.array) == len(this.array) {
        return err
    }
    this.array = append(this.array, val)
    return nil
}

func (this *Stack) pop() (val int, err error) {
    if len(this.array) == 0 {
        return -1, err
    }
    ret := this.array[len(this.array)-1]
    this.array = this.array[:len(this.array)-1]
    return ret, nil
}
