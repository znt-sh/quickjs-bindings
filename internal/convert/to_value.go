package convert

import (
	"fmt"

	"github.com/znt-sh/quickjs-bindings/internal/context"
	"github.com/znt-sh/quickjs-bindings/internal/value"
)

func ToJS(ctx *context.Context, in any) (value.Value, error) {
	switch v := in.(type) {
	case nil:
		return ctx.NewNull(), nil
	case bool:
		return ctx.NewBool(v), nil
	case int:
		return ctx.NewInt32(int32(v)), nil
	case int32:
		return ctx.NewInt32(v), nil
	case int64:
		return ctx.NewFloat64(float64(v)), nil
	case float32:
		return ctx.NewFloat64(float64(v)), nil
	case float64:
		return ctx.NewFloat64(v), nil
	case string:
		return ctx.NewString(v), nil
	case []any:
		arr := ctx.NewArray()
		for i, item := range v {
			itemValue, err := ToJS(ctx, item)
			if err != nil {
				arr.Free()
				return value.Value{}, err
			}
			if err := ctx.SetPropertyUint32(arr, uint32(i), itemValue); err != nil {
				itemValue.Free()
				arr.Free()
				return value.Value{}, err
			}
			itemValue.Free()
		}
		return arr, nil
	case map[string]any:
		obj := ctx.NewObject()
		for key, item := range v {
			itemValue, err := ToJS(ctx, item)
			if err != nil {
				obj.Free()
				return value.Value{}, err
			}
			if err := ctx.SetPropertyString(obj, key, itemValue); err != nil {
				itemValue.Free()
				obj.Free()
				return value.Value{}, err
			}
			itemValue.Free()
		}
		return obj, nil
	default:
		return value.Value{}, fmt.Errorf("quickjs: unsupported Go value type %T", in)
	}
}
