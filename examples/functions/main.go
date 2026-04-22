package main

import (
	"fmt"
	"log"
	"time"

	"github.com/znt-sh/quickjs-bindings/quickjs"
)

func main() {
	currentTime := time.Now()
	rt, err := quickjs.NewRuntime()
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	ctx, err := rt.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Close()

	err = ctx.RegisterFunction("add", 2, func(ctx *quickjs.Context, _ quickjs.Value, args []quickjs.Value) (quickjs.Value, error) {
		if err := ctx.AssertArgsAtLeast(args, 2); err != nil {
			return quickjs.Value{}, err
		}
		a := ctx.MustFloat(args[0])
		b := ctx.MustFloat(args[1])
		return ctx.NewFloat(a + b), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := ctx.Eval("add(20, 22)")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Free()

	number, err := ctx.ToFloat64(result)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(currentTime)

	fmt.Printf("add result: %.0f\n", number)
	fmt.Printf("elapsed: %v\n", elapsed)
}
