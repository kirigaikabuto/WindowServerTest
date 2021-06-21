package main

import (
	"fmt"
	griffon_lib "git.dar.tech/griffon-open/griffon-lib"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"os"
	"strings"
)

var (
	bucket            = "1fe099cc-cd7d-4556-b90d-1a9822395eba"
	importType        = "windows-server-2012"
	adminUsername     = "tleugazy_erasil@gmail.com"
	adminPassword     = "i77GPf#%"
	adminClientId     = "griffon"
	adminClientSecret = "$2a$10$qC9dtMHqvgbA/Rn10UV49OY4Lp6yETBsNKPTAdp4mnQcVL/.bDbQS"
	adminGrantType    = "password"
)

type WindowUser struct {
	RID      string `json:"rid"`
	Username string `json:"username"`
	HashNTLM string `json:"hash_ntlm"`
}

func main() {
	windowsUsers := []WindowUser{}
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
	_, _, err = shell.Execute("Remove-Item system2.hiv")
	if err != nil {
		panic(err.Error())
	}
	_, _, err = shell.Execute("Remove-Item sam2.hiv")
	if err != nil {
		panic(err.Error())
	}
	for i, v := range parts {
		if strings.Contains(v, "RID") {
			ridParts := strings.Split(strings.Split(v, ": ")[1], " ")
			userParts := strings.Split(strings.Split(parts[i+1], ": ")[1], " ")
			ntlmString := ""
			if strings.Contains(parts[i+2], "NTLM") {
				ntlmString = strings.Split(parts[i+2], ": ")[1]
			}

			windowsUsers = append(windowsUsers, WindowUser{
				RID:      ridParts[0],
				Username: userParts[0],
				HashNTLM: ntlmString,
			})
		}
	}
	service, err := griffon_lib.NewGriffonConnect(&griffon_lib.GriffonConnectCommand{
		ClientId:     adminClientId,
		ClientSecret: adminClientSecret,
		Username:     adminUsername,
		Password:     adminPassword,
		GrantType:    adminGrantType,
	})
	if err != nil {
		panic(err)
		return
	}
	users := []griffon_lib.GriffonUser{}
	for _, v := range windowsUsers {
		users = append(users, griffon_lib.GriffonUser{
			Email:    v.Username,
			Password: v.HashNTLM,
			Bucket:   bucket,
		})
	}
	response, err := service.CreateBunchWithPasswords(&griffon_lib.CreateBunchWithPasswordsCommand{
		Bucket:     bucket,
		ImportType: importType,
		Users:      nil,
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(response)
}
