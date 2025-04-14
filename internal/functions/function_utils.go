package functions

import (
	"fmt"
	"reflect"
	"unicode"
)

// Validation functions

func FunctionNameValid(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r == '_':
		case i == 0 && !unicode.IsLetter(r):
			return false
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			return false
		}
	}
	return true
}

func FunctionValid(fn any) error {
	// Check that fn is a function
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("not a function")
	}

	// We allow functions with 1 result or 2 results where the second is an error.
	typ := v.Type()
	switch numOut := typ.NumOut(); {
	case numOut == 1:
		return nil
	case numOut == 2 && typ.Out(1) == reflect.TypeFor[error]():
		return nil
	case numOut == 2:
		return fmt.Errorf("invalid function signature - second return value should be error; is %s", typ.Out(1))
	default:
		return fmt.Errorf("function has %d return values; should be 1 or 2", typ.NumOut())
	}
}
