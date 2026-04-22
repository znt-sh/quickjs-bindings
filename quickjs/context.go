package quickjs

import (
	"fmt"

	icontext "github.com/znt-sh/quickjs-bindings/internal/context"
	"github.com/znt-sh/quickjs-bindings/internal/convert"
)

const defaultEvalFilename = "<eval>"

type Context struct {
	inner *icontext.Context
}

func NewContext(runtime *Runtime) (*Context, error) {
	if runtime == nil || runtime.inner == nil {
		return nil, &Error{Message: "quickjs: runtime is nil"}
	}

	ctx, err := icontext.New(runtime.inner)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Context{inner: ctx}, nil
}

func (c *Context) Close() {
	if c == nil || c.inner == nil {
		return
	}
	c.inner.Close()
}

func (c *Context) Eval(code string) (Value, error) {
	return c.EvalFile(code, defaultEvalFilename)
}

func (c *Context) EvalFile(code, filename string) (Value, error) {
	out, err := c.inner.Eval(code, filename)
	if err != nil {
		return Value{}, wrapError(err)
	}
	return wrapValue(c, out), nil
}

func (c *Context) ToString(v Value) (string, error) {
	out, err := c.inner.ToString(v.inner)
	return out, wrapError(err)
}

func (c *Context) ToInt32(v Value) (int32, error) {
	out, err := c.inner.ToInt32(v.inner)
	return out, wrapError(err)
}

func (c *Context) ToFloat64(v Value) (float64, error) {
	out, err := c.inner.ToFloat64(v.inner)
	return out, wrapError(err)
}

func (c *Context) ToBool(v Value) (bool, error) {
	out, err := c.inner.ToBool(v.inner)
	return out, wrapError(err)
}

func (c *Context) NewInt(v int32) Value {
	return wrapValue(c, c.inner.NewInt32(v))
}

func (c *Context) NewFloat(v float64) Value {
	return wrapValue(c, c.inner.NewFloat64(v))
}

func (c *Context) NewBool(v bool) Value {
	return wrapValue(c, c.inner.NewBool(v))
}

func (c *Context) NewString(v string) Value {
	return wrapValue(c, c.inner.NewString(v))
}

func (c *Context) NewNull() Value {
	return wrapValue(c, c.inner.NewNull())
}

func (c *Context) NewUndefined() Value {
	return wrapValue(c, c.inner.NewUndefined())
}

func (c *Context) NewObject() Value {
	return wrapValue(c, c.inner.NewObject())
}

func (c *Context) NewArray() Value {
	return wrapValue(c, c.inner.NewArray())
}

func (c *Context) GetProperty(obj Value, prop string) (Value, error) {
	out, err := c.inner.GetPropertyString(obj.inner, prop)
	if err != nil {
		return Value{}, wrapError(err)
	}
	return wrapValue(c, out), nil
}

func (c *Context) SetProperty(obj Value, prop string, v Value) error {
	return wrapError(c.inner.SetPropertyString(obj.inner, prop, v.inner))
}

func (c *Context) GetIndex(obj Value, idx uint32) (Value, error) {
	out, err := c.inner.GetPropertyUint32(obj.inner, idx)
	if err != nil {
		return Value{}, wrapError(err)
	}
	return wrapValue(c, out), nil
}

func (c *Context) SetIndex(obj Value, idx uint32, v Value) error {
	return wrapError(c.inner.SetPropertyUint32(obj.inner, idx, v.inner))
}

func (c *Context) ToValue(in any) (Value, error) {
	out, err := convert.ToJS(c.inner, in)
	if err != nil {
		return Value{}, wrapError(err)
	}
	return wrapValue(c, out), nil
}

func (c *Context) NewObjectFromMap(in map[string]any) (Value, error) {
	return c.ToValue(in)
}

func (c *Context) NewArrayFromSlice(in []any) (Value, error) {
	return c.ToValue(in)
}

func (c *Context) FromValue(v Value) (any, error) {
	out, err := convert.FromJS(c.inner, v.inner)
	return out, wrapError(err)
}

func (c *Context) MustString(v Value) string {
	out, err := c.ToString(v)
	if err != nil {
		return ""
	}
	return out
}

func (c *Context) MustFloat(v Value) float64 {
	out, err := c.ToFloat64(v)
	if err != nil {
		return 0
	}
	return out
}

func (c *Context) MustBool(v Value) bool {
	out, err := c.ToBool(v)
	if err != nil {
		return false
	}
	return out
}

func (c *Context) GetAny(obj Value, prop string) (any, error) {
	v, err := c.GetProperty(obj, prop)
	if err != nil {
		return nil, err
	}
	defer v.Free()
	return c.FromValue(v)
}

func (c *Context) GetString(obj Value, prop string) (string, error) {
	v, err := c.GetProperty(obj, prop)
	if err != nil {
		return "", err
	}
	defer v.Free()
	return c.ToString(v)
}

func (c *Context) GetFloat(obj Value, prop string) (float64, error) {
	v, err := c.GetProperty(obj, prop)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	return c.ToFloat64(v)
}

func (c *Context) GetInt(obj Value, prop string) (int32, error) {
	v, err := c.GetProperty(obj, prop)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	return c.ToInt32(v)
}

func (c *Context) GetBool(obj Value, prop string) (bool, error) {
	v, err := c.GetProperty(obj, prop)
	if err != nil {
		return false, err
	}
	defer v.Free()
	return c.ToBool(v)
}

func (c *Context) SetAny(obj Value, prop string, in any) error {
	v, err := c.ToValue(in)
	if err != nil {
		return err
	}
	defer v.Free()
	return c.SetProperty(obj, prop, v)
}

func (c *Context) ObjectToMap(obj Value, keys ...string) (map[string]any, error) {
	if !obj.IsObject() {
		return nil, &Error{Message: "quickjs: value is not an object"}
	}
	out := make(map[string]any, len(keys))
	for _, key := range keys {
		v, err := c.GetProperty(obj, key)
		if err != nil {
			return nil, err
		}
		item, convErr := c.FromValue(v)
		v.Free()
		if convErr != nil {
			return nil, convErr
		}
		out[key] = item
	}
	return out, nil
}

func (c *Context) AssertArgsAtLeast(args []Value, n int) error {
	if len(args) < n {
		return &Error{Message: fmt.Sprintf("quickjs: expected at least %d args, got %d", n, len(args))}
	}
	return nil
}
