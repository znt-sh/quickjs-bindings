package value

import "github.com/znt-sh/quickjs-bindings/internal/cgo"

type Value struct {
	handle cgo.ValueHandle
}

func Wrap(handle cgo.ValueHandle) Value {
	return Value{handle: handle}
}

func (v Value) Handle() cgo.ValueHandle {
	return v.handle
}

func (v Value) Dup() Value {
	return Value{handle: v.handle.Dup()}
}

func (v *Value) Free() {
	if v == nil {
		return
	}
	v.handle.Free()
}

func (v Value) IsNull() bool {
	return v.handle.IsNull()
}

func (v Value) IsUndefined() bool {
	return v.handle.IsUndefined()
}

func (v Value) IsBool() bool {
	return v.handle.IsBool()
}

func (v Value) IsNumber() bool {
	return v.handle.IsNumber()
}

func (v Value) IsString() bool {
	return v.handle.IsString()
}

func (v Value) IsObject() bool {
	return v.handle.IsObject()
}

func (v Value) IsArray() bool {
	return v.handle.IsArray()
}

func (v Value) IsFunction() bool {
	return v.handle.IsFunction()
}
