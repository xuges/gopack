# GoPack

**Pack your program and it's external dependents to single-file!**



## Overview

Go projects will lose deployed-as-single-file when using CGO libraries, this `gopack` library can get it back!

You can also use this library to deploy C/C++, Java, Python and other programs as single-file!

Go version: go1.16 later.



## Usage

Now you have the `libvalue.c` file

```c
//libvalue.c
int getValue() { return 96; }
```

It will generate the `libvalue.so`:

```bash
gcc -shared -fPIC -o libvalue.so libvalue.c
```

And you want call `getValue` in Go:

```go
//main.go
package main

//#cgo LDFLAG: -L. -lvalue
//extern int getValue();
import "C"

func main() {
	value := int(C.getValue())
	println(value)
}
```

Compile Go program:

```bash
go build -o main.run main.go
```

The Go program will dependent the `libvalue.so`:

```bash
./main.run
#./main.run: error while loading shared libraries: libvalue.so: cannot open shared object file: No such file or directory
ldd main.run
#libvalue.so => not found
```

You need to provide both  `main.run`  and  `libvalue.so` , and instructions for settings the `LD_LIBRARY_PATH` environment variable :(

But  `gopack`  solves for now!

Look the `single.go`:

```go
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
	gopack.SetWorkerDir("extract/bin")
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
```

Compile the `single.go` and run it:

```bash
go build -o single.run single.go
./single.run
#program exit code: 96
```

The `single.run` packed `main.run` and `libvalue.so`, and then runs `main.run` after unpacking at runtime.

Now `main.run` can be deployed as a single file :)

```bash
ldd single.run
#not a dynamic executable
```

