package cmd

import (
	"bufio"
	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type fileStatus struct {
	file   string
	status git.StatusCode
}

// Print fileStatus print
func (fs fileStatus) Print(i int) {
	switch fs.status {
	case git.Untracked:
		// print serial number and file name and format alignment
		_, err := tm.Printf("%d. %s\n", i, tm.Color(fs.file, tm.WHITE))
		cobra.CheckErr(err)
	case git.Added:
		_, err := tm.Printf("%d. %s\n", i, tm.Color(fs.file, tm.GREEN))
		cobra.CheckErr(err)
	case git.Modified:
		_, err := tm.Printf("%d. %s\n", i, tm.Color(fs.file, tm.YELLOW))
		cobra.CheckErr(err)
	case git.Deleted:
		_, err := tm.Printf("%d. %s\n", i, tm.Color(fs.file, tm.RED))
		cobra.CheckErr(err)
	default:
		_, err := tm.Printf("%d. %s\n", i, tm.Color(fs.file, tm.WHITE))
		cobra.CheckErr(err)
	}
}

func commit(fileStatusList []fileStatus) {
	for i, fs := range fileStatusList {
		fs.Print(i + 1)
	}

	tm.Flush()
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "commit your code",
	//DisableAutoGenTag:  true,
	//DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := tm.Println("The following is the file status list:")
		cobra.CheckErr(err)
		commit(fileStatusList)
		tm.Flush()
		var fileIndex string
		prompt := &survey.Input{
			Message: "Please input the serial number of the file you want to commit:",
		}
		err = survey.AskOne(prompt, &fileIndex)
		cobra.CheckErr(err)
		// split by ,
		for _, index := range strings.Split(fileIndex, ",") {
			i, err := strconv.Atoi(strings.TrimSpace(index))
			cobra.CheckErr(err)
			_, err = tm.Println(fileStatusList[i-1].file)
			cobra.CheckErr(err)
			tm.Flush()

			// add file to staging area
			_, err = workTree.Add(fileStatusList[i-1].file)
		}
		// commit input email, name and message
		//var email, name, message string
		//prompt = &survey.Input{
		//	Message: "Please input your email:",
		//}
		//err = survey.AskOne(prompt, &email)
		//cobra.CheckErr(err)
		//prompt = &survey.Input{
		//	Message: "Please input your name:",
		//}
		//err = survey.AskOne(prompt, &name)
		//cobra.CheckErr(err)

		var message string
		prompt = &survey.Input{
			Message: "Please input your commit message:",
		}
		err = survey.AskOne(prompt, &message)
		cobra.CheckErr(err)
		_, err = workTree.Commit(message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  globalName,
				Email: globalMail,
				When:  time.Now(),
			},
		})
		cobra.CheckErr(err)
	},
}

// git the status of the current repository
func gitStatus() []fileStatus {
	var err error
	path, err := os.Getwd()
	cobra.CheckErr(err)

	r, err := git.PlainOpen(path)
	cobra.CheckErr(err)

	// getCmd the worktree
	workTree, err = r.Worktree()
	cobra.CheckErr(err)

	status, err := workTree.Status()
	cobra.CheckErr(err)

	var statusList []fileStatus
	for file, s := range status {
		var fs fileStatus
		if s.Staging == git.Deleted || s.Worktree == git.Deleted {
			fs = fileStatus{file: file, status: git.Deleted}
		} else if s.Staging == git.Added || s.Worktree == git.Added {
			fs = fileStatus{file: file, status: git.Added}
		} else if s.Staging == git.Modified || s.Worktree == git.Modified {
			fs = fileStatus{file: file, status: git.Modified}
		} else if s.Staging == git.Untracked || s.Worktree == git.Untracked {
			fs = fileStatus{file: file, status: git.Untracked}
		}
		statusList = append(statusList, fs)
	}

	return statusList
}

// 获取git的用户名和邮箱
func getGitConfig() (string, string) {
	var n, e string
	// 查看.git config文件是否存在
	_, err := os.Stat(".git/config")
	if err == os.ErrNotExist {
		goto git
	} else if err != nil {
		cobra.CheckErr(err)
	} else {
		// 读取文件内容
		file, err := os.Open(".git/config")
		cobra.CheckErr(err)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				cobra.CheckErr(err)
			}
		}(file)

		var name, email string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "name") {
				name = line[strings.Index(line, "=")+1:]
			}
			if strings.Contains(line, "email") {
				email = line[strings.Index(line, "=")+1:]
			}
		}
		cobra.CheckErr(scanner.Err())

		if name != "" && email != "" {
			return name, email
		} else {
			goto git
		}
	}

git:
	var email, name []byte
	// 获取系统的用户名
	name, err = exec.Command("git", "config", "--global", "user.name").Output()
	cobra.CheckErr(err)
	// 获取系统的邮箱
	email, err = exec.Command("git", "config", "--global", "user.email").Output()
	cobra.CheckErr(err)

	if e == "" && n == "" {
		return strings.TrimSpace(string(name)), strings.TrimSpace(string(email))
	} else if e == "" {
		return n, strings.TrimSpace(string(email))
	} else {
		return strings.TrimSpace(string(name)), e
	}
}
