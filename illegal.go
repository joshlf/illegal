// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package illegal provides runtime support for operations
// which are disallowed by the compiler.
package illegal

import (
	"reflect"
)

// Convenience reflect.Type values
var (
	InterfaceType reflect.Type
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
		panic("illegal.FuncEqual: passed non-function value")
	}
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

// Convert each element in slc to example's
// underlying type, and return the conversions in a new slice.
//
// ConvertSlice panics if slc is not a slice value,
// or if the conversion is illegal.
//
// If ConvertSlice returns without panicking,
// the return value's underlying value is
// guaranteed to be a slice, and the element
// type is guaranteed to be of the same type
// as example.
func ConvertSlice(slc, example interface{}) interface{} {
	defer func() {
		r := recover()
		if r != nil {
			str, ok := r.(string)
			if !ok {
				panic("illegal: internal error: recovered from non-string panic in ConvertSlice")
			}
			panic("illegal.ConvertSlice: " + str)
		}
	}()
	return convertSliceType(slc, reflect.TypeOf(example))
}

// Convert each element in slc to the given type,
// and return the conversions in a new slice.
//
// ConvertSlice panics if slc is not a slice value,
// or if the conversion is illegal.
//
// If ConvertSliceType returns without panicking,
// the return value's underlying value is
// guaranteed to be a slice, and the element
// type is guaranteed to be of the same type
// as typ.
func ConvertSliceType(slc interface{}, typ reflect.Type) interface{} {
	defer func() {
		r := recover()
		if r != nil {
			str, ok := r.(string)
			if !ok {
				panic("illegal: internal error: recovered from non-string panic in ConvertSliceType")
			}
			panic("illegal.ConvertSliceType: " + str)
		}
	}()
	return convertSliceType(slc, typ)
}

// Since both ConvertSlice and ConvertSliceType
// call this, and need to be able to name
// their own panics ("illegal.ConvertSlice: "
// vs "illegal.ConvertSliceType: "), this function
// only panics using the message itself, letting
// the exported functions which call it prepend
// the "illegal.FuncName: " prefix.
func convertSliceType(slc interface{}, typ reflect.Type) interface{} {
	slice := reflect.ValueOf(slc)
	if slice.Kind() != reflect.Slice {
		panic("passed non-slice value")
	}

	// reflect.Value.Convert panics if it
	// is called on a value which cannot
	// be converted to the given type. We
	// don't want to display their panic
	// message, so we create our own.
	//
	// Note that we don't do this check
	// ahead of time, letting the first
	// conversion take place, since
	// it is more performant, and the
	// normal use case is that the user
	// passes the correct arguments (so
	// such a check would, in theory,
	// never actually panic anyway).
	defer func() {
		r := recover()
		if r != nil {
			panic("cannot convert type " + slice.Type().Elem().String() + " to " + typ.String())
		}
	}()

	ret := reflect.MakeSlice(reflect.SliceOf(typ), slice.Len(), slice.Cap())
	args := make([]reflect.Value, 1)
	for i := 0; i < slice.Len(); i++ {
		args[0] = slice.Index(i)
		ret.Index(i).Set(args[0].Convert(typ))
	}

	// If slice.Len() == 0, then no conversions have
	// been attempted, which means that it's possible
	// that the conversion is illegal, but the function
	// hasn't panicked yet. Thus, check explicitly.
	if slice.Len() == 0 {
		if !slice.Type().ConvertibleTo(typ) {
			panic(0) // This panic will be caught by recover, so its value is irrelevant
		}
	}

	return ret.Interface()
}

func init() {
	// Credit to http://golang.org/src/pkg/net/rpc/server.go?s=4244:4436#L145 (build version go1.1.2)
	InterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
}
