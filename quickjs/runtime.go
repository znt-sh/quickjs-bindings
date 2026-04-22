package quickjs

import iruntime "github.com/znt-sh/quickjs-bindings/internal/runtime"

type Runtime struct {
	inner *iruntime.Runtime
}

func NewRuntime(opts ...RuntimeOption) (*Runtime, error) {
	rt, err := iruntime.New()
	if err != nil {
		return nil, wrapError(err)
	}
	out := &Runtime{inner: rt}
	for _, opt := range opts {
		if opt != nil {
			opt(out)
		}
	}
	return out, nil
}

func (r *Runtime) Close() {
	if r == nil || r.inner == nil {
		return
	}
	r.inner.Close()
}

func (r *Runtime) SetMemoryLimit(limitBytes uint64) {
	if r == nil || r.inner == nil {
		return
	}
	r.inner.SetMemoryLimit(limitBytes)
}

func (r *Runtime) NewContext() (*Context, error) {
	return NewContext(r)
}
