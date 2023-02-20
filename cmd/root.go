// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

var cmd = &cobra.Command{
	Use:   "ggit",
	Short: "A simple git tool",
	Run: func(cmd *cobra.Command, args []string) {
		// 创建选择问题
		question := &survey.Select{
			Message: "Select Your Action:",
			Options: []string{"Commit", "Commit And Push", "Push"},
		}

		// 提示用户选择选项
		var answer string
		if err := survey.AskOne(question, &answer); err != nil {
			fmt.Println("Failed to prompt user:", err)
			os.Exit(1)
		}

		// 根据用户选择的选项提示用户输入
		switch answer {
		case "Commit":
			_, err := tm.Println("The following is the file status list:")
			cobra.CheckErr(err)
			commit(fileStatusList)
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

				// add file to staging area
				_, err = workTree.Add(fileStatusList[i-1].file)
			}
			// commit input email, name and message
			var email, name, message string
			prompt = &survey.Input{
				Message: "Please input your email:",
			}
			err = survey.AskOne(prompt, &email)
			cobra.CheckErr(err)
			prompt = &survey.Input{
				Message: "Please input your name:",
			}
			err = survey.AskOne(prompt, &name)
			cobra.CheckErr(err)
			prompt = &survey.Input{
				Message: "Please input your commit message:",
			}
			err = survey.AskOne(prompt, &message)
			cobra.CheckErr(err)
			_, err = workTree.Commit(message, &git.CommitOptions{
				Author: &object.Signature{
					Name:  name,
					Email: email,
					When:  time.Now(),
				},
			})
			cobra.CheckErr(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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

var (
	fileStatusList []fileStatus
	workTree       *git.Worktree
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ggit.yaml)")

	fileStatusList = gitStatus()
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// git the status of the current repository
func gitStatus() []fileStatus {
	var err error
	path, err := os.Getwd()
	cobra.CheckErr(err)

	r, err := git.PlainOpen(path)
	cobra.CheckErr(err)

	// get the worktree
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
