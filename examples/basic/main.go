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

	result, err := ctx.Eval("1 + 2 * 3")
	if err != nil {
		log.Fatal(err)
	}
	defer result.Free()

	number, err := ctx.ToFloat64(result)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(currentTime)
	fmt.Printf("result: %.0f\n", number)
	fmt.Printf("elapsed: %v\n", elapsed)
}
