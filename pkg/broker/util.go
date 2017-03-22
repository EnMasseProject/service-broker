package broker

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func RunCommand(bin string, args ...string) {
	cmd := exec.Command(bin, args...) //.CombinedOutput()
	cmd.Stdin = os.Stdin
	cmd.Stdin = os.Stdout
	cmd.Stdin = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}

	return
}

func ProjectRoot() string {
	gopath := os.Getenv("GOPATH")
	rootPath := path.Join(gopath, "src", "github.com", "fusor",
		"ansible-service-broker")
	return rootPath
}

func SanitizePlanName(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	return name
}
