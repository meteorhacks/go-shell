package main

import (
	"github.com/meteorhacks/go-shell/gosh"
	"github.com/robertkrimen/otto"
)

var (
	sh gosh.Shell
)

func main() {
	sh = gosh.New()
	sh.SetVar("name", "John Doe")
	sh.SetVar("hello", CmdHello)
	sh.Start()
}

func CmdHello(call otto.FunctionCall) (val otto.Value) {
	name, err := call.Argument(0).ToString()
	if err != nil {
		sh.PrintError(err)
		return gosh.Undefined
	}

	return sh.Value("Hello " + name)
}
