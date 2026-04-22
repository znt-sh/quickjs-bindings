package quickjs

import (
	"fmt"

	icontext "github.com/znt-sh/quickjs-bindings/internal/context"
	ivalue "github.com/znt-sh/quickjs-bindings/internal/value"
)

type Function func(ctx *Context, this Value, args []Value) (Value, error)

func (c *Context) RegisterFunction(name string, length int, fn Function) error {
	if c == nil || c.inner == nil {
		return &Error{Message: "quickjs: context is nil"}
	}
	if fn == nil {
		return &Error{Message: "quickjs: function is nil"}
	}

	return wrapError(c.inner.RegisterFunction(name, length, func(_ *icontext.Context, this ivalue.Value, args []ivalue.Value) (ivalue.Value, error) {
		publicArgs := make([]Value, len(args))
		for i := range args {
			publicArgs[i] = wrapValue(c, args[i])
		}

		out, callErr := fn(c, wrapValue(c, this), publicArgs)
		if callErr != nil {
			return ivalue.Value{}, callErr
		}
		return out.inner, nil
	}))
}

func ThrowTypeError(message string) error {
	return &Error{Message: fmt.Sprintf("TypeError: %s", message)}
}
