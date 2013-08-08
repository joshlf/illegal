package illegal

import (
	"errors"
	"reflect"
)

var (
	ErrNotFunc       = errors.New("Argument was not a function")
	ErrNotSlice      = errors.New("Argument was not a slice")
	ErrWrongFuncType = errors.New("Argument function did not have the correct signature")
	ErrWrongZeroType = errors.New("Zero argument type did not match function return type")
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
// FuncEqual returns only one error, ErrNotFunc,
// if one or both of the arguments are not actually
// function pointers. Otherwise, the error is nil.
func FuncEqual(f1, f2 interface{}) (bool, error) {
	if reflect.TypeOf(f1).Kind() != reflect.Func || reflect.TypeOf(f2).Kind() != reflect.Func {
		return false, ErrNotFunc
	}
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer(), nil
}

// Map takes slice, of type []T, and fn, of type
// func(T)S, and returns a slice whose contents
// are the result of applying fn to each element
// of slice successively. The resulting slice
// has type []S, and has the same length as slice.
//
// It is an error to pass a non-slice value as slice
// (ErrNotSlice), or a non-function value as fn'
// (ErrNotFunc). It is an error for the argument type
// of fn to be anything but T, and it is an error for
// fn to return anything but 1 value (ErrWrongFuncType).
func Map(slice, fn interface{}) (interface{}, error) {
	slc := reflect.ValueOf(slice)
	if slc.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		return nil, ErrNotFunc
	}

	slcType := slc.Type()
	fType := f.Type()

	// f must take a single parameter of the same type as
	// the given slice, and return a single result
	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != slcType.Elem() {
		return nil, ErrWrongFuncType
	}

	ret := reflect.MakeSlice(reflect.SliceOf(fType.Out(0)), slc.Len(), slc.Cap())

	args := make([]reflect.Value, 1)
	for i := 0; i < slc.Len(); i++ {
		args[0] = slc.Index(i)
		ret.Index(i).Set(f.Call(args)[0])
	}

	return ret.Interface(), nil
}

// Map takes slice, of type []T, and fn, of type
// func(T)bool, and returns a slice whose contents
// are those elements, elem1, elem2, ... elemn,
// for which fn(elemi) == true. The resulting slice
// has type []T, and has length less than or equal
// to the length of slice.
//
// It is an error to pass a non-slice value as slice
// (ErrNotSlice), or a non-function value as fn'
// (ErrNotFunc). It is an error for the argument type
// of fn to be anything but T, and it is an error for
// fn to return anything but bool (ErrWrongFuncType).
func Filter(slice, fn interface{}) (interface{}, error) {
	slc := reflect.ValueOf(slice)
	if slc.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		return nil, ErrNotFunc
	}

	slcType := slc.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 1 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.Out(0) != boolType {
		return nil, ErrWrongFuncType
	}

	ret := reflect.MakeSlice(slcType, 0, 0)

	args := make([]reflect.Value, 1)
	for i := 0; i < slc.Len(); i++ {
		args[0] = slc.Index(i)
		if f.Call(args)[0].Bool() {
			ret = reflect.Append(ret, args[0])
		}
	}

	return ret.Interface(), nil
}

func Foldr(slice, zero, fn interface{}) (interface{}, error) {
	slc := reflect.ValueOf(slice)
	if slc.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	f := reflect.ValueOf(fn)
	if f.Kind() != reflect.Func {
		return nil, ErrNotFunc
	}

	z := reflect.ValueOf(zero)

	slcType := slc.Type()
	elemType := slcType.Elem()
	fType := f.Type()

	if fType.NumIn() != 2 || fType.NumOut() != 1 || fType.In(0) != elemType || fType.In(1) != fType.Out(0) {
		return nil, ErrWrongFuncType
	}

	// It's possible to have a valid function
	// (that is, func(A, B)B) and have the type
	// of zero not be equal to B
	if fType.Out(0) != z.Type() {
		return nil, ErrWrongZeroType
	}

	args := make([]reflect.Value, 2)
	args[1] = z
	for i := 0; i < slc.Len(); i++ {
		args[0] = slc.Index(i)
		args[1] = f.Call(args)[0]
	}

	return args[1].Interface(), nil
}
