package common

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"os/exec"
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// RemoveStringSlice 通过字符串切片2，去掉字符串切片1中的元素
func RemoveStringSlice(a []string, b []string) []string {
	// 将b装进map
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}
	// 遍历a，如果a中的元素在b中，就删除
	for i := 0; i < len(a); i++ {
		if _, ok := bMap[a[i]]; ok {
			a = append(a[:i], a[i+1:]...)
			i--
		}
	}
	return a
}

func Push() (err error) {
	//git push
	cmd := exec.Command("git", "push")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func ShowDiff(file, hash string, add bool) (err error) {
	//git push
	cmd := exec.Command("git", "diff", hash, file)
	if add {
		cmd = exec.Command("git", "diff", file)
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func ShowLog(head, pHash, file string) (err error) {
	//git push
	cmd := exec.Command("git", "diff", pHash, head, file)
	if pHash == "" {
		cmd = exec.Command("git", "diff", head, file)
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func Indent(t string) string {
	var output []string
	for _, line := range strings.Split(t, "\n") {
		if len(line) != 0 {
			line = "    " + line
		}

		output = append(output, line)
	}

	return strings.Join(output, "\n")
}

func GetFileChangeByCommit(repo *git.Repository, commitHash string) (fileChange []string, err error) {
	commit, err := repo.CommitObject(plumbing.NewHash(commitHash))
	if err != nil {
		return
	}
	fileState, err := commit.Stats()
	if err != nil {
		return
	}

	for _, file := range fileState {
		fileChange = append(fileChange, file.Name)
	}
	return
}
