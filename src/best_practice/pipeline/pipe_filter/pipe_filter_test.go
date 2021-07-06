package pipe_filter

import (
    "reflect"
    "testing"
)

func TestStringSplit(t *testing.T) {
    sf := NewSplitFilter(",")
    resp, err := sf.Process("1,2,3")
    if err != nil {
        t.Fatal(err)
    }
    parts, ok := resp.([]string)
    if !ok {
        t.Fatalf("Repsonse type is %T, but the expected type is string", parts)
    }
    if !reflect.DeepEqual(parts, []string{"1", "2", "3"}) {
        t.Errorf("Expected value is {\"1\",\"2\",\"3\"}, but actual is %v", parts)
    }
}

func TestWrongInput(t *testing.T) {
    sf := NewSplitFilter(",")
    _, err := sf.Process(123)
    if err == nil {
        t.Fatal("An error is expected.")
    }
}

func TestSumElems(t *testing.T) {
    sf := NewSumFilter()
    ret, err := sf.Process([]int{1, 2, 3})
    if err != nil {
        t.Fatal(err)
    }
    expected := 6
    if ret != expected {
        t.Fatalf("The expected is %d, but actual is %d", expected, ret)
    }
}

func TestWrongInputForSumFilter(t *testing.T) {
    sf := NewSumFilter()
    _, err := sf.Process([]float32{1.1, 2.2, 3.1})

    if err == nil {
        t.Fatal("An error is expected.")
    }
}

func TestConvertToInt(t *testing.T) {
    tif := NewToIntFilter()
    resp, err := tif.Process([]string{"1", "2", "3"})
    if err != nil {
        t.Fatal(err)
    }
    expected := []int{1, 2, 3}
    if !reflect.DeepEqual(expected, resp) {
        t.Fatalf("The expected is %v, the actual is %v", expected, resp)
    }
}

func TestWrongInputForTIF(t *testing.T) {
    tif := NewToIntFilter()
    _, err := tif.Process([]string{"1", "2.2", "3"})
    if err == nil {
        t.Fatal("An error is expected for wrong input")
    }
    _, err = tif.Process([]int{1, 2, 3})
    if err == nil {
        t.Fatal("An error is expected for wrong input")
    }
}

func TestStraightPipeline(t *testing.T) {
    spliter := NewSplitFilter(",")
    converter := NewToIntFilter()
    sum := NewSumFilter()
    sp := NewStraightPipeline("p1", spliter, converter, sum)
    ret, err := sp.Process("1,2,3")
    if err != nil {
        t.Fatal(err)
    }
    if ret != 6 {
        t.Fatalf("The expected is 6, but the actual is %d", ret)
    }
}
