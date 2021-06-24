package main
//
//import (
//	"fmt"
//	griffon_lib "git.dar.tech/griffon-open/griffon-lib"
//	ps "github.com/bhendo/go-powershell"
//	"github.com/bhendo/go-powershell/backend"
//	"github.com/joho/godotenv"
//	"os"
//	"strings"
//)
//
//var (
//	configPath        = ".env"
//	bucket            = ""
//	importType        = ""
//	adminUsername     = ""
//	adminPassword     = ""
//	adminClientId     = ""
//	adminClientSecret = ""
//	adminGrantType    = ""
//)
//
//type WindowUser struct {
//	RID      string `json:"rid"`
//	Username string `json:"username"`
//	HashNTLM string `json:"hash_ntlm"`
//}
//
//func setValues() {
//	if configPath != "" {
//		godotenv.Overload(configPath)
//	}
//	bucket = strings.TrimSpace(os.Getenv("BUCKET"))
//	adminUsername = strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))
//	adminPassword = strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
//	adminClientId = strings.TrimSpace(os.Getenv("ADMIN_CLIENT_ID"))
//	adminClientSecret = strings.TrimSpace(os.Getenv("ADMIN_CLIENT_SECRET"))
//	adminGrantType = strings.TrimSpace(os.Getenv("ADMIN_GRANT_TYPE"))
//	importType = strings.TrimSpace(os.Getenv("IMPORT_TYPE"))
//}
//
//func main() {
//
//}
