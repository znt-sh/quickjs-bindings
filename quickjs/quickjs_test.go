package quickjs_test

import (
	"strings"
	"testing"

	"github.com/znt-sh/quickjs-bindings/quickjs"
)

func TestEvalPrimitive(t *testing.T) {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext failed: %v", err)
	}
	defer ctx.Close()

	value, err := ctx.Eval("40 + 2")
	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}
	defer value.Free()

	out, err := ctx.ToFloat64(value)
	if err != nil {
		t.Fatalf("ToFloat64 failed: %v", err)
	}

	if out != 42 {
		t.Fatalf("expected 42, got %v", out)
	}
}

func TestExceptionIncludesStack(t *testing.T) {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext failed: %v", err)
	}
	defer ctx.Close()

	_, err = ctx.Eval("function fail(){ throw new Error('boom') }; fail();")
	if err == nil {
		t.Fatal("expected error")
	}

	qerr, ok := err.(*quickjs.Error)
	if !ok {
		t.Fatalf("expected *quickjs.Error, got %T", err)
	}
	if !strings.Contains(qerr.Message, "boom") {
		t.Fatalf("expected message to contain boom, got %q", qerr.Message)
	}
	if qerr.Stack == "" {
		t.Fatal("expected JS stack to be present")
	}
}

func TestTypedObjectMappingHelpers(t *testing.T) {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext failed: %v", err)
	}
	defer ctx.Close()

	obj, err := ctx.NewObjectFromMap(map[string]any{
		"name": "quickjs",
		"ok":   true,
		"n":    7,
	})
	if err != nil {
		t.Fatalf("NewObjectFromMap failed: %v", err)
	}
	defer obj.Free()

	if err := ctx.SetAny(obj, "pi", 3.14); err != nil {
		t.Fatalf("SetAny failed: %v", err)
	}

	name, err := ctx.GetString(obj, "name")
	if err != nil || name != "quickjs" {
		t.Fatalf("GetString mismatch: %q, err=%v", name, err)
	}

	ok, err := ctx.GetBool(obj, "ok")
	if err != nil || !ok {
		t.Fatalf("GetBool mismatch: %v, err=%v", ok, err)
	}

	n, err := ctx.GetFloat(obj, "n")
	if err != nil || n != 7 {
		t.Fatalf("GetFloat mismatch: %v, err=%v", n, err)
	}

	mapped, err := ctx.ObjectToMap(obj, "name", "ok", "n", "pi")
	if err != nil {
		t.Fatalf("ObjectToMap failed: %v", err)
	}
	if mapped["name"] != "quickjs" {
		t.Fatalf("unexpected name: %#v", mapped["name"])
	}
}

func TestContextBoundToOneGoroutine(t *testing.T) {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext failed: %v", err)
	}
	defer ctx.Close()

	ch := make(chan error, 1)
	go func() {
		_, goroutineErr := ctx.Eval("1+1")
		ch <- goroutineErr
	}()

	goroutineErr := <-ch
	if goroutineErr == nil {
		t.Fatal("expected guard error for cross-goroutine context usage")
	}
	if !strings.Contains(goroutineErr.Error(), "bound to one goroutine") {
		t.Fatalf("unexpected error: %v", goroutineErr)
	}
}

func TestRegisterFunction(t *testing.T) {
	rt, err := quickjs.NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime failed: %v", err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		t.Fatalf("NewContext failed: %v", err)
	}
	defer ctx.Close()

	if err := ctx.RegisterFunction("mul", 2, func(ctx *quickjs.Context, _ quickjs.Value, args []quickjs.Value) (quickjs.Value, error) {
		if err := ctx.AssertArgsAtLeast(args, 2); err != nil {
			return quickjs.Value{}, err
		}
		a, err := ctx.ToFloat64(args[0])
		if err != nil {
			return quickjs.Value{}, err
		}
		b, err := ctx.ToFloat64(args[1])
		if err != nil {
			return quickjs.Value{}, err
		}
		return ctx.NewFloat(a * b), nil
	}); err != nil {
		t.Fatalf("RegisterFunction failed: %v", err)
	}

	out, err := ctx.Eval("mul(6, 7)")
	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}
	defer out.Free()

	got, err := ctx.ToFloat64(out)
	if err != nil {
		t.Fatalf("ToFloat64 failed: %v", err)
	}

	if got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}
