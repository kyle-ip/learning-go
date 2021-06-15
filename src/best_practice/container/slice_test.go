package container

import (
    "bytes"
    "fmt"
    "reflect"
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

    // 切分后赋值给 b，然后赋值 b[1]。
    b := a[1:4]
    b[1] = 99

    // 假使对 a 进行 append 操作，会对 a 重新分配内存，导致 a 和 b 不再共享。
    // append 只有在 cap 不够用时，才会重新分配内存以扩大容量。
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

// ========== 深度比较 ==========

func TestObjectDeepEqual(t *testing.T) {
    // 在复制一个对象时，它可以是内建数据类型、数组、结构体、Map。
    // 比如在复制结构体时，如果需要比较两个结构体中的数据是否相同，就要使用深度比较（使用反射 reflect.DeepEqual()）。
    v1 := data{}
    v2 := data{}
    fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))
    //prints: v1 == v2: true

    m1 := map[string]string{"one": "a", "two": "b"}
    m2 := map[string]string{"two": "b", "one": "a"}
    fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))
    // prints: m1 == m2: true

    s1 := []int{1, 2, 3}
    s2 := []int{1, 2, 3}
    fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))
    // prints: s1 == s2: true
}
