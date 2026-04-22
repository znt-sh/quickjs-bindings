package quickjs

type RuntimeOption func(*Runtime)

func WithMemoryLimit(limitBytes uint64) RuntimeOption {
	return func(r *Runtime) {
		if r == nil || r.inner == nil {
			return
		}
		r.inner.SetMemoryLimit(limitBytes)
	}
}
