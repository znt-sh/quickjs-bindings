package convert

import (
	"fmt"

	"github.com/znt-sh/quickjs-bindings/internal/context"
	"github.com/znt-sh/quickjs-bindings/internal/value"
)

func FromJS(ctx *context.Context, in value.Value) (any, error) {
	if in.IsNull() || in.IsUndefined() {
		return nil, nil
	}
	if in.IsBool() {
		return ctx.ToBool(in)
	}
	if in.IsNumber() {
		return ctx.ToFloat64(in)
	}
	if in.IsString() {
		return ctx.ToString(in)
	}
	if in.IsArray() {
		lengthValue, err := ctx.GetPropertyString(in, "length")
		if err != nil {
			return nil, err
		}
		defer lengthValue.Free()

		length, err := ctx.ToInt32(lengthValue)
		if err != nil {
			return nil, err
		}

		out := make([]any, length)
		for i := int32(0); i < length; i++ {
			item, err := ctx.GetPropertyUint32(in, uint32(i))
			if err != nil {
				return nil, err
			}
			converted, err := FromJS(ctx, item)
			item.Free()
			if err != nil {
				return nil, err
			}
			out[i] = converted
		}
		return out, nil
	}
	if in.IsObject() {
		return nil, fmt.Errorf("quickjs: generic object conversion is not automatic; use explicit property access")
	}
	return nil, fmt.Errorf("quickjs: unsupported JS value")
}
