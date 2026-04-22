package quickjs

import icgo "github.com/znt-sh/quickjs-bindings/internal/cgo"

type Error struct {
	Message string
	Stack   string
}

func (e *Error) Error() string {
	if e == nil {
		return "quickjs: unknown error"
	}
	if e.Stack != "" {
		return e.Message + "\n" + e.Stack
	}
	return e.Message
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	if jsErr, ok := err.(*icgo.JSError); ok {
		return &Error{Message: jsErr.Message, Stack: jsErr.Stack}
	}
	return &Error{Message: err.Error()}
}
