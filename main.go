package main

import (
	"fmt"

	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
)

func main() {
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		panic(err)
	}
	defer shell.Exit()

	// ... and interact with it
	stdout, stderr, err := shell.Execute("reg save HKLM\\SYSTEM system.hiv")
	if err != nil {
		panic(err.Error())
	}
	stout, stderr, err := shell.Execute("reg save HKLM\\SAM sam.hiv")
	fmt.Println("error", stderr)
	fmt.Println(stdout)
}
