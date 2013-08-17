// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package illegal

import (
	"reflect"
)

// Pre-computed type literals
var (
	boolType = reflect.TypeOf(bool(true))
)

// FuncEqual figures out if two function pointers
// reference the same function. For two closures
// created from the same original function, it will
// say that they are equal.
//
// FuncEqual also works on methods. Note that if
// two interface variables of the same interface
// type but of two separate concrete types are
// used to get a function pointer, that pointer
// will point at the interface methods, not the
// methods associated with the concrete types,
// so they will register as equal.
//
// FuncEqual panics if either argument is not
// a function.
func FuncEqual(f1, f2 interface{}) bool {
	if reflect.TypeOf(f1).Kind() != reflect.Func || reflect.TypeOf(f2).Kind() != reflect.Func {
		panic("illegal: passed non-function value to FuncEqual")
	}
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

// Map takes slice, of type []T, and fn, of type
// func(T) U, and returns a slice whose contents
// are the result of applying fn to each element
// of slice successively. The resulting slice
// has type []U, and has the same length as slice.
//
// Map panics if slice is not actually a slice,
// or if fn is not actually a function. Map panics
// if the argument type of fn is anything but T,
// or if fn does not return exactly 1 value.
func Map(slice, fn interface{}) interface{} {
	slc := reflect.ValueOf(slice)
	if slc.Kind() != reflect.Slice {
		panic("illegal: passed non-slice value to Map")
	}

	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Map")
	}

	slcType := slc.Type()
	fType := f.Type()

	// f must take a single parameter of the same type as
	// the given slice, and return a single result
	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != slcType.Elem() {
		panic("illegal: function type and slice type do not match in call to Map(slice []T, fn func(T) U) []U")
	}

	ret := reflect.MakeSlice(reflect.SliceOf(fType.Out(0)), slc.Len(), slc.Cap())

	args := make([]reflect.Value, 1)
	for i := 0; i < slc.Len(); i++ {
		args[0] = slc.Index(i)
		ret.Index(i).Set(f.Call(args)[0])
	}

	return ret.Interface()
}

// Map takes slice, of type []T, and fn, of type
// func(T) bool, and returns a slice whose contents
// are those elements, elem1, elem2, ... elemn,
// for which fn(elemi) == true. The resulting slice
// has type []T, and has length less than or equal
// to the length of slice.
//
// Filter panics if slice is not actually a slice,
// or if fn is not actually a function. Filter panics
// if the argument type of fn is anything but T, or if
// fn returns anything but bool.
func Filter(slice, fn interface{}) interface{} {
	slc := reflect.ValueOf(slice)
	if slc.Kind() != reflect.Slice {
		panic("illegal: passed non-slice value to Filter")
	}

	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		panic("illegal: passed non-function value to Filter")
	}

	slcType := slc.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		panic("illegal: function type and slice type do not match in call to Filter(slice []T, fn func(T) bool) []T")
	}

	ret := reflect.MakeSlice(slcType, 0, 0)

	args := make([]reflect.Value, 1)
	for i := 0; i < slc.Len(); i++ {
		args[0] = slc.Index(i)
		if f.Call(args)[0].Bool() {
			ret = reflect.Append(ret, args[0])
		}
	}

	return ret.Interface()
}
