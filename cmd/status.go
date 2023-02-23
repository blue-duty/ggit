package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"got/common"
	"strconv"
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
		_, err := tm.Println(`The following is the file which can be committed ('Green' means added, 'Yellow' means modified, 'Red' means deleted, 'White' means untracked):`)
		cobra.CheckErr(err)
		tm.Flush()
		commit(fileStatusList)

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
				add := false
				if files[statusOpts.selectedFile].status == git.Added {
					add = true
				}
				if files[statusOpts.selectedFile].status == git.Deleted {
					_, err := tm.Println(tm.Color("The file has been deleted, so there is no diff.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					continue
					//add = true
				}
				err := common.ShowDiff(files[statusOpts.selectedFile].file, parentCommit, add)
				if err != nil {
					return
				}
			} else {
				_, err := tm.Println(tm.Color("You must input the serial number of the file to show the diff, or input 'q' to quit.", tm.RED))
				cobra.CheckErr(err)
				tm.Flush()
			}
		}
	},
}

func init() {
	for k, f := range fileStatusList {
		files[strconv.Itoa(k+1)] = f
	}

	head, err := workRepo.Head()
	cobra.CheckErr(err)

	c, err := workRepo.CommitObject(head.Hash())
	cobra.CheckErr(err)

	cc, err := c.Parent(0)
	cobra.CheckErr(err)
	parentCommit = cc.Hash.String()
}
