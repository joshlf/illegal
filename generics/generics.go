// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package generics implements common generic functions.
//
// Since many of the types are not specified at compile-time,
// each function's documentation will begin with its equivalent
// generic type definition. For example, for the Map function,
// the type definition is:
//
//	func Map(slc []T, pred func(T) U) []U
//
// For all functions, if the types of the arguments do not match
// the generic argument types of the function, it will cause a
// runtime panic.
//
// Except in certain documented cases, the documented return types
// are guaranteed to be valid. Thus, type assertions are guaranteed
// to succeed.
// For example, for func Identity(x T) T:
//
//	yVal := Identity(1)
//	y, _ := yVal.(int) // guaranteed to succeed
package generics

import (
	"reflect"
)

// Pre-computed type literals
var (
	boolType = reflect.TypeOf(bool(true))
)

//	func Identity(x T) T
//
// Identity returns its argument.
func Identity(x interface{}) interface{} { return x }

//	func Map(slc []T, pred func(T) U) []U
//
// Map applies pred to each element of slc
// successively, and returns the results.
func Map(slc, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(mapSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(mapFunctionError)
	}

	slcType := slice.Type()
	fType := f.Type()

	// f must take a single parameter of the same type as
	// the given slice, and return a single result
	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != slcType.Elem() {
		panic(mapTypeError)
	}

	ret := reflect.MakeSlice(reflect.SliceOf(fType.Out(0)), slice.Len(), slice.Cap())

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		ret.Index(i).Set(f.Call(args)[0])
	}

	return ret.Interface()
}

//	func Filter(slc []T, pred func(T) bool) []T
//
// Filter applies pred to each element of slc,
// and returns those elements for which pred
// returned true. The returned slice will be
// only as long as it needs to to store all
// "true" elements, not necessarily as long
// as slc.
func Filter(slc, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(filterSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(filterFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(filterTypeError)
	}

	ret := reflect.MakeSlice(slcType, 0, 0)

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			ret = reflect.Append(ret, args[0])
		}
	}

	return ret.Interface()
}

//	func Reject(slc []T, pred func(T) bool) []T
//
// Reject applies pred to each element of slc,
// and returns those elements for which pred
// returned false. The returned slice will be
// only as long as it needs to to store all
// "false" elements, not necessarily as long
// as slc.
func Reject(slc, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(rejectSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(rejectFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(rejectTypeError)
	}

	ret := reflect.MakeSlice(slcType, 0, 0)

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if !f.Call(args)[0].Bool() {
			ret = reflect.Append(ret, args[0])
		}
	}

	return ret.Interface()
}

//	func foldl(slc []T, zero U, pred func(T, U) U) U
//
// Foldr applies pred to each element of slc, using
// the previous call's return value as its second
// argument. In other words, it does:
//
//	tmp := pred(slc[0], zero)
//	tmp = pred(slc[1], tmp)
//	tmp = pred(slc[2], tmp)
//	...
//	return tmp
func Foldr(slc, zero, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(foldrSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(foldrFunctionError)
	}

	z := reflect.ValueOf(zero)

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != fType.Out(0) {
		panic(foldrTypeError)
	}

	// It's possible to have a valid function
	// (that is, func(A, B)B) and have the type
	// of zero not be equal to B
	if fType.Out(0) != z.Type() {
		panic(foldrZeroError)
	}

	args := make([]reflect.Value, 2)
	args[1] = z
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		args[1] = f.Call(args)[0]
	}

	return args[1].Interface()
}

