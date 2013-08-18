// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIdentity(t *testing.T) {
	testIdentity(nil, t)
	testIdentity(1, t)
	testIdentity(3.4, t)
	testIdentity(make([]int, 3, 6), t)
	testIdentity(new(int), t)
	testIdentity(make(map[int](map[int]string)), t)
}

func testIdentity(x interface{}, t *testing.T) {
	y := Identity(x)
	if !reflect.DeepEqual(x, y) {
		t.Errorf("Expected %v; got %v", x, y)
	}
}

func TestMap(t *testing.T) {
	// Map should succeed
	testMap([]int{1, 2, 3}, []int{1, 4, 9}, func(i int) int { return i * i }, nil, t)
	testMap([]int{}, []int{}, func(i int) int { return 3 }, nil, t)
	testMap([]int{1, 2, 3}, []bool{false, false, false}, func(i int) bool { return false }, nil, t)

	// Map should panic
	testMap([]int{}, nil, func(b bool) int { return 3 }, mapTypeError, t)
	testMap([]int{}, nil, 3, mapFunctionError, t)
	testMap(3, nil, func() {}, mapSliceError, t)
	testMap([]int{}, nil, func(i, j int) int { return i * j }, mapTypeError, t)
	testMap([]int{}, nil, func(i int) (int, int) { return i, i }, mapTypeError, t)
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
	testFilter([]int{1, 2, 3, 4}, nil, func(b bool) bool { return false }, filterTypeError, t)
	testFilter([]int{1, 2, 3, 4}, nil, func(i int) int { return i }, filterTypeError, t)
	testFilter([]int{}, nil, 3, filterFunctionError, t)
	testFilter([]int{}, nil, func(i, j int) bool { return false }, filterTypeError, t)
	testFilter([]int{}, nil, func(i, j int) (bool, bool) { return false, false }, filterTypeError, t)
	testFilter(3, nil, func() {}, filterSliceError, t)
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

func TestReject(t *testing.T) {
	// Reject should succeed
	testReject([]int{1, 2, 3, 4}, []int{1, 3}, func(i int) bool { return i%2 == 0 }, nil, t)
	testReject([]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, func(i int) bool { return false }, nil, t)
	testReject([]int{1, 2, 3, 4}, []int{}, func(i int) bool { return true }, nil, t)

	// Reject should panic
	testReject([]int{1, 2, 3, 4}, nil, func(b bool) bool { return false }, rejectTypeError, t)
	testReject([]int{1, 2, 3, 4}, nil, func(i int) int { return i }, rejectTypeError, t)
	testReject([]int{}, nil, 3, rejectFunctionError, t)
	testReject([]int{}, nil, func(i, j int) bool { return false }, rejectTypeError, t)
	testReject([]int{}, nil, func(i, j int) (bool, bool) { return false, false }, rejectTypeError, t)
	testReject(3, nil, func() {}, rejectSliceError, t)
}

func testReject(slc1, slc2, f interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	slc3 := Reject(slc1, f)

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
	testFoldr(3, 0, 0, nil, foldrSliceError, t)
	testFoldr([]int{}, 0, 0, nil, foldrFunctionError, t)
	testFoldr([]int{}, 0, func(i, j, k int) int { return 0 }, nil, foldrTypeError, t)
	testFoldr([]int{}, 0, func(i, j int) (int, int) { return 0, 0 }, nil, foldrTypeError, t)
	testFoldr([]int{}, false, func(i, j int) bool { return false }, nil, foldrTypeError, t)
	testFoldr([]int{}, 0, func(i int, b bool) bool { return false }, nil, foldrZeroError, t)
	testFoldr([]int{}, 0, func(i, j int) bool { return false }, nil, foldrTypeError, t)
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
	// Foldl should succeed
	testFoldl([]int{1, 2, 3}, 1, func(i, j int) int { return i + j }, 7, nil, t)
	testFoldl([]int{3, 2, 1}, 6, func(i, j int) int { return i - j }, 0, nil, t)
	testFoldl([]int{1, 2, 3}, "", func(s string, i int) string { return fmt.Sprintf("%d%s", i, s) }, "123", nil, t)
	testFoldl([]int{}, "", func(s string, i int) string { return "foo" }, "", nil, t)

	// Foldl should fail
	testFoldl(3, 0, 0, nil, foldlSliceError, t)
	testFoldl([]int{}, 0, 0, nil, foldlFunctionError, t)
	testFoldl([]int{}, 0, func(i, j, k int) int { return 0 }, nil, foldlTypeError, t)
	testFoldl([]int{}, 0, func(i, j int) (int, int) { return 0, 0 }, nil, foldlTypeError, t)
	testFoldl([]int{}, false, func(i, j int) bool { return false }, nil, foldlTypeError, t)
	testFoldl([]int{}, 0, func(b bool, i int) bool { return false }, nil, foldlZeroError, t)
	testFoldl([]int{}, 0, func(i, j int) bool { return false }, nil, foldlTypeError, t)
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

func TestFind(t *testing.T) {
	// Find should succeed
	testFind([]int{1, 2, 3}, func(i int) bool { return i%2 == 0 }, 2, nil, t)
	testFind([]int{1, 2, 3}, func(i int) bool { return i > 4 }, nil, nil, t)
	testFind([]bool{true, false, true}, func(b bool) bool { return b }, true, nil, t)
	testFind([]int{}, func(i int) bool { return true }, nil, nil, t)

	// Find should fail
	testFind(3, nil, nil, findSliceError, t)
	testFind([]int{1, 2, 3}, 3, nil, findFunctionError, t)
	testFind([]int{1, 2, 3}, func(b bool) bool { return b }, nil, findTypeError, t)
	testFind([]int{1, 2, 3}, func(i int) int { return i }, nil, findTypeError, t)
	testFind([]int{1, 2, 3}, func(i, j int) bool { return i == j }, nil, findTypeError, t)
}

func testFind(slc, pred, target interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Find(slc, pred)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return value %v; got %v", target, ret)
	}
}

func TestFindIndex(t *testing.T) {
	// FindIndex should succeed
	testFindIndex([]int{1, 2, 3}, func(i int) bool { return i%2 == 0 }, 1, nil, t)
	testFindIndex([]int{1, 2, 3}, func(i int) bool { return i > 4 }, -1, nil, t)
	testFindIndex([]bool{true, false, true}, func(b bool) bool { return b }, 0, nil, t)
	testFindIndex([]int{}, func(i int) bool { return true }, -1, nil, t)

	// FindIndex should fail
	testFindIndex(3, nil, -1, findIndexSliceError, t)
	testFindIndex([]int{1, 2, 3}, 3, -1, findIndexFunctionError, t)
	testFindIndex([]int{1, 2, 3}, func(b bool) bool { return b }, -1, findIndexTypeError, t)
	testFindIndex([]int{1, 2, 3}, func(i int) int { return i }, -1, findIndexTypeError, t)
	testFindIndex([]int{1, 2, 3}, func(i, j int) bool { return i == j }, -1, findIndexTypeError, t)
}

func testFindIndex(slc, pred interface{}, target int, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := FindIndex(slc, pred)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

func TestSome(t *testing.T) {
	// Some should succeed
	testSome([]int{1, 2, 3}, func(i int) bool { return i%2 == 0 }, true, nil, t)
	testSome([]int{1, 2, 3}, func(i int) bool { return i > 4 }, false, nil, t)
	testSome([]bool{true, false, true}, func(b bool) bool { return b }, true, nil, t)
	testSome([]int{}, func(i int) bool { return true }, false, nil, t)

	// Some should fail
	testSome(3, nil, false, someSliceError, t)
	testSome([]int{1, 2, 3}, 3, false, someFunctionError, t)
	testSome([]int{1, 2, 3}, func(b bool) bool { return b }, false, someTypeError, t)
	testSome([]int{1, 2, 3}, func(i int) int { return i }, false, someTypeError, t)
	testSome([]int{1, 2, 3}, func(i, j int) bool { return i == j }, false, someTypeError, t)
}

func testSome(slc, pred interface{}, target bool, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Some(slc, pred)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

func TestEvery(t *testing.T) {
	// Some should succeed
	testEvery([]int{1, 2, 3}, func(i int) bool { return i%2 == 0 }, false, nil, t)
	testEvery([]int{1, 2, 3}, func(i int) bool { return i > 0 }, true, nil, t)
	testEvery([]int{1, 2, 3}, func(i int) bool { return i > 4 }, false, nil, t)
	testEvery([]bool{true, false, true}, func(b bool) bool { return b }, false, nil, t)
	testEvery([]bool{true, true, true}, func(b bool) bool { return b }, true, nil, t)
	testEvery([]int{}, func(i int) bool { return true }, true, nil, t)

	// Some should fail
	testEvery(3, nil, false, everySliceError, t)
	testEvery([]int{1, 2, 3}, 3, false, everyFunctionError, t)
	testEvery([]int{1, 2, 3}, func(b bool) bool { return b }, false, everyTypeError, t)
	testEvery([]int{1, 2, 3}, func(i int) int { return i }, false, everyTypeError, t)
	testEvery([]int{1, 2, 3}, func(i, j int) bool { return i == j }, false, everyTypeError, t)
}

func testEvery(slc, pred interface{}, target bool, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Every(slc, pred)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

func TestCount(t *testing.T) {
	// Count should succeed
	testCount([]int{1, 2, 3}, func(i int) bool { return i%2 == 0 }, 1, nil, t)
	testCount([]int{1, 2, 3}, func(i int) bool { return i > 4 }, 0, nil, t)
	testCount([]bool{true, false, true}, func(b bool) bool { return b }, 2, nil, t)
	testCount([]bool{false, false, false}, func(b bool) bool { return b }, 0, nil, t)
	testCount([]int{}, func(i int) bool { return true }, 0, nil, t)

	// Count should fail
	testCount(3, nil, 0, countSliceError, t)
	testCount([]int{1, 2, 3}, 3, 0, countFunctionError, t)
	testCount([]int{1, 2, 3}, func(b bool) bool { return b }, 0, countTypeError, t)
	testCount([]int{1, 2, 3}, func(i int) int { return i }, 0, countTypeError, t)
	testCount([]int{1, 2, 3}, func(i, j int) bool { return i == j }, 0, countTypeError, t)
	testCount([]int{1, 2, 3}, func(i int) (bool, bool) { return true, true }, 0, countTypeError, t)
}

func testCount(slc, pred interface{}, target int, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Count(slc, pred)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

func TestMax(t *testing.T) {
	// Max should succeed
	testMax([]int{1, 2, 3}, func(i, j int) bool { return i < j }, 3, nil, t)
	testMax([]int{1, 2, 3}, func(i, j int) bool { return i > j }, 1, nil, t)
	testMax([]int{1, 2, 3}, func(i, j int) bool { return true }, 3, nil, t)
	testMax([]int{}, func(i, j int) bool { return true }, nil, nil, t)

	// Max should fail
	testMax(3, nil, nil, maxSliceError, t)
	testMax([]int{}, 3, nil, maxFunctionError, t)
	testMax([]int{}, func(i, j bool) bool { return true }, nil, maxTypeError, t)
	testMax([]int{}, func(i int, b bool) bool { return true }, nil, maxTypeError, t)
	testMax([]int{}, func(b bool, j int) bool { return true }, nil, maxTypeError, t)
	testMax([]int{}, func(i, j int) int { return 0 }, nil, maxTypeError, t)
	testMax([]int{}, func(i, j, k int) bool { return true }, nil, maxTypeError, t)
	testMax([]int{}, func(i, j int) (bool, bool) { return true, true }, nil, maxTypeError, t)
}

func testMax(slc, greater, target interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Max(slc, greater)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

func TestMin(t *testing.T) {
	// Max should succeed
	testMin([]int{1, 2, 3}, func(i, j int) bool { return i < j }, 1, nil, t)
	testMin([]int{1, 2, 3}, func(i, j int) bool { return i > j }, 3, nil, t)
	testMin([]int{1, 2, 3}, func(i, j int) bool { return true }, 3, nil, t)
	testMin([]int{}, func(i, j int) bool { return true }, nil, nil, t)
	testMin([]int{1}, func(i, j int) bool { return true }, 1, nil, t)

	// Max should fail
	testMin(3, nil, nil, minSliceError, t)
	testMin([]int{}, 3, nil, minFunctionError, t)
	testMin([]int{}, func(i, j bool) bool { return true }, nil, minTypeError, t)
	testMin([]int{}, func(i int, b bool) bool { return true }, nil, minTypeError, t)
	testMin([]int{}, func(b bool, j int) bool { return true }, nil, minTypeError, t)
	testMin([]int{}, func(i, j int) int { return 0 }, nil, minTypeError, t)
	testMin([]int{}, func(i, j, k int) bool { return true }, nil, minTypeError, t)
	testMin([]int{}, func(i, j int) (bool, bool) { return true, true }, nil, minTypeError, t)
}

func testMin(slc, greater, target interface{}, err interface{}, t *testing.T) {
	defer func() {
		r := recover()
		if !reflect.DeepEqual(r, err) {
			t.Errorf("Expected error %v; got %v", err, r)
		}
	}()

	ret := Min(slc, greater)
	if !reflect.DeepEqual(target, ret) {
		t.Errorf("Expected return index %v; got %v", target, ret)
	}
}

// Since we don't get to see the strings written
// as literals anywhere, do this so we can double-check
// that the error strings were composed properly.
//
// Capitalize to make it actually run
// (it's annoying to have it run every time).
func testVerifyErrorStrings(t *testing.T) {
	toPrint := []string{
		mapSliceError,
		mapFunctionError,
		mapTypeError,

		filterSliceError,
		filterFunctionError,
		filterTypeError,

		rejectSliceError,
		rejectFunctionError,
		rejectTypeError,

		foldrSliceError,
		foldrFunctionError,
		foldrTypeError,
		foldrZeroError,

		foldlSliceError,
		foldlFunctionError,
		foldlTypeError,
		foldlZeroError,

		findSliceError,
		findFunctionError,
		findTypeError,

		findIndexSliceError,
		findIndexFunctionError,
		findIndexTypeError,

		someSliceError,
		someFunctionError,
		someTypeError,

		everySliceError,
		everyFunctionError,
		everyTypeError,

		countSliceError,
		countFunctionError,
		countTypeError,

		maxSliceError,
		maxFunctionError,
		maxTypeError,

		minSliceError,
		minFunctionError,
		minTypeError,
	}
	fmt.Println("Error strings:")
	for _, s := range toPrint {
		fmt.Println(s)
	}
}
