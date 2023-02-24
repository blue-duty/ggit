package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"strings"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long:  `Show commit logs`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取commit log
		clogs, err := workRepo.Log(&git.LogOptions{
			All: true,
		})
		cobra.CheckErr(err)
		// 遍历commit log
		var commitLogs []string
		err = clogs.ForEach(func(c *object.Commit) error {
			commitLogs = append(commitLogs, fmt.Sprintf(
				"%s %s Author: %sDate:   %s %s",
				plumbing.CommitObject, c.Hash, c.Author.String(),
				c.Author.When.Format(DateFormat), indent(c.Message),
			))
			return nil
		})
		cobra.CheckErr(err)

		prompt := &survey.Select{
			Message: "Select a commit log:",
			Options: commitLogs,
		}
		var commitLog string
		err = survey.AskOne(prompt, &commitLog)
		cobra.CheckErr(err)
		_, err = tm.Println(tm.Color(commitLog, tm.GREEN))
		cobra.CheckErr(err)
		tm.Flush()
	},
}

const DateFormat = "Mon Jan 02 15:04:05 2006 -0700"

func indent(t string) string {
	var output []string
	for _, line := range strings.Split(t, "\n") {
		if len(line) != 0 {
			line = "    " + line
		}

		output = append(output, line)
	}

	return strings.Join(output, "\n")
}
