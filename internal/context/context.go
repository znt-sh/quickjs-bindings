package context

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/znt-sh/quickjs-bindings/internal/cgo"
	iruntime "github.com/znt-sh/quickjs-bindings/internal/runtime"
	"github.com/znt-sh/quickjs-bindings/internal/value"
)

type Context struct {
	handle      *cgo.ContextHandle
	callbackIDs []int32
	ownerGID    uint64
	closed      bool
	mu          sync.RWMutex
}

func New(rt *iruntime.Runtime) (*Context, error) {
	if rt == nil || rt.Handle() == nil {
		return nil, fmt.Errorf("quickjs: runtime is nil")
	}
	h, err := cgo.NewContext(rt.Handle())
	if err != nil {
		return nil, err
	}
	return &Context{handle: h, ownerGID: currentGID()}, nil
}

func (c *Context) Close() {
	if c == nil {
		return
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	for _, id := range c.callbackIDs {
		cgo.UnregisterCallback(id)
	}
	c.callbackIDs = nil
	h := c.handle
	c.mu.Unlock()

	if h != nil {
		h.Close()
	}
}

func (c *Context) Eval(code, filename string) (value.Value, error) {
	if err := c.guard(); err != nil {
		return value.Value{}, err
	}
	v, err := c.handle.Eval(code, filename, cgo.EvalTypeGlobal)
	if err != nil {
		return value.Value{}, err
	}
	return value.Wrap(v), nil
}

func (c *Context) ToString(v value.Value) (string, error) {
	if err := c.guardValue(v); err != nil {
		return "", err
	}
	return c.handle.ToString(v.Handle())
}

func (c *Context) ToInt32(v value.Value) (int32, error) {
	if err := c.guardValue(v); err != nil {
		return 0, err
	}
	return c.handle.ToInt32(v.Handle())
}

func (c *Context) ToFloat64(v value.Value) (float64, error) {
	if err := c.guardValue(v); err != nil {
		return 0, err
	}
	return c.handle.ToFloat64(v.Handle())
}

func (c *Context) ToBool(v value.Value) (bool, error) {
	if err := c.guardValue(v); err != nil {
		return false, err
	}
	return c.handle.ToBool(v.Handle())
}

func (c *Context) NewInt32(v int32) value.Value {
	return value.Wrap(c.handle.NewInt32(v))
}

func (c *Context) NewFloat64(v float64) value.Value {
	return value.Wrap(c.handle.NewFloat64(v))
}

func (c *Context) NewBool(v bool) value.Value {
	return value.Wrap(c.handle.NewBool(v))
}

func (c *Context) NewString(v string) value.Value {
	return value.Wrap(c.handle.NewString(v))
}

func (c *Context) NewNull() value.Value {
	return value.Wrap(c.handle.NewNull())
}

func (c *Context) NewUndefined() value.Value {
	return value.Wrap(c.handle.NewUndefined())
}

func (c *Context) NewObject() value.Value {
	return value.Wrap(c.handle.NewObject())
}

func (c *Context) NewArray() value.Value {
	return value.Wrap(c.handle.NewArray())
}

func (c *Context) GetPropertyString(obj value.Value, prop string) (value.Value, error) {
	if err := c.guardValue(obj); err != nil {
		return value.Value{}, err
	}
	out, err := c.handle.GetPropertyString(obj.Handle(), prop)
	if err != nil {
		return value.Value{}, err
	}
	return value.Wrap(out), nil
}

func (c *Context) SetPropertyString(obj value.Value, prop string, v value.Value) error {
	if err := c.guardValue(obj); err != nil {
		return err
	}
	if err := c.guardValue(v); err != nil {
		return err
	}
	return c.handle.SetPropertyString(obj.Handle(), prop, v.Handle())
}

func (c *Context) GetPropertyUint32(obj value.Value, idx uint32) (value.Value, error) {
	if err := c.guardValue(obj); err != nil {
		return value.Value{}, err
	}
	out, err := c.handle.GetPropertyUint32(obj.Handle(), idx)
	if err != nil {
		return value.Value{}, err
	}
	return value.Wrap(out), nil
}

func (c *Context) SetPropertyUint32(obj value.Value, idx uint32, v value.Value) error {
	if err := c.guardValue(obj); err != nil {
		return err
	}
	if err := c.guardValue(v); err != nil {
		return err
	}
	return c.handle.SetPropertyUint32(obj.Handle(), idx, v.Handle())
}

type GoFunction func(ctx *Context, this value.Value, args []value.Value) (value.Value, error)

func (c *Context) RegisterFunction(name string, length int, fn GoFunction) error {
	if err := c.guard(); err != nil {
		return err
	}
	id := cgo.RegisterCallback(func(ch *cgo.ContextHandle, this cgo.ValueHandle, args []cgo.ValueHandle) (cgo.ValueHandle, error) {
		wargs := make([]value.Value, len(args))
		for i := range args {
			wargs[i] = value.Wrap(args[i])
		}
		out, err := fn(c, value.Wrap(this), wargs)
		if err != nil {
			return cgo.ValueHandle{}, err
		}
		return out.Handle(), nil
	})

	if err := c.handle.SetGlobalFunction(name, length, id); err != nil {
		cgo.UnregisterCallback(id)
		return err
	}

	c.mu.Lock()
	c.callbackIDs = append(c.callbackIDs, id)
	c.mu.Unlock()
	return nil
}

func (c *Context) Handle() *cgo.ContextHandle {
	if c == nil {
		return nil
	}
	return c.handle
}

func (c *Context) guard() error {
	if c == nil || c.handle == nil {
		return fmt.Errorf("quickjs: context is nil")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.closed || c.handle == nil || c.handle.Ptr() == nil {
		return fmt.Errorf("quickjs: context is closed")
	}
	if c.ownerGID != 0 && currentGID() != c.ownerGID {
		return fmt.Errorf("quickjs: context is bound to one goroutine")
	}
	return nil
}

func (c *Context) guardValue(v value.Value) error {
	if err := c.guard(); err != nil {
		return err
	}
	h := v.Handle()
	if h.IsReleased() {
		return fmt.Errorf("quickjs: value already released")
	}
	if !h.BelongsTo(c.handle) {
		return fmt.Errorf("quickjs: value belongs to a different context")
	}
	return nil
}

func currentGID() uint64 {
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	line := strings.TrimPrefix(string(buf[:n]), "goroutine ")
	idx := strings.IndexByte(line, ' ')
	if idx < 0 {
		return 0
	}
	id, err := strconv.ParseUint(line[:idx], 10, 64)
	if err != nil {
		return 0
	}
	return id
}