//	func Foldl(slc []T, zero U, pred func(U, T) U) U
//
// Foldr applies pred to each element of slc in reverse
// order, using the previous call's return value as its
// first argument. In other words, it does:
//
//	tmp := pred(zero, slc[len(slc)-1])
//	tmp = pred(tmp, slc[len(slc)-2])
//	tmp = pred(tmp, slc[len(slc)-3])
//	...
//	return tmp
func Foldl(slc, zero, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(foldlSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(foldlFunctionError)
	}

	z := reflect.ValueOf(zero)

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(1) != elemType || fType.In(0) != fType.Out(0) {
		panic(foldlTypeError)
	}

	// It's possible to have a valid function
	// (that is, func(B, A)B) and have the type
	// of zero not be equal to B
	if fType.Out(0) != z.Type() {
		panic(foldlZeroError)
	}

	args := make([]reflect.Value, 2)
	args[0] = z
	for i := slice.Len() - 1; i > -1; i-- {
		args[1] = slice.Index(i)
		args[0] = f.Call(args)[0]
	}

	return args[0].Interface()
}

//	func Find(slc []T, pred func(T) bool) T
//
// Find applies pred to each element in slc,
// returning the first element for which pred
// returns true. If pred never returns true,
// Find returns a nil interface. Note that type
// assertions of the form
//
//	f := find([]int{1, 2, 3}, pred)
//	g, ok := f.(int)
//
// may fail, which breaks the contract which
// most other functions in this package obey.
func Find(slc, pred interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(findSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(findFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(findTypeError)
	}

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			return args[0].Interface()
		}
	}

	return nil
}

//	func FindIndex(slc []T, pred func(T) bool) int
//
// FindIndex applies pred to each element in slc,
// returning the index of the first element for which
// pred returns true. If pred never returns true,
// FindIndex returns -1.
func FindIndex(slc, pred interface{}) int {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(findIndexSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(findIndexFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(findIndexTypeError)
	}

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			return i
		}
	}

	return -1
}

//	func Some(slc []T, pred func(T) bool) bool
//
// Some applies pred to each element in slc.
// If any of those calls returns true, Contains
// returns true. Otherwise, it returns false.
func Some(slc, pred interface{}) bool {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(someSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(someFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(someTypeError)
	}

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			return true
		}
	}

	return false
}

//	func Every(slc []T, pred func(T) bool) bool
//
// Every applies pred to each element in slc.
// If any of those calls returns false, Every
// returns false. Otherwise, it returns true.
func Every(slc, pred interface{}) bool {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(everySliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(everyFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(everyTypeError)
	}

	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if !f.Call(args)[0].Bool() {
			return false
		}
	}

	return true
}

//	func Count(slc []T, pred func(T) bool) int
//
// Count applies pred to each element in slc,
// and returns the number of elements for which
// the call returned true.
func Count(slc, pred interface{}) int {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(countSliceError)
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic(countFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic(countTypeError)
	}

	ret := 0
	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			ret++
		}
	}

	return ret
}

//	func Max(slc []T, less func(T, T) bool) T
//
// Max finds the largest element in slc according
// to less (less(a, b) returns (a < b)).
// If len(slc) == 0, Max will return a nil interface,
// thus breaking the type assertion guarantee.
// However, so long as len(slc) > 0, the type
// assertion guarantee holds.
func Max(slc, less interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(maxSliceError)
	}

	f := reflect.ValueOf(less)
	if f.Kind() != reflect.Func {
		panic(maxFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != elemType || fType.Out(0) != boolType {
		panic(maxTypeError)
	}

	if slice.Len() == 0 {
		return nil
	}

	args := make([]reflect.Value, 2)
	args[0] = slice.Index(0)
	for i := 1; i < slice.Len(); i++ {
		args[1] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			args[0] = args[1]
		}
	}

	return args[0].Interface()
}

//	func Min(slc []T, less func(T, T) bool) T
//
// Min finds the smallest element in slc according
// to less (less(a, b) returns (a < b)).
// If len(slc) == 0, Min will return a nil interface,
// thus breaking the type assertion guarantee.
// However, so long as len(slc) > 0, the type
// assertion guarantee holds.
func Min(slc, less interface{}) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic(minSliceError)
	}

	f := reflect.ValueOf(less)
	if f.Kind() != reflect.Func {
		panic(minFunctionError)
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != elemType || fType.Out(0) != boolType {
		panic(minTypeError)
	}

	if slice.Len() == 0 {
		return nil
	}

	args := make([]reflect.Value, 2)
	args[1] = slice.Index(0)
	for i := 1; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		if f.Call(args)[0].Bool() {
			args[1] = args[0]
		}
	}

	return args[1].Interface()
}

