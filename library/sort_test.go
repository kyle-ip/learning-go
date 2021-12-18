package library

import (
	"fmt"
	"sort"
	"testing"
)

// 学生结构体
type Student struct {
	name  string // 姓名
	score int    // 成绩
}

type StuScores []Student

//Len()
func (s StuScores) Len() int {
	return len(s)
}

//Less(): 成绩将有低到高排序
func (s StuScores) Less(i, j int) bool {
	return s[i].score < s[j].score
}

// Swap()
func (s StuScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 定义了一个 reverse 结构类型，嵌入 Interface 接口
type reverse struct {
	sort.Interface
}

// reverse 结构类型的 Less() 方法拥有嵌入的 Less() 方法相反的行为
// Len() 和 Swap() 方法则会保持嵌入类型的方法行为
func (r reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}

func TestSorting(t *testing.T) {
	s := []int{5, 2, 6, 3, 1, 4}

	sort.Ints(s)
	fmt.Println(s)

	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	fmt.Println(s)
}

func TestStructSorting(t *testing.T) {
	students := StuScores{
		{"Kylo", 95},
		{"Janice", 91},
		{"Ben", 96},
		{"Jack", 90},
	}

	// 打印未排序的 students 数据
	fmt.Println("Default:\n\t", students)

	//StuScores 已经实现了 sort.Interface 接口 , 所以可以调用 Sort 函数进行排序
	sort.Sort(students)

	// 判断是否已经排好顺序，将会打印 true
	fmt.Println("IS Sorted?\n\t", sort.IsSorted(students))

	// 打印排序后的 students 数据
	fmt.Println("Sorted:\n\t", students)

	sort.Sort(sort.Reverse(students))
	fmt.Println("Reverse:\n\t", students)
}

func TestInterfaceSorting(t *testing.T) {
	people := []struct {
		Name string
		Age  int
	}{
		{"Gopher", 7},
		{"Alice", 55},
		{"Vera", 24},
		{"Bob", 75},
	}

	sort.SliceStable(people, func(i, j int) bool { return people[i].Age > people[j].Age }) // 按年龄降序排序
	fmt.Println("Sort by age:", people)

	sort.Slice(people, func(i, j int) bool { return people[i].Age < people[j].Age }) // 按年龄升序排序
	fmt.Println("Sort by age:", people)

	fmt.Println("Sorted:", sort.SliceIsSorted(people, func(i, j int) bool { return people[i].Age < people[j].Age }))
}

func TestBinarySearch(t *testing.T) {
	a := []int{2, 3, 4, 200, 100, 21, 234, 56}
	x := 21

	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })           // 升序排序
	index := sort.Search(len(a), func(i int) bool { return a[i] >= x }) // 查找元素

	if index < len(a) && a[index] == x {
		fmt.Printf("found %d at index %d in %v\n", x, index, a)
	} else {
		fmt.Printf("%d not found in %v,index:%d\n", x, a, index)
	}
}

func TestGuessingGame(t *testing.T) {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		_, err := fmt.Scanf("%s", &s)
		if err != nil {
			return false
		}
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
