package run

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	OutputBaseDirectory = getUserHome()
	OutputDirectory     = path.Join(OutputBaseDirectory, InstallBaseDirectory)
)

const (
	InstallBaseDirectory string      = "spike/"
	FilePermission       fs.FileMode = 00775
)

func getUserHome() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	return homedir
}

// Create a new task based on the data input via the CLI
func Create(name string, description string, status string) int {
    return 0
}

// Git command wrapper for `git add`
func Add(files []string) {
	out, err := exec.Command("git", "-C", OutputDirectory, "add", strings.Join(files, " ")).Output()
	if err != nil {
		fmt.Println("`git add` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git init`
func Init() {
	if _, err := os.Stat(OutputDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(OutputDirectory, FilePermission)
		if err != nil {
			fmt.Println("Failed to create directories")
			fmt.Println(err)
		}
	}

	out, err := exec.Command("git", "-C", OutputDirectory, "init").Output()
	if err != nil {
		fmt.Println("`git init` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git push`
func Push() {
	out, err := exec.Command("git", "-C", OutputDirectory, "push").Output()
	if err != nil {
		fmt.Println("`git push` exited abnormally")
		fmt.Println(err)
	}

	output := string(out)

	fmt.Print(output)
}

// Git command wrapper for `git pull`
func Pull() {
	out, err := exec.Command("git", "-C", OutputDirectory, "pull").Output()
	if err != nil {
		fmt.Println("`git push` exited abnormally")
		fmt.Println(err)
	}

	output := string(out[:])

	fmt.Print(output)
}
