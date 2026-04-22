# quickjs-bindings

[![CI](https://github.com/znt-sh/quickjs-bindings/actions/workflows/ci.yml/badge.svg)](https://github.com/znt-sh/quickjs-bindings/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.22%2B-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Production-focused Go bindings for the QuickJS JavaScript engine.

This project provides a small, idiomatic, cgo-backed API that keeps QuickJS internals hidden behind a stable Go package.

## Highlights

- Embedded QuickJS source via Git submodule
- Runtime and context lifecycle management
- Eval and value conversion helpers
- Go function registration callable from JavaScript
- Structured JS errors with stack traces
- Safety guards for value ownership and context goroutine affinity

## Requirements

- Go 1.22+
- A C compiler (cgo is required)

Windows:
- Recommended: MSYS2 UCRT64 with gcc and g++ in PATH

Linux:
- build-essential (gcc, g++, make)

macOS:
- Xcode Command Line Tools (clang)

## Installation

Add the module:

go get github.com/znt-sh/quickjs-bindings

If you clone this repository directly, initialize submodules:

git submodule update --init --recursive

## Quick Start

Run tests:

go test ./...

Run examples:

go run ./examples/basic
go run ./examples/functions
go run ./examples/sandbox

## Minimal Example

package main

import (
	"fmt"
	"log"

	"github.com/znt-sh/quickjs-bindings/quickjs"
)

func main() {
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

	v, err := ctx.Eval("1 + 2")
	if err != nil {
		log.Fatal(err)
	}
	defer v.Free()

	out, err := ctx.ToFloat64(v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}

## API Overview

Core package:
- quickjs

Main entry points:
- quickjs.NewRuntime
- runtime.NewContext
- context.Eval
- context.RegisterFunction
- context.ToValue and context.FromValue

Typed helpers:
- context.NewObjectFromMap
- context.NewArrayFromSlice
- context.GetString, context.GetFloat, context.GetInt, context.GetBool, context.GetAny
- context.SetAny
- context.ObjectToMap

## Documentation

- docs/getting-started.md
- docs/architecture.md
- docs/api.md
- docs/memory-and-threading.md
- CONTRIBUTING.md
- SECURITY.md

## Stability and Scope

Current scope targets a pragmatic MVP plus safety hardening.

Not included yet:
- ES module loader
- worker threads
- async event loop integration
- bytecode caching and snapshots

## Contributing

Contributions are welcome. Please read CONTRIBUTING.md before opening a pull request.

## License

MIT. See LICENSE.
