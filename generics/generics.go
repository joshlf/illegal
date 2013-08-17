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
		panic("illegal: passed non-slice value to Map")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Map")
	}

	slcType := slice.Type()
	fType := f.Type()

	// f must take a single parameter of the same type as
	// the given slice, and return a single result
	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != slcType.Elem() {
		panic("illegal: function type and slice type do not match in call to Map(slc []T, pred func(T) U) []U")
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
		panic("illegal: passed non-slice value to Filter")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Filter")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Filter(slc []T, pred func(T) bool) []T")
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
		panic("illegal: passed non-slice value to Reject")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Reject")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Reject(slc []T, pred func(T) bool) []T")
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

//	func Foldr(slc []T, zero U, pred func(T, U) U) U
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
		panic("illegal: passed non-slice value to Foldr")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Foldr")
	}

	z := reflect.ValueOf(zero)

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != fType.Out(0) {
		panic("illegal: function type and slice type do not match in call to Foldr(slc []T, zero U, pred func(T, U) U) U")
	}

	// It's possible to have a valid function
	// (that is, func(A, B)B) and have the type
	// of zero not be equal to B
	if fType.Out(0) != z.Type() {
		panic("illegal: zero type and function return type do not match in call to Foldr(slc []T, zero U, pred func(T, U) U) U")
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
		panic("illegal: passed non-slice value to Foldl")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Foldl")
	}

	z := reflect.ValueOf(zero)

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(1) != elemType || fType.In(0) != fType.Out(0) {
		panic("illegal: function type and slice type do not match in call to Foldl(slc []T, zero U, pred func(U, T) U) U")
	}

	// It's possible to have a valid function
	// (that is, func(B, A)B) and have the type
	// of zero not be equal to B
	if fType.Out(0) != z.Type() {
		panic("illegal: zero type and function return type do not match in call to Foldl(slc []T, zero U, pred func(U, T) U) U")
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
		panic("illegal: passed non-slice value to Find")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Find")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Find(slc []T, pred func(T) bool) T")
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
		panic("illegal: passed non-slice value to FindIndex")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to FindIndex")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to FindIndex(slc []T, pred func(T) bool) int")
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
		panic("illegal: passed non-slice value to Some")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Some")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Some(slc []T, pred func(T) bool) bool")
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
		panic("illegal: passed non-slice value to Every")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Every")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Every(slc []T, pred func(T) bool) bool")
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
		panic("illegal: passed non-slice value to Count")
	}

	f := reflect.ValueOf(pred)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Count")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Count(slc []T, pred func(T) bool) int")
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
		panic("illegal: passed non-slice value to Max")
	}

	f := reflect.ValueOf(less)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Max")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Max(slc []T, less func(T, T) bool) T")
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
		panic("illegal: passed non-slice value to Min")
	}

	f := reflect.ValueOf(less)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Min")
	}

	slcType := slice.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Min(slc []T, less func(T, T) bool) T")
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
