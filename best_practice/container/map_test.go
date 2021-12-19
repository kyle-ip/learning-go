package container

import (
	"fmt"
	"reflect"
	"testing"
)

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
