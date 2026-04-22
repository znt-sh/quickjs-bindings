package cgo

/*
#cgo CFLAGS: -I${SRCDIR}/../../third_party/quickjs
#cgo windows CFLAGS: -D_WIN32
#cgo linux LDFLAGS: -lm

#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

const EvalTypeGlobal = int(C.JS_EVAL_TYPE_GLOBAL)

type JSError struct {
	Message string
	Stack   string
}

func (e *JSError) Error() string {
	if e == nil {
		return "quickjs: unknown exception"
	}
	if e.Stack == "" {
		return e.Message
	}
	return e.Message + "\n" + e.Stack
}

type RuntimeHandle struct {
	ptr *C.JSRuntime
}

func NewRuntime() (*RuntimeHandle, error) {
	ptr := C.qjs_new_runtime()
	if ptr == nil {
		return nil, errors.New("quickjs: failed to create runtime")
	}
	return &RuntimeHandle{ptr: ptr}, nil
}

func (r *RuntimeHandle) Ptr() *C.JSRuntime {
	if r == nil {
		return nil
	}
	return r.ptr
}

func (r *RuntimeHandle) SetMemoryLimit(limitBytes uint64) {
	if r == nil || r.ptr == nil {
		return
	}
	C.qjs_set_memory_limit(r.ptr, C.size_t(limitBytes))
}

func (r *RuntimeHandle) Close() {
	if r == nil || r.ptr == nil {
		return
	}
	C.qjs_free_runtime(r.ptr)
	r.ptr = nil
}

type ContextHandle struct {
	ptr *C.JSContext
}

func NewContext(rt *RuntimeHandle) (*ContextHandle, error) {
	if rt == nil || rt.ptr == nil {
		return nil, errors.New("quickjs: runtime is nil")
	}
	ptr := C.qjs_new_context(rt.ptr)
	if ptr == nil {
		return nil, errors.New("quickjs: failed to create context")
	}
	h := &ContextHandle{ptr: ptr}
	registerContext(h)
	return h, nil
}

func (c *ContextHandle) Ptr() *C.JSContext {
	if c == nil {
		return nil
	}
	return c.ptr
}

func (c *ContextHandle) Close() {
	if c == nil || c.ptr == nil {
		return
	}
	unregisterContext(c.ptr)
	C.qjs_free_context(c.ptr)
	c.ptr = nil
}

type ValueHandle struct {
	ctx      *ContextHandle
	raw      C.JSValue
	owned    bool
	released bool
}

func (v ValueHandle) Raw() C.JSValue {
	return v.raw
}

func (v ValueHandle) Dup() ValueHandle {
	if v.released || v.ctx == nil || v.ctx.ptr == nil {
		return ValueHandle{}
	}
	return ValueHandle{
		ctx:   v.ctx,
		raw:   C.qjs_dup_value(v.ctx.ptr, C.JSValueConst(v.raw)),
		owned: true,
	}
}

func (v *ValueHandle) Free() {
	if v == nil || v.released || !v.owned || v.ctx == nil || v.ctx.ptr == nil {
		return
	}
	C.qjs_free_value(v.ctx.ptr, v.raw)
	v.owned = false
	v.released = true
}

func (v ValueHandle) IsReleased() bool {
	return v.released
}

func (v ValueHandle) BelongsTo(ctx *ContextHandle) bool {
	if v.ctx == nil || v.ctx.ptr == nil || ctx == nil || ctx.ptr == nil {
		return false
	}
	return v.ctx.ptr == ctx.ptr
}

func (v ValueHandle) IsException() bool {
	return int(C.qjs_is_exception(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsNull() bool {
	return int(C.qjs_is_null(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsUndefined() bool {
	return int(C.qjs_is_undefined(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsBool() bool {
	return int(C.qjs_is_bool(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsNumber() bool {
	return int(C.qjs_is_number(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsString() bool {
	return int(C.qjs_is_string(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsObject() bool {
	return int(C.qjs_is_object(C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsArray() bool {
	if v.ctx == nil || v.ctx.ptr == nil {
		return false
	}
	return int(C.qjs_is_array(v.ctx.ptr, C.JSValueConst(v.raw))) != 0
}

func (v ValueHandle) IsFunction() bool {
	if v.ctx == nil || v.ctx.ptr == nil {
		return false
	}
	return int(C.qjs_is_function(v.ctx.ptr, C.JSValueConst(v.raw))) != 0
}

func (c *ContextHandle) Eval(code, filename string, flags int) (ValueHandle, error) {
	if c == nil || c.ptr == nil {
		return ValueHandle{}, errors.New("quickjs: context is nil")
	}

	cCode := C.CString(code)
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cCode))
	defer C.free(unsafe.Pointer(cFilename))

	v := C.qjs_eval(c.ptr, cCode, C.size_t(len(code)), cFilename, C.int(flags))
	out := ValueHandle{ctx: c, raw: v, owned: true}
	if out.IsException() {
		out.Free()
		return ValueHandle{}, c.Exception()
	}
	return out, nil
}

func (c *ContextHandle) Exception() error {
	if c == nil || c.ptr == nil {
		return errors.New("quickjs: unknown exception")
	}

	ex := ValueHandle{ctx: c, raw: C.qjs_get_exception(c.ptr), owned: true}
	defer ex.Free()

	message := "quickjs: exception"
	if msg, err := c.ToString(ex); err == nil && msg != "" {
		message = msg
	}

	stack := ""
	if stackV, err := c.GetPropertyString(ex, "stack"); err == nil {
		if !stackV.IsUndefined() && !stackV.IsNull() {
			if text, convErr := c.ToString(stackV); convErr == nil {
				stack = text
			}
		}
		stackV.Free()
	}

	return &JSError{Message: message, Stack: stack}
}

func (c *ContextHandle) ToString(v ValueHandle) (string, error) {
	if c == nil || c.ptr == nil {
		return "", errors.New("quickjs: context is nil")
	}
	cstr := C.qjs_to_cstring(c.ptr, C.JSValueConst(v.raw))
	if cstr == nil {
		return "", c.Exception()
	}
	defer C.qjs_free_cstring(c.ptr, cstr)
	return C.GoString(cstr), nil
}

func (c *ContextHandle) ToInt32(v ValueHandle) (int32, error) {
	if c == nil || c.ptr == nil {
		return 0, errors.New("quickjs: context is nil")
	}
	var out C.int32_t
	if int(C.qjs_to_int32(c.ptr, &out, C.JSValueConst(v.raw))) < 0 {
		return 0, c.Exception()
	}
	return int32(out), nil
}

func (c *ContextHandle) ToFloat64(v ValueHandle) (float64, error) {
	if c == nil || c.ptr == nil {
		return 0, errors.New("quickjs: context is nil")
	}
	var out C.double
	if int(C.qjs_to_float64(c.ptr, &out, C.JSValueConst(v.raw))) < 0 {
		return 0, c.Exception()
	}
	return float64(out), nil
}

func (c *ContextHandle) ToBool(v ValueHandle) (bool, error) {
	if c == nil || c.ptr == nil {
		return false, errors.New("quickjs: context is nil")
	}
	rc := int(C.qjs_to_bool(c.ptr, C.JSValueConst(v.raw)))
	if rc < 0 {
		return false, c.Exception()
	}
	return rc != 0, nil
}

func (c *ContextHandle) NewInt32(v int32) ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_int32(c.ptr, C.int32_t(v)), owned: true}
}

func (c *ContextHandle) NewFloat64(v float64) ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_float64(c.ptr, C.double(v)), owned: true}
}

func (c *ContextHandle) NewBool(v bool) ValueHandle {
	iv := 0
	if v {
		iv = 1
	}
	return ValueHandle{ctx: c, raw: C.qjs_new_bool(c.ptr, C.int(iv)), owned: true}
}

func (c *ContextHandle) NewString(v string) ValueHandle {
	cs := C.CString(v)
	defer C.free(unsafe.Pointer(cs))
	return ValueHandle{ctx: c, raw: C.qjs_new_string(c.ptr, cs), owned: true}
}

func (c *ContextHandle) NewNull() ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_null(), owned: true}
}

func (c *ContextHandle) NewUndefined() ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_undefined(), owned: true}
}

func (c *ContextHandle) NewObject() ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_object(c.ptr), owned: true}
}

func (c *ContextHandle) NewArray() ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_new_array(c.ptr), owned: true}
}

func (c *ContextHandle) GetPropertyString(obj ValueHandle, prop string) (ValueHandle, error) {
	cs := C.CString(prop)
	defer C.free(unsafe.Pointer(cs))
	out := ValueHandle{ctx: c, raw: C.qjs_get_property_str(c.ptr, C.JSValueConst(obj.raw), cs), owned: true}
	if out.IsException() {
		out.Free()
		return ValueHandle{}, c.Exception()
	}
	return out, nil
}

func (c *ContextHandle) SetPropertyString(obj ValueHandle, prop string, value ValueHandle) error {
	cs := C.CString(prop)
	defer C.free(unsafe.Pointer(cs))
	if int(C.qjs_set_property_str_dup(c.ptr, obj.raw, cs, C.JSValueConst(value.raw))) < 0 {
		return c.Exception()
	}
	return nil
}

func (c *ContextHandle) GetPropertyUint32(obj ValueHandle, idx uint32) (ValueHandle, error) {
	out := ValueHandle{ctx: c, raw: C.qjs_get_property_uint32(c.ptr, C.JSValueConst(obj.raw), C.uint32_t(idx)), owned: true}
	if out.IsException() {
		out.Free()
		return ValueHandle{}, c.Exception()
	}
	return out, nil
}

func (c *ContextHandle) SetPropertyUint32(obj ValueHandle, idx uint32, value ValueHandle) error {
	if int(C.qjs_set_property_uint32_dup(c.ptr, obj.raw, C.uint32_t(idx), C.JSValueConst(value.raw))) < 0 {
		return c.Exception()
	}
	return nil
}

func (c *ContextHandle) GetGlobalObject() ValueHandle {
	return ValueHandle{ctx: c, raw: C.qjs_get_global_object(c.ptr), owned: true}
}

type GoFunction func(ctx *ContextHandle, this ValueHandle, args []ValueHandle) (ValueHandle, error)

var (
	callbackSeq int32 = 1

	callbacksMu sync.RWMutex
	callbacks   = map[int32]GoFunction{}

	contextsMu sync.RWMutex
	contexts   = map[uintptr]*ContextHandle{}
)

func RegisterCallback(fn GoFunction) int32 {
	id := atomic.AddInt32(&callbackSeq, 1)
	callbacksMu.Lock()
	callbacks[id] = fn
	callbacksMu.Unlock()
	return id
}

func UnregisterCallback(id int32) {
	callbacksMu.Lock()
	delete(callbacks, id)
	callbacksMu.Unlock()
}

func (c *ContextHandle) SetGlobalFunction(name string, length int, callbackID int32) error {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	if int(C.qjs_set_global_function(c.ptr, cs, C.int(length), C.int(callbackID))) < 0 {
		return c.Exception()
	}
	return nil
}

func registerContext(ctx *ContextHandle) {
	contextsMu.Lock()
	contexts[uintptr(unsafe.Pointer(ctx.ptr))] = ctx
	contextsMu.Unlock()
}

func unregisterContext(ptr *C.JSContext) {
	contextsMu.Lock()
	delete(contexts, uintptr(unsafe.Pointer(ptr)))
	contextsMu.Unlock()
}

func findContext(ptr *C.JSContext) *ContextHandle {
	contextsMu.RLock()
	ctx := contexts[uintptr(unsafe.Pointer(ptr))]
	contextsMu.RUnlock()
	if ctx != nil {
		return ctx
	}
	return &ContextHandle{ptr: ptr}
}

func callbackByID(id int32) GoFunction {
	callbacksMu.RLock()
	fn := callbacks[id]
	callbacksMu.RUnlock()
	return fn
}

//export goQuickjsInvoke
func goQuickjsInvoke(cctx *C.JSContext, thisVal C.JSValueConst, argc C.int, argv *C.JSValueConst, magic C.int) C.JSValue {
	fn := callbackByID(int32(magic))
	if fn == nil {
		msg := C.CString("quickjs: callback not found")
		defer C.free(unsafe.Pointer(msg))
		return C.qjs_throw_type_error(cctx, msg)
	}

	ctx := findContext(cctx)
	args := make([]ValueHandle, int(argc))
	if argc > 0 {
		argSlice := unsafe.Slice(argv, int(argc))
		for i := range argSlice {
			args[i] = ValueHandle{ctx: ctx, raw: C.JSValue(argSlice[i]), owned: false}
		}
	}

	thisHandle := ValueHandle{ctx: ctx, raw: C.JSValue(thisVal), owned: false}
	out, err := fn(ctx, thisHandle, args)
	if err != nil {
		msg := C.CString(err.Error())
		defer C.free(unsafe.Pointer(msg))
		return C.qjs_throw_type_error(cctx, msg)
	}

	if out.ctx == nil {
		msg := C.CString("quickjs: callback returned value without context")
		defer C.free(unsafe.Pointer(msg))
		return C.qjs_throw_type_error(cctx, msg)
	}

	if out.owned {
		out.owned = false
		return out.raw
	}
	return C.qjs_dup_value(cctx, C.JSValueConst(out.raw))
}

func (c *ContextHandle) ThrowTypeError(message string) ValueHandle {
	cs := C.CString(message)
	defer C.free(unsafe.Pointer(cs))
	return ValueHandle{ctx: c, raw: C.qjs_throw_type_error(c.ptr, cs), owned: true}
}

func (c *ContextHandle) MustPrimitive(v any) (ValueHandle, error) {
	switch x := v.(type) {
	case nil:
		return c.NewNull(), nil
	case bool:
		return c.NewBool(x), nil
	case int:
		return c.NewInt32(int32(x)), nil
	case int32:
		return c.NewInt32(x), nil
	case int64:
		return c.NewFloat64(float64(x)), nil
	case float32:
		return c.NewFloat64(float64(x)), nil
	case float64:
		return c.NewFloat64(x), nil
	case string:
		return c.NewString(x), nil
	default:
		return ValueHandle{}, fmt.Errorf("quickjs: unsupported primitive type %T", v)
	}
}
