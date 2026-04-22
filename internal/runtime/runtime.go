package runtime

import "github.com/znt-sh/quickjs-bindings/internal/cgo"

type Runtime struct {
	handle *cgo.RuntimeHandle
}

func New() (*Runtime, error) {
	h, err := cgo.NewRuntime()
	if err != nil {
		return nil, err
	}
	return &Runtime{handle: h}, nil
}

func (r *Runtime) Close() {
	if r == nil || r.handle == nil {
		return
	}
	r.handle.Close()
}

func (r *Runtime) SetMemoryLimit(limitBytes uint64) {
	if r == nil || r.handle == nil {
		return
	}
	r.handle.SetMemoryLimit(limitBytes)
}

func (r *Runtime) Handle() *cgo.RuntimeHandle {
	if r == nil {
		return nil
	}
	return r.handle
}
