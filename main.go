package main

import (
	"fmt"
	griffon_lib "git.dar.tech/griffon-open/griffon-lib"
	ps "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	configPath        = ".env"
	bucket            = ""
	importType        = ""
	adminUsername     = ""
	adminPassword     = ""
	adminClientId     = ""
	adminClientSecret = ""
	adminGrantType    = ""
)

type WindowUser struct {
	RID      string `json:"rid"`
	Username string `json:"username"`
	HashNTLM string `json:"hash_ntlm"`
}

func setValues() {
	if configPath != "" {
		godotenv.Overload(configPath)
	}
	bucket = strings.TrimSpace(os.Getenv("BUCKET"))
	adminUsername = strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
	adminPassword = strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	adminClientId = strings.TrimSpace(os.Getenv("ADMIN_CLIENT_ID"))
	adminClientSecret = strings.TrimSpace(os.Getenv("ADMIN_CLIENT_SECRET"))
	adminGrantType = strings.TrimSpace(os.Getenv("ADMIN_GRANT_TYPE"))
	importType = strings.TrimSpace(os.Getenv("IMPORT_TYPE"))
}

func work() {
	setValues()
	fmt.Println(bucket)
	fmt.Println(adminGrantType)
	fmt.Println(adminPassword)
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
			Email:    v.Username + "@gmail.com",
			Password: v.HashNTLM,
			Bucket:   bucket,
		})
	}
	response, err := service.CreateBunchWithPasswords(&griffon_lib.CreateBunchWithPasswordsCommand{
		Bucket:     bucket,
		ImportType: importType,
		Users:      users,
	})
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(response)
	if response.Status == 200 {
		usersUpdate := []*griffon_lib.UserUpdate{}
		for _, v := range users {
			userUpdate := &griffon_lib.UserUpdate{}
			searchedUser, err := service.SearchUser(&griffon_lib.SearchUserCommand{
				BucketId:  bucket,
				Parameter: v.Email,
			})
			if err != nil {
				panic(err)
				return
			}
			if len(searchedUser) != 0 {
				currentUser, err := service.GetUser(&griffon_lib.GetUserCommand{
					Bucket: searchedUser[0].Bucket,
					Id:     searchedUser[0].ID,
				})
				if err != nil {
					panic(err)
					return
				}
				userUpdate.ID = searchedUser[0].ID
				if v.FirstName != currentUser.FirstName {
					userUpdate.FirstName = &v.FirstName
				}
				if v.LastName != currentUser.LastName {
					userUpdate.LastName = &v.LastName
				}
				if v.Password != currentUser.Password {
					userUpdate.Password = &v.Password
				}
				usersUpdate = append(usersUpdate, userUpdate)
				fmt.Printf("FROM AD FirstName:%s,LastName:%s,Password:%s \n", v.FirstName, v.LastName, v.Password)
				fmt.Printf("FOR DB FirstName:%s,LastName:%s,Password:%s \n", userUpdate.FirstName, userUpdate.LastName, userUpdate.Password)
			}
		}
		cmd := &griffon_lib.UpdateUsersCommand{
			Users:  usersUpdate,
			Bucket: bucket,
		}
		res, err := service.UpdateUsers(cmd)
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("update users ", res)
	}
	os.Chdir("../../")
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

		time.Sleep(0 * time.Second)
		command <- "Pause"

		time.Sleep(10 * time.Second)
		command <- "Play"

		time.Sleep(0 * time.Second)
		command <- "Stop"

		wg.Wait()
	}
}
