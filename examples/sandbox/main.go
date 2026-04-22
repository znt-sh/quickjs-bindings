package main

import (
	"fmt"
	"log"
	"time"

	"github.com/znt-sh/quickjs-bindings/quickjs"
)

func main() {
	currentTime := time.Now()

	rt, err := quickjs.NewRuntime(quickjs.WithMemoryLimit(32 * 1024 * 1024))
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Close()

	ctxA, err := rt.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctxA.Close()

	ctxB, err := rt.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer ctxB.Close()

	_, err = ctxA.Eval("globalThis.value = 123")
	if err != nil {
		log.Fatal(err)
	}

	valB, err := ctxB.Eval("typeof globalThis.value")
	if err != nil {
		log.Fatal(err)
	}
	defer valB.Free()

	kind, err := ctxB.ToString(valB)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(currentTime)

	fmt.Printf("ctxB sees value as: %s\n", kind)
	fmt.Printf("elapsed: %v\n", elapsed)
}