var (
	// The same basic types of errors are used
	// over and over again, and must be checked
	// against in testing code. Writing out the
	// string literals every single time is tiresome,
	// and has led to consistency issues. Making
	// them variables creates a single point of
	// truth, and allows the error string scheme
	// to be edited /far/ more easily than if we
	// were using string literals everywhere.
	//
	// The general form of errors is:
	//     package.Function: error
	// So, for example,
	//     generics.Map: passed non-slice value
	sliceError        = "passed non-slice value"
	functionError     = "passed non-function value"
	typeError         = "function type and slice type do not match"
	zeroError         = "zero type and function return type do not match"
	packageNamePrefix = "generics."

	mapErrorPrefix   = packageNamePrefix + "Map: "
	mapSliceError    = mapErrorPrefix + sliceError
	mapFunctionError = mapErrorPrefix + functionError
	mapTypeError     = mapErrorPrefix + typeError

	filterErrorPrefix   = packageNamePrefix + "Filter: "
	filterSliceError    = filterErrorPrefix + sliceError
	filterFunctionError = filterErrorPrefix + functionError
	filterTypeError     = filterErrorPrefix + typeError

	rejectErrorPrefix   = packageNamePrefix + "Reject: "
	rejectSliceError    = rejectErrorPrefix + sliceError
	rejectFunctionError = rejectErrorPrefix + functionError
	rejectTypeError     = rejectErrorPrefix + typeError

	foldrErrorPrefix   = packageNamePrefix + "Foldr: "
	foldrSliceError    = foldrErrorPrefix + sliceError
	foldrFunctionError = foldrErrorPrefix + functionError
	foldrTypeError     = foldrErrorPrefix + typeError
	foldrZeroError     = foldrErrorPrefix + zeroError

	foldlErrorPrefix   = packageNamePrefix + "Foldl: "
	foldlSliceError    = foldlErrorPrefix + sliceError
	foldlFunctionError = foldlErrorPrefix + functionError
	foldlTypeError     = foldlErrorPrefix + typeError
	foldlZeroError     = foldlErrorPrefix + zeroError

	findErrorPrefix   = packageNamePrefix + "Find: "
	findSliceError    = findErrorPrefix + sliceError
	findFunctionError = findErrorPrefix + functionError
	findTypeError     = findErrorPrefix + typeError

	findIndexErrorPrefix   = packageNamePrefix + "FindIndex: "
	findIndexSliceError    = findIndexErrorPrefix + sliceError
	findIndexFunctionError = findIndexErrorPrefix + functionError
	findIndexTypeError     = findIndexErrorPrefix + typeError

	someErrorPrefix   = packageNamePrefix + "Some: "
	someSliceError    = someErrorPrefix + sliceError
	someFunctionError = someErrorPrefix + functionError
	someTypeError     = someErrorPrefix + typeError

	everyErrorPrefix   = packageNamePrefix + "Every: "
	everySliceError    = everyErrorPrefix + sliceError
	everyFunctionError = everyErrorPrefix + functionError
	everyTypeError     = everyErrorPrefix + typeError

	countErrorPrefix   = packageNamePrefix + "Count: "
	countSliceError    = countErrorPrefix + sliceError
	countFunctionError = countErrorPrefix + functionError
	countTypeError     = countErrorPrefix + typeError

	maxErrorPrefix   = packageNamePrefix + "Max: "
	maxSliceError    = maxErrorPrefix + sliceError
	maxFunctionError = maxErrorPrefix + functionError
	maxTypeError     = maxErrorPrefix + typeError

	minErrorPrefix   = packageNamePrefix + "Min: "
	minSliceError    = minErrorPrefix + sliceError
	minFunctionError = minErrorPrefix + functionError
	minTypeError     = minErrorPrefix + typeError
)
