package common

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// url is the adress of the remote repository
// p is the path of the local repository
func Clone(url, p string) (err error) {
	// 判断是否是url
	if !strings.HasPrefix(url, "https://") || !strings.HasSuffix(url, ".git") {
		return errors.New("repository url is not correct")
	}

	var c *exec.Cmd
	if p == "" {
		c = exec.Command("git", "clone", url)
	} else {
		c = exec.Command("git", "clone", url, p)
	}

	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	err = c.Run()
	return
}

func Pull() (err error) {
	//git pull
	cmd := exec.Command("git", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
