// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package illegal

import (
	"fmt"
	"reflect"
	"runtime"
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
	/*
		f1 := func() {}
		testFuncEqual(f1, f1, fname(f1), fname(f1), true, nil, t)

		f2 := func() bool { return false }
		testFuncEqual(f2, f2, fname(f2), fname(f2), true, nil, t)

		f3 := func(b bool) (bool, error) { return b, nil }
		testFuncEqual(f3, f3, fname(f3), fname(f3), true, nil, t)

		testFuncEqual(f1, f2, fname(f1), fname(f2), false, nil, t)
		testFuncEqual(f1, f3, fname(f1), fname(f3), false, nil, t)
		testFuncEqual(f2, f3, fname(f2), fname(f3), false, nil, t)

		testFuncEqual(f1, 3, fname(f1), "", false, "illegal.FuncEqual: passed non-function value", t)
		testFuncEqual(3, f1, "", fname(f1), false, "illegal.FuncEqual: passed non-function value", t)

		f4 := func(i int) func() int { return func() int { return i } }

		testFuncEqual(f4, f4, fname(f4), fname(f4), true, nil, t)
	*/
	f1 := func() {}
	testFuncEqual(f1, f1, true, nil, t)

	f2 := func() bool { return false }
	testFuncEqual(f2, f2, true, nil, t)

	f3 := func(b bool) (bool, error) { return b, nil }
	testFuncEqual(f3, f3, true, nil, t)

	testFuncEqual(f1, f2, false, nil, t)
	testFuncEqual(f1, f3, false, nil, t)
	testFuncEqual(f2, f3, false, nil, t)

	testFuncEqual(f1, 3, false, "illegal.FuncEqual: passed non-function value", t)
	testFuncEqual(3, f1, false, "illegal.FuncEqual: passed non-function value", t)

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

	testFuncEqual(t1.Test1, t2.Test2, false, nil, t)

	// testFuncEqual()
}

func testFuncEqual(f1, f2 interface{}, eq bool, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	eq1 := FuncEqual(f1, f2)

	if eq1 != eq {
		n1, n2 := "", ""
		v1, v2 := reflect.ValueOf(f1), reflect.ValueOf(f2)
		if v1.Kind() == reflect.Func {
			n1 = fname(v1.Pointer())
		}
		if v2.Kind() == reflect.Func {
			n2 = fname(v2.Pointer())
		}
		if eq {
			t.Errorf("Functions %v (%v) and %v (%v) are equal; FuncEqual said they weren't", f1, n1, f2, n2)
		} else {
			t.Errorf("Functions %v (%v) and %v (%v) are not equal; FuncEqual said they were", f1, n1, f2, n2)
		}
	}
}

type IntAlias int
type IntAlias2 int
type EmptyStructAlias struct{}

var InterfaceReflectType = reflect.TypeOf([]interface{}{}).Elem()

type TypeWithMethod int

func (t TypeWithMethod) Int() int { return int(t) }

// Tests both ConvertSlice and ConvertSliceType
func TestConvertSlice(t *testing.T) {
	testConvertSlice([]int{1, 2, 3}, []int{1, 2, 3}, int(0), nil, t)
	testConvertSlice([]int{1, 2, 3}, []int64{1, 2, 3}, int64(0), nil, t)
	testConvertSlice([]int{1, 2, 3}, []float64{1.0, 2.0, 3.0}, float64(0), nil, t)
	testConvertSlice([]int{1, 2, 3}, []IntAlias{1, 2, 3}, IntAlias(0), nil, t)
	testConvertSlice([]struct{}{struct{}{}}, []EmptyStructAlias{EmptyStructAlias{}}, EmptyStructAlias{}, nil, t)
	testConvertSlice([]IntAlias{1, 2, 3}, []IntAlias2{1, 2, 3}, IntAlias2(0), nil, t)
	testConvertSlice([]int{1, 2, 3}, []interface{}{1, 2, 3}, InterfaceReflectType, nil, t)

	method1 := TypeWithMethod.Int
	// method2 := (TypeWithMethod(3)).Int

	result1 := ConvertSlice([](func(TypeWithMethod) int){method1}, method1).([](func(TypeWithMethod) int))[0]
	fmt.Println(reflect.TypeOf(result1))
	fmt.Println(result1)
	fmt.Println(FuncEqual(method1, result1))
	// result2 := ConvertSlice([](func() int){method2}, method2).([]func() int)[0]
	// fmt.Println(reflect.TypeOf(result2))
	// fmt.Println(result2)
	// fmt.Println(FuncEqual(method2, result2))
}

// Tests both ConvertSlice and ConvertSliceType:
// example can either be an example value,
// in which case ConvertSlice will be called,
// or a reflect.Type value, in which case
// ConvertSliceType will be called.
func testConvertSlice(input, target, example interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	var result interface{}
	if typ, ok := example.(reflect.Type); ok {
		result = ConvertSliceType(input, typ)
	} else {
		result = ConvertSlice(input, example)
	}
	if !reflect.DeepEqual(target, result) {
		t.Errorf("Expected %s(%v); got %s(%v)", reflect.TypeOf(target).String(), target, reflect.TypeOf(result).String(), result)
	}
}

func TestConvenienceVariables(t *testing.T) {
	expect := "interface {}"
	get := fmt.Sprint(InterfaceType)
	if get != expect {
		t.Errorf("Expected %s; got %s", expect, get)
	}
}

func fname(p uintptr) string {
	f := runtime.FuncForPC(p)
	if f == nil {
		return ""
	}
	return f.Name()

}
