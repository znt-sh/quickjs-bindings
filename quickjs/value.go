package quickjs

import ivalue "github.com/znt-sh/quickjs-bindings/internal/value"

type Value struct {
	inner ivalue.Value
	ctx   *Context
}

func wrapValue(ctx *Context, inner ivalue.Value) Value {
	return Value{inner: inner, ctx: ctx}
}

func (v Value) Dup() Value {
	return Value{inner: v.inner.Dup(), ctx: v.ctx}
}

func (v *Value) Free() {
	if v == nil {
		return
	}
	v.inner.Free()
}

func (v Value) IsNull() bool {
	return v.inner.IsNull()
}

func (v Value) IsUndefined() bool {
	return v.inner.IsUndefined()
}

func (v Value) IsBool() bool {
	return v.inner.IsBool()
}

func (v Value) IsNumber() bool {
	return v.inner.IsNumber()
}

func (v Value) IsString() bool {
	return v.inner.IsString()
}

func (v Value) IsObject() bool {
	return v.inner.IsObject()
}

func (v Value) IsArray() bool {
	return v.inner.IsArray()
}

func (v Value) IsFunction() bool {
	return v.inner.IsFunction()
}
