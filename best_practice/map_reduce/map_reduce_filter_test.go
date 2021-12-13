package map_reduce

import (
	"fmt"
	"strings"
	"testing"
)

// ========== Map ==========
// Map 函数在遍历第一个参数的数组，然后调用第二个参数的函数，把它的值组合成另一个数组返回。

func MapStrToStr(arr []string, fn func(s string) string) []string {
	var newArray []string
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func MapStrToInt(arr []string, fn func(s string) int) []int {
	var newArray []int
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func TestMap(t *testing.T) {

	var list = []string{"Kylo", "Yip", ""}

	x := MapStrToStr(list, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("%v\n", x)
	//["Kylo", "Yip", ""]

	y := MapStrToInt(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", y)
	//[3, 4, 8]
}

// ========== Reduce ==========

func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}
	return sum
}

func TestReduce(t *testing.T) {
	var list = []string{"Kylo", "Yip", ""}
	x := Reduce(list, func(s string) int {
		return len(s)
	})
	fmt.Printf("%v\n", x)
	// 15
}

// ========== Filter ==========

func Filter(arr []int, fn func(n int) bool) []int {
	var newArray []int
	for _, it := range arr {
		if fn(it) {
			newArray = append(newArray, it)
		}
	}
	return newArray
}

func TestFilter(t *testing.T) {
	var intSet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	out := Filter(intSet, func(n int) bool {
		return n%2 == 1
	})
	fmt.Printf("%v\n", out)

	out = Filter(intSet, func(n int) bool {
		return n > 5
	})
	fmt.Printf("%v\n", out)
}

// ========== 综合运用 ==========

type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

var list = []Employee{
	{"Kylo", 26, 0, 8000},
	{"Bob", 34, 10, 5000},
	{"Alice", 23, 5, 9000},
	{"Jack", 26, 0, 4000},
	{"Tom", 48, 9, 7500},
	{"Marry", 29, 0, 6000},
	{"Mike", 32, 8, 4000},
}

func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i := range list {
		if fn(&list[i]) {
			count += 1
		}
	}
	return count
}

func EmployeeFilterIn(list []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i := range list {
		if fn(&list[i]) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

func EmployeeSumIf(list []Employee, fn func(e *Employee) int) int {
	var sum = 0
	for i := range list {
		sum += fn(&list[i])
	}
	return sum
}

func TestMapReduceFilter(t *testing.T) {

	// 统计有多少员工大于 40 岁
	old := EmployeeCountIf(list, func(e *Employee) bool {
		return e.Age > 40
	})
	fmt.Printf("old people: %d\n", old)

	// 统计有多少员工的薪水大于 6000
	highPay := EmployeeCountIf(list, func(e *Employee) bool {
		return e.Salary > 6000
	})
	fmt.Printf("High Salary people: %d\n", highPay)

	// 列出有没有休假的员工
	noVacation := EmployeeFilterIn(list, func(e *Employee) bool {
		return e.Vacation == 0
	})
	fmt.Printf("People no vacation: %v\n", noVacation)

	// 统计所有员工的薪资总和
	totalPay := EmployeeSumIf(list, func(e *Employee) int {
		return e.Salary
	})
	fmt.Printf("Total Salary: %d\n", totalPay)

	// 统计 30 岁以下员工的薪资总和
	youngerPay := EmployeeSumIf(list, func(e *Employee) int {
		if e.Age < 30 {
			return e.Salary
		}
		return 0
	})
	fmt.Printf("Younger Salary: %d\n", youngerPay)
}

// https://github.com/robpike/filter
