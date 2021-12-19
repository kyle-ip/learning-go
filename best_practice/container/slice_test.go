package container

import (
	"bytes"
	"fmt"
	"testing"
)

/**
type slice struct {
    array unsafe.Pointer // 指向底层数组的指针
    len   int            // 长度
    cap   int            // 容量
}
*/

type data struct {
}

// ========== Slice 共享内存 ==========

func TestSlice(t *testing.T) {
	// 创建一个长度和容量都是 5 的 Slice，对下标 3、4 的位置赋值。
	a := make([]int, 5)
	a[3] = 42
	a[4] = 100
	fmt.Println(a)
	// [0 0 0 42 100]

	// 切分后赋值给 b，然后赋值 b[1]，由于此时 a、b 仍共享内存，所以会影响到 a 中的值。
	b := a[1:4]
	b[1] = 99
	// a: [0 0 0 42 100]   =>   [0 0 99 42 100]
	// b:   [0 0 42]              [0 99 42]
	fmt.Println(a)

	// 此时对 a 进行 append 操作，会对 a 重新分配内存，导致 a 和 b 不再共享。
	// append 只有在 cap 不够用时，才会重新分配内存以扩大容量。
	a = append(a, 11)
	b[0] = 52
	fmt.Println(a)
	fmt.Println(b)
	// a: [0 0 0 42 100 11]
	// b: [52 99 42]

}

func TestSliceSharedMemory(t *testing.T) {
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	// dir1 和 dir2 共享内存。
	// AAAA
	dir1 := path[:sepIndex]
	// BBBBBBBBB
	dir2 := path[sepIndex+1:]
	dir3 := path[:sepIndex:sepIndex]

	// Full Slice Expression：要使 dir1、dir2 独立占有内存，只需要把 dir1 的定义改为：
	//      dir1 := path[:sepIndex:sepIndex]
	// 最后一个参数 Limited Capacity，可使得后续的 append() 操作重新分配内存。

	dir3 = append(dir3, "suffix"...)
	fmt.Println("dir1 =>", string(dir1)) // prints: AAAA
	fmt.Println("dir2 =>", string(dir2)) // prints: BBBBBBBBB
	fmt.Println("dir3 =>", string(dir3)) // prints: AAAAsuffix

	// 因为 cap 足够，数据扩展到了 dir2 的空间。
	dir1 = append(dir1, "suffix"...)
	fmt.Println("dir1 =>", string(dir1)) // prints: AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) // prints: uffixBBBB
	fmt.Println("dir3 =>", string(dir3)) // prints: AAAAsuffix
}
