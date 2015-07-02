package gosh

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

var (
	ErrReadErr = errors.New("read error")
	Undefined  = otto.Value{}
)

type Shell interface {
	SetVar(name string, val interface{})
	Value(val interface{}) (out otto.Value)
	Prompt() (prm string)
	SetPrompt(prm string)
	PrintPrompt()
	PrintValue(val interface{})
	PrintError(err error)
	Start()
}

type shell struct {
	vm *otto.Otto

	prm string
	sin *bufio.Reader

	pprm *color.Color
	pval *color.Color
	perr *color.Color
}

func New() (sh Shell) {
	s := shell{
		vm:   otto.New(),
		prm:  "> ",
		sin:  bufio.NewReader(os.Stdin),
		pprm: color.New(color.FgWhite, color.Bold),
		pval: color.New(color.FgYellow),
		perr: color.New(color.FgRed),
	}

	return &s
}

func (sh *shell) SetVar(name string, val interface{}) {
	sh.vm.Set(name, val)
}

func (sh *shell) Value(val interface{}) (out otto.Value) {
	out, err := sh.vm.ToValue(val)
	if err != nil {
		sh.PrintError(err)
		return Undefined
	}

	return out
}

func (sh *shell) Prompt() (prm string) {
	return sh.prm
}

func (sh *shell) SetPrompt(prm string) {
	sh.prm = prm
}

//   PRINTERS
// ------------

func (sh *shell) PrintPrompt() {
	sh.pprm.Printf(sh.prm)
}

func (sh *shell) PrintValue(val interface{}) {
	sh.pval.Printf("%+v\n", val)
}

func (sh *shell) PrintError(err error) {
	sh.perr.Printf("ERR: %s\n", err.Error())
}

//   MAIN LOOP
// -------------

func (sh *shell) Start() {
	for {
		sh.PrintPrompt()

		str, err := sh.sin.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			sh.PrintError(ErrReadErr)
			continue
		} else if strings.TrimSpace(str) == "" {
			continue
		}

		val, err := sh.vm.Run(str)
		if err != nil {
			sh.PrintError(err)
			continue
		} else if val.IsUndefined() {
			continue
		}

		sh.PrintValue(val)
	}
}
