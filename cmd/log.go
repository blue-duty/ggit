package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"got/common"
	"strings"
	"time"
)

type logOptions struct {
	number int
	author string
	date   string
	email  string
}

var logOpts = &logOptions{}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long: `Show commit logs.
You can input a commit hash to show the file which is changed in this commit.
After that, you can select a file to show the diff in this commit.
Also you can use '-n <number>' to show the last <number> commit logs.
And '-e','-d','-a' can be used to filter the commit logs, but you can't use '-a' and '-e' at the same time.
If you use '-d' to filter the commit logs, you can input a date or a date range to filter the commit logs, like '2020-01-01' or '2020-01-01..2020-01-31'.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取commit log
		var gitLogOpts = git.LogOptions{
			All: true,
		}
		var (
			startDate, endDate time.Time
			err                error
		)
		if logOpts.number == 0 {
			logOpts.number = 10
		}
		if logOpts.author != "" && logOpts.email != "" {
			_, err := tm.Println(tm.Color("You can't use '-a' and '-e' at the same time", tm.RED))
			cobra.CheckErr(err)
			return
		}
		if logOpts.date != "" {
			if strings.Contains(logOpts.date, "..") {
				dateRange := strings.Split(logOpts.date, "..")
				if len(dateRange) != 2 {
					fmt.Println("Invalid date range")
					return
				}
				startDate, err = time.Parse("2006-01-02", dateRange[0])
				cobra.CheckErr(err)
				endDate, err = time.Parse("2006-01-02", dateRange[1])
				cobra.CheckErr(err)
				if startDate.After(endDate) {
					fmt.Println("Invalid date range")
					return
				}
				gitLogOpts.Since = &startDate
				gitLogOpts.Until = &endDate
			} else {
				startDate, err = time.Parse("2006-01-02", logOpts.date)
				if err != nil {
					fmt.Println("Invalid date")
					return
				}
				//nt := time.Now()
				gitLogOpts.Since = &startDate
				//gitLogOpts.Until = &nt
			}
		}

		clogs, err := workRepo.Log(&gitLogOpts)
		cobra.CheckErr(err)

		totals := tm.NewTable(0, 10, 5, ' ', 0)
		_, err = fmt.Fprintf(totals, "Commit\tTime\tAuther\tMessage\n")
		cobra.CheckErr(err)
		// 遍历commit log
		var commitLogs []string
		err = clogs.ForEach(func(c *object.Commit) error {
			if logOpts.author != "" && c.Author.Name != logOpts.author {
				return nil
			}
			if logOpts.email != "" && c.Author.Email != logOpts.email {
				return nil
			}
			// 放入commitLogs
			commitLogs = append(commitLogs, c.Hash.String()[0:8])
			// 放入logHashes
			logHashes[c.Hash.String()[0:8]] = c.Hash.String()
			_, err2 := fmt.Fprintf(totals, "%s\t%s\t%s\t%s\n", c.Hash.String()[:8], c.Author.When.Format(DateFormat), c.Author.Name+" <"+c.Author.Email+">", c.Message)
			cobra.CheckErr(err2)
			return nil
		})
		cobra.CheckErr(err)

		tm.Clear() // Clear current screen
		// 从第一行开始打印
		tm.MoveCursor(0, 0)
		_, err = tm.Println(totals)
		cobra.CheckErr(err)
		tm.Flush()

		for {
			pp := &survey.Input{
				Message: "Input a commit hash to show the file which is changed in this commit:",
				Suggest: func(toComplete string) []string {
					var suggestions []string
					for _, s := range commitLogs {
						if strings.HasPrefix(s, toComplete) {
							suggestions = append(suggestions, s)
						}
					}
					return suggestions
				},
			}

			var ch string
			err = survey.AskOne(pp, &ch)
			cobra.CheckErr(err)

			if ch == "" || ch == "q" {
				break
			}

			// 查询ch的文件更改
			fileChanges, err := common.GetFileChangeByCommit(workRepo, logHashes[ch])
			cobra.CheckErr(err)

			fileChanges = append(fileChanges, "quit")

			for {
				// 选择文件
				p := &survey.Select{
					Message: "Select a file to show the file change detail:",
					Options: fileChanges,
				}

				var fc string
				err = survey.AskOne(p, &fc)
				cobra.CheckErr(err)

				if fc == "quit" {
					break
				}

				// 获取Commit的父节点
				var ph string
				for k, h := range completedLogHashes {
					if h == logHashes[ch] {
						if k == len(completedLogHashes)-1 {
							ph = ""
						} else {
							ph = completedLogHashes[k+1]
						}
						break
					}
				}

				// 使用git log -p命令获取commit的详细信息
				err = common.ShowLog(ch, ph, fc)
				cobra.CheckErr(err)
			}
		}
	},
}

var (
	logHashes          = make(map[string]string)
	completedLogHashes []string
)

func init() {
	clogs, err := workRepo.Log(&git.LogOptions{
		All: true,
	})
	cobra.CheckErr(err)

	err = clogs.ForEach(func(c *object.Commit) error {
		completedLogHashes = append(completedLogHashes, c.Hash.String())
		return nil
	})
	cobra.CheckErr(err)

	logCmd.Flags().IntVarP(&logOpts.number, "number", "n", 0, "show the last <number> commit logs")
	logCmd.Flags().StringVarP(&logOpts.author, "author", "a", "", "filter the commit logs by author")
	logCmd.Flags().StringVarP(&logOpts.date, "date", "d", "", "filter the commit logs by date")
	logCmd.Flags().StringVarP(&logOpts.email, "email", "e", "", "filter the commit logs by email")
}

const DateFormat = "2006-01-02 15:04:05"
