package cmd

import (
	"got/common"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	selectedFile string
}

var statusOpts = &statusOptions{}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long: `Show the working tree status.
This command shows the working tree status.`,
	Run: func(cmd *cobra.Command, args []string) {
		newWorkTree()
		fileStatus, err := workTree.Status()
		cobra.CheckErr(err)
		initFiles(fileStatus)

		head, err := workRepo.Head()
		cobra.CheckErr(err)
		parentCommit = head.Hash().String()

		_, err = tm.Println(`The following is the all status of the files which are not Unmodified.`)
		cobra.CheckErr(err)
		tm.Flush()
		printStatus()
		for {
			prompt := &survey.Input{
				Message: "You can input the serial number of the file to show the diff, or input 'q' to quit:",
			}

			err := survey.AskOne(prompt, &statusOpts.selectedFile)
			if err != nil {
				return
			}

			if statusOpts.selectedFile == "q" {
				return
			}

			if _, ok := files[statusOpts.selectedFile]; ok {
				if files[statusOpts.selectedFile].staging == git.Added || files[statusOpts.selectedFile].worktree == git.Added {
					_, err := tm.Println(tm.Color("The file has been added, so there is no diff.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					continue
				}
				if files[statusOpts.selectedFile].staging == git.Untracked || files[statusOpts.selectedFile].worktree == git.Untracked {
					_, err := tm.Println(tm.Color("The file is untracked, so there is no diff.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					continue
				}
				if files[statusOpts.selectedFile].staging == git.Deleted || files[statusOpts.selectedFile].worktree == git.Deleted {
					_, err := tm.Println(tm.Color("The file has been deleted, so there is no diff.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					continue
				}
				err := common.ShowDiff(files[statusOpts.selectedFile].file, parentCommit)
				if err != nil && !strings.Contains(err.Error(), "broken pipe") {
					cobra.CheckErr(err)
				}
			} else {
				_, err := tm.Println(tm.Color("You must input the serial number of the file to show the diff, or input 'q' to quit.", tm.RED))
				cobra.CheckErr(err)
				tm.Flush()
			}
		}
	},
}
