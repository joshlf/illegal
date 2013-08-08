package illegal

import (
	// "fmt"
	"reflect"
	// "runtime"
	"testing"
)

type FuncEqualTestInterface interface {
	Test1()
	Test2()
}

type FuncEqualTestType1 struct{}

func (f FuncEqualTestType1) Test1() {}

func (f FuncEqualTestType1) Test2() {}

type FuncEqualTestType2 struct{}

func (f FuncEqualTestType2) Test1() {}

func (f FuncEqualTestType2) Test2() {}

func TestFuncEqual(t *testing.T) {
	f1 := func() {}
	testFuncEqual(f1, f1, true, nil, t)

	f2 := func() bool { return false }
	testFuncEqual(f2, f2, true, nil, t)

	f3 := func(b bool) (bool, error) { return b, nil }
	testFuncEqual(f3, f3, true, nil, t)

	testFuncEqual(f1, f2, false, nil, t)
	testFuncEqual(f1, f3, false, nil, t)
	testFuncEqual(f2, f3, false, nil, t)

	testFuncEqual(f1, 3, false, ErrNotFunc, t)
	testFuncEqual(3, f1, false, ErrNotFunc, t)

	f4 := func(i int) func() int { return func() int { return i } }

	testFuncEqual(f4, f4, true, nil, t)

	f5 := f4(1)
	f6 := f4(2)

	testFuncEqual(f5, f6, true, nil, t)

	f7 := FuncEqualTestType1.Test1
	f8 := FuncEqualTestType1.Test2

	testFuncEqual(f7, f7, true, nil, t)
	testFuncEqual(f7, f8, false, nil, t)

	t1 := FuncEqualTestType1{}
	t2 := FuncEqualTestType1{}

	testFuncEqual(t1.Test1, t2.Test1, true, nil, t)

	t3 := FuncEqualTestType2{}
	t4 := FuncEqualTestType2{}

	testFuncEqual(t3.Test1, t4.Test1, true, nil, t)
	testFuncEqual(t1.Test1, t3.Test1, false, nil, t)

	var i1 FuncEqualTestInterface = t1
	var i2 FuncEqualTestInterface = t2

	testFuncEqual(i1.Test1, i2.Test1, true, nil, t)

	i2 = t3

	testFuncEqual(i1.Test1, i2.Test1, true, nil, t)
}

func testFuncEqual(f1, f2 interface{}, eq bool, err error, t *testing.T) {
	eq1, err1 := FuncEqual(f1, f2)

	if eq1 != eq {
		if eq {
			t.Errorf("Functions %v and %v are equal; FuncEqual said they weren't", f1, f2)
		} else {
			t.Errorf("Functions %v and %v are not equal; FuncEqual said they were", f1, f2)
		}
	}

	if err1 != err {
		t.Errorf("Expected error %v; got %v", err, err1)
	}
}

func TestMap(t *testing.T) {
	// Map should succeed
	testMap([]int{1, 2, 3}, []int{1, 4, 9}, func(i int) int { return i * i }, nil, t)
	testMap([]int{}, []int{}, func(i int) int { return 3 }, nil, t)
	testMap([]int{1, 2, 3}, []bool{false, false, false}, func(i int) bool { return false }, nil, t)

	// Map should return an error
	testMap([]int{}, nil, func(b bool) int { return 3 }, ErrWrongFuncType, t)
	testMap([]int{}, nil, 3, ErrNotFunc, t)
	testMap(3, nil, func() {}, ErrNotSlice, t)
	testMap([]int{}, nil, func(i, j int) int { return i * j }, ErrWrongFuncType, t)
	testMap([]int{}, nil, func(i int) (int, int) { return i, i }, ErrWrongFuncType, t)
}

func testMap(slc1, slc2, f interface{}, err error, t *testing.T) {
	slc3, err1 := Map(slc1, f)

	if !reflect.DeepEqual(slc2, slc3) {
		t.Errorf("Expected result %v; got %v", slc2, slc3)
	}

	if err1 != err {
		t.Errorf("Expected error %v; got %v", err, err1)
	}
}

func TestFilter(t *testing.T) {
	// Filter should succeed
	testFilter([]int{1, 2, 3, 4}, []int{2, 4}, func(i int) bool { return i%2 == 0 }, nil, t)
	testFilter([]int{1, 2, 3, 4}, []int{}, func(i int) bool { return false }, nil, t)
	testFilter([]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, func(i int) bool { return true }, nil, t)

	// Filter should return an error
	testFilter([]int{1, 2, 3, 4}, nil, func(b bool) bool { return false }, ErrWrongFuncType, t)
	testFilter([]int{1, 2, 3, 4}, nil, func(i int) int { return i }, ErrWrongFuncType, t)
	testFilter([]int{}, nil, 3, ErrNotFunc, t)
	testFilter([]int{}, nil, func(i, j int) bool { return false }, ErrWrongFuncType, t)
	testFilter([]int{}, nil, func(i, j int) (bool, bool) { return false, false }, ErrWrongFuncType, t)
}

func testFilter(slc1, slc2, f interface{}, err error, t *testing.T) {
	slc3, err1 := Filter(slc1, f)

	if !reflect.DeepEqual(slc2, slc3) {
		t.Errorf("Expected result %v; got %v", slc2, slc3)
	}

	if err1 != err {
		t.Errorf("Expected error %v; got %v", err, err1)
	}
}
