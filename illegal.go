// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package illegal provides runtime support for operations
// which are disallowed by the compiler.
package illegal

import (
	"reflect"
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
