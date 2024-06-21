package main

import (
	_ "embed"
	"github.com/xuges/gopack"
)

//go:embed libvalue.so
var lib []byte

//go:embed main.run
var run []byte

func main() {
	gopack.SetUnpackPath("unpacked")
	gopack.AddDependency("lib/libvalue.so", lib)
	gopack.AddExecutable("bin/main.run", run)
	gopack.SetWorkerDir("unpacked/bin")
	gopack.AddEnv("LD_LIBRARY_PATH=../lib")
	err := gopack.Unpack()
	if err != nil {
		panic(err)
	}

	code, err := gopack.Run()
	if err != nil {
		panic(err)
	}

	println("program exit code:", code)
}
