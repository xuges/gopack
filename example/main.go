package main

//extern int getValue();
//#cgo LDFLAGS: -L. -lvalue
import "C"
import (
	"os"
)

func main() {
	os.Exit(int(C.getValue()))
}
