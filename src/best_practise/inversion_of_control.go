package best_practise

import "errors"

// 控制反转：一种软件设计的方法，把控制逻辑与业务逻辑分开，不在业务逻辑里写控制逻辑（会让控制逻辑依赖于业务逻辑），而是反过来让业务逻辑依赖控制逻辑。
// 这样的编程方式可以有效降低程序复杂度，并提升代码重用度。

// ========== 控制反转 ==========

type IntSet struct {
	data map[int]bool
}

func NewIntSet() IntSet {
	return IntSet{make(map[int]bool)}
}

func (set *IntSet) Add(x int) {
	set.data[x] = true
}

func (set *IntSet) Delete(x int) {
	delete(set.data, x)
}

func (set *IntSet) Contains(x int) bool {
	return set.data[x]
}

// IntSet 的扩展：加入 Undo 功能。

type UndoableIntSet struct { // Poor style
	IntSet    // Embedding (delegation)
	functions []func()
}

func NewUndoableIntSet() UndoableIntSet {
	return UndoableIntSet{NewIntSet(), nil}
}

// 在 UndoableIntSet 中嵌入了IntSet ，Override 了 它的 Append() 和 Delete() 方法；
// Contains() 方法没有 Override，直接被带到 UndoableInSet 中。

func (set *UndoableIntSet) Add(x int) { // Override
	// 在 Append() 中，记录 Delete 操作。
	if !set.Contains(x) {
		set.data[x] = true
		set.functions = append(set.functions, func() { set.Delete(x) })
	} else {
		set.functions = append(set.functions, nil)
	}
}

func (set *UndoableIntSet) Delete(x int) { // Override
	// 在 Delete() 中，记录 Append 操作。
	if set.Contains(x) {
		delete(set.data, x)
		set.functions = append(set.functions, func() { set.Add(x) })
	} else {
		set.functions = append(set.functions, nil)
	}
}

func (set *UndoableIntSet) Undo() error {
	if len(set.functions) == 0 {
		return errors.New("No functions to undo")
	}
	// 读取此前记录的 Append 和 Delete 操作，实现 Undo。
	index := len(set.functions) - 1
	if function := set.functions[index]; function != nil {
		function()
		set.functions[index] = nil // For garbage collection
	}
	set.functions = set.functions[:index]
	return nil
}

// ========== 反转依赖 ==========

// 声明函数接口，表示 Undo 控制可以接受的函数签名。

type Undo []func()

func (undo *Undo) Append(function func()) {
	*undo = append(*undo, function)
}

func (undo *Undo) Undo() error {
	functions := *undo
	if len(functions) == 0 {
		return errors.New("No functions to undo")
	}
	// 取 function 数组的最后一个，不为空则执行，执行完成就置空。
	index := len(functions) - 1
	if function := functions[index]; function != nil {
		function()
		functions[index] = nil // For gc
	}
	// 截去刚才被执行的那个 function。
	*undo = functions[:index]
	return nil
}

// 反转依赖：由业务逻辑 IntSet 依赖 Undo（即一个协议，表示为没有参数的函数数组，目的是服用 Undo 的代码）。

type UndoIntSet struct {
	data map[int]bool
	undo Undo
}

func NewUndoIntSet() UndoIntSet {
	return UndoIntSet{data: make(map[int]bool)}
}

func (set *UndoIntSet) Undo() error {
	return set.undo.Undo()
}

func (set *UndoIntSet) Contains(x int) bool {
	return set.data[x]
}

func (set *UndoIntSet) Add(x int) {
	if !set.Contains(x) {
		set.data[x] = true
		// 在 undo 中加入一个 Delete(x) 的函数，用于执行 Undo 时弹出、复原（Add 同理）。
		set.undo.Append(func() { set.Delete(x) })
	} else {
		set.undo.Append(nil)
	}
}

func (set *UndoIntSet) Delete(x int) {
	if set.Contains(x) {
		delete(set.data, x)
		set.undo.Append(func() { set.Add(x) })
	} else {
		set.undo.Append(nil)
	}
}
