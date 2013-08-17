// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package illegal

import (
	"fmt"
	"reflect"
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

	testFuncEqual(f1, 3, false, "illegal: passed non-function value to FuncEqual", t)
	testFuncEqual(3, f1, false, "illegal: passed non-function value to FuncEqual", t)

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

func testFuncEqual(f1, f2 interface{}, eq bool, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	eq1 := FuncEqual(f1, f2)

	if eq1 != eq {
		if eq {
			t.Errorf("Functions %v and %v are equal; FuncEqual said they weren't", f1, f2)
		} else {
			t.Errorf("Functions %v and %v are not equal; FuncEqual said they were", f1, f2)
		}
	}
}

func TestMap(t *testing.T) {
	// Map should succeed
	testMap([]int{1, 2, 3}, []int{1, 4, 9}, func(i int) int { return i * i }, nil, t)
	testMap([]int{}, []int{}, func(i int) int { return 3 }, nil, t)
	testMap([]int{1, 2, 3}, []bool{false, false, false}, func(i int) bool { return false }, nil, t)

	// Map should panic
	testMap([]int{}, nil, func(b bool) int { return 3 }, "illegal: function type and slice type do not match in call to Map(slice []T, fn func(T) S) []S", t)
	testMap([]int{}, nil, 3, "illegal: passed non-function value to Map", t)
	testMap(3, nil, func() {}, "illegal: passed non-slice value to Map", t)
	testMap([]int{}, nil, func(i, j int) int { return i * j }, "illegal: function type and slice type do not match in call to Map(slice []T, fn func(T) S) []S", t)
	testMap([]int{}, nil, func(i int) (int, int) { return i, i }, "illegal: function type and slice type do not match in call to Map(slice []T, fn func(T) S) []S", t)
}

func testMap(slc1, slc2, f interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	slc3 := Map(slc1, f)

	if !reflect.DeepEqual(slc2, slc3) {
		t.Errorf("Expected result %v; got %v", slc2, slc3)
	}
}

func TestFilter(t *testing.T) {
	// Filter should succeed
	testFilter([]int{1, 2, 3, 4}, []int{2, 4}, func(i int) bool { return i%2 == 0 }, nil, t)
	testFilter([]int{1, 2, 3, 4}, []int{}, func(i int) bool { return false }, nil, t)
	testFilter([]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, func(i int) bool { return true }, nil, t)

	// Filter should panic
	testFilter([]int{1, 2, 3, 4}, nil, func(b bool) bool { return false }, "illegal: function type and slice type do not match in call to Filter(slice []T, fn func(T) bool) []T", t)
	testFilter([]int{1, 2, 3, 4}, nil, func(i int) int { return i }, "illegal: function type and slice type do not match in call to Filter(slice []T, fn func(T) bool) []T", t)
	testFilter([]int{}, nil, 3, "illegal: passed non-function value to Filter", t)
	testFilter([]int{}, nil, func(i, j int) bool { return false }, "illegal: function type and slice type do not match in call to Filter(slice []T, fn func(T) bool) []T", t)
	testFilter([]int{}, nil, func(i, j int) (bool, bool) { return false, false }, "illegal: function type and slice type do not match in call to Filter(slice []T, fn func(T) bool) []T", t)
	testFilter(3, nil, func() {}, "illegal: passed non-slice value to Filter", t)
}

func testFilter(slc1, slc2, f interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	slc3 := Filter(slc1, f)

	if !reflect.DeepEqual(slc2, slc3) {
		t.Errorf("Expected result %v; got %v", slc2, slc3)
	}
}

func TestFoldr(t *testing.T) {
	// Foldr should succeed
	testFoldr([]int{1, 2, 3}, 1, func(i, j int) int { return i + j }, 7, nil, t)
	testFoldr([]int{1, 2, 3}, 6, func(i, j int) int { return j - i }, 0, nil, t)
	testFoldr([]int{1, 2, 3}, "", func(i int, s string) string { return fmt.Sprintf("%d%s", i, s) }, "321", nil, t)
	testFoldr([]int{}, "", func(i int, s string) string { return "foo" }, "", nil, t)

	// Foldr should fail
	testFoldr(3, 0, 0, nil, "illegal: passed non-slice value to Foldr", t)
	testFoldr([]int{}, 0, 0, nil, "illegal: passed non-function value to Foldr", t)
	testFoldr([]int{}, 0, func(i, j, k int) int { return 0 }, nil, "illegal: function type and slice type do not match in call to Foldr(slice []T, zero S, fn func(T, S) S) S", t)
	testFoldr([]int{}, 0, func(i, j int) (int, int) { return 0, 0 }, nil, "illegal: function type and slice type do not match in call to Foldr(slice []T, zero S, fn func(T, S) S) S", t)
	testFoldr([]int{}, false, func(i, j int) bool { return false }, nil, "illegal: function type and slice type do not match in call to Foldr(slice []T, zero S, fn func(T, S) S) S", t)
	testFoldr([]int{}, 0, func(i int, b bool) bool { return false }, nil, "illegal: zero type and function return type do not match in call to Foldr(slice []T, zero S, fn func(T, S) S) S", t)
	testFoldr([]int{}, 0, func(i, j int) bool { return false }, nil, "illegal: function type and slice type do not match in call to Foldr(slice []T, zero S, fn func(T, S) S) S", t)
}

func testFoldr(slc, z, f interface{}, res interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	res1 := Foldr(slc, z, f)

	if !reflect.DeepEqual(res, res1) {
		t.Errorf("Expected result %v; got %v", res, res1)
	}
}

func TestFoldl(t *testing.T) {
	// Foldr should succeed
	testFoldl([]int{1, 2, 3}, 1, func(i, j int) int { return i + j }, 7, nil, t)
	testFoldl([]int{3, 2, 1}, 6, func(i, j int) int { return i - j }, 0, nil, t)
	testFoldl([]int{1, 2, 3}, "", func(s string, i int) string { return fmt.Sprintf("%d%s", i, s) }, "123", nil, t)
	testFoldl([]int{}, "", func(s string, i int) string { return "foo" }, "", nil, t)

	// Foldr should fail
	testFoldl(3, 0, 0, nil, "illegal: passed non-slice value to Foldl", t)
	testFoldl([]int{}, 0, 0, nil, "illegal: passed non-function value to Foldl", t)
	testFoldl([]int{}, 0, func(i, j, k int) int { return 0 }, nil, "illegal: function type and slice type do not match in call to Foldl(slice []T, zero S, fn func(S, T) S) S", t)
	testFoldl([]int{}, 0, func(i, j int) (int, int) { return 0, 0 }, nil, "illegal: function type and slice type do not match in call to Foldl(slice []T, zero S, fn func(S, T) S) S", t)
	testFoldl([]int{}, false, func(i, j int) bool { return false }, nil, "illegal: function type and slice type do not match in call to Foldl(slice []T, zero S, fn func(S, T) S) S", t)
	testFoldl([]int{}, 0, func(b bool, i int) bool { return false }, nil, "illegal: zero type and function return type do not match in call to Foldl(slice []T, zero S, fn func(S, T) S) S", t)
	testFoldl([]int{}, 0, func(i, j int) bool { return false }, nil, "illegal: function type and slice type do not match in call to Foldl(slice []T, zero S, fn func(S, T) S) S", t)
}

func testFoldl(slc, z, f interface{}, res interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	res1 := Foldl(slc, z, f)

	if !reflect.DeepEqual(res, res1) {
		t.Errorf("Expected result %v; got %v", res, res1)
	}
}
