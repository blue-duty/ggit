package cmd

import (
	"fmt"
	"got/common"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type logOptions struct {
	number int
	author string
	date   string
	email  string
	file   string
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
If you use '-d' to filter the commit logs, you can input a date or a date range to filter the commit logs, like '2020-01-01' or '2020-01-01..2020-01-31'.
If you use '-f' to filter the commit logs, you can input a file path to filter the commit logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		newWorkTree()
		var (
			startDate, endDate time.Time
			err                error
		)

		ref, err := workRepo.Head()
		cobra.CheckErr(err)
		var gitLogOpts = git.LogOptions{
			All:   true,
			From:  ref.Hash(),
			Order: git.LogOrderCommitterTime,
		}

		clogs, err := workRepo.Log(&gitLogOpts)
		cobra.CheckErr(err)
		err = clogs.ForEach(func(c *object.Commit) error {
			completedLogHashes = append(completedLogHashes, c.Hash.String())
			return nil
		})
		cobra.CheckErr(err)
		// 获取commit log

		if logOpts.number == 0 {
			logOpts.number = 10
		}
		if logOpts.number < 0 {
			_, err := tm.Println(tm.Color("Invalid number", tm.RED))
			cobra.CheckErr(err)
			return
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

		if logOpts.file != "" {
			gitLogOpts.PathFilter = func(path string) bool {
				return path == logOpts.file
			}
		}

		clogs, err = workRepo.Log(&gitLogOpts)
		cobra.CheckErr(err)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Commit", "Time", "Author", "Message"})
		table.SetRowLine(true)
		table.SetRowSeparator("-")
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		//totals := tm.NewTable(0, 10, 5, ' ', 0)
		//_, err = fmt.Fprintf(totals, "Commit\tTime\tAuthor\tMessage\n")
		//cobra.CheckErr(err)
		// 遍历commit log
		var data [][]string
		var commitLogs []string
		err = clogs.ForEach(func(c *object.Commit) error {
			if logOpts.author != "" && c.Author.Name != logOpts.author {
				return nil
			}
			if logOpts.email != "" && c.Author.Email != logOpts.email {
				return nil
			}
			// 放入logHashes
			logHashes[c.Hash.String()[0:8]] = c.Hash.String()
			if logOpts.number > len(commitLogs) {
				// 放入commitLogs
				commitLogs = append(commitLogs, c.Hash.String()[0:8])
				// 放入data
				data = append(data, []string{c.Hash.String()[0:8], c.Author.When.Format(DateFormat), c.Author.Name + " <" + c.Author.Email + ">", common.Indent(c.Message)})
				//_, err2 := fmt.Fprintf(totals, "%s\t%s\t%s\t%s\n", c.Hash.String()[:8], c.Author.When.Format(DateFormat), c.Author.Name+" <"+c.Author.Email+">", c.Message)
				//cobra.CheckErr(err2)
			}
			return nil
		})
		cobra.CheckErr(err)

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output

		//tm.Clear() // Clear current screen
		//// 从第一行开始打印
		//tm.MoveCursor(0, 0)
		//_, err = tm.Println(totals)
		//cobra.CheckErr(err)
		//tm.Flush()

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

			if ch == "q" {
				break
			}

			if _, ok := logHashes[ch]; !ok {
				tm.Println(tm.Color("Invalid input", tm.RED))
				tm.Flush()
				continue
			}

			if len(logOpts.file) > 0 {
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
				err = common.ShowLog(ch, ph, logOpts.file)
				// 检查err是否为broken pipe，如果是则不报错
				if err != nil && !strings.Contains(err.Error(), "broken pipe") {
					cobra.CheckErr(err)
				}
				continue
			}
			// 查询ch的文件更改
			fileChanges, err := common.GetFileChangeByCommit(workRepo, logHashes[ch])
			cobra.CheckErr(err)
			fileChanges = append(fileChanges, "quit")
			for {
				// 选择文件
				p := &survey.Select{
					Message:  "Select a file to show the file change detail:",
					Options:  fileChanges,
					PageSize: 20,
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
				if err != nil && !strings.Contains(err.Error(), "broken pipe") {
					cobra.CheckErr(err)
				}
			}
		}
	},
}

var (
	logHashes          = make(map[string]string)
	completedLogHashes []string
)

func init() {
	logCmd.Flags().IntVarP(&logOpts.number, "number", "n", 0, "show the last <number> commit logs")
	logCmd.Flags().StringVarP(&logOpts.author, "author", "a", "", "filter the commit logs by author")
	logCmd.Flags().StringVarP(&logOpts.date, "date", "d", "", "filter the commit logs by date")
	logCmd.Flags().StringVarP(&logOpts.email, "email", "e", "", "filter the commit logs by email")
	logCmd.Flags().StringVarP(&logOpts.file, "file", "f", "", "filter the commit logs by file")
}

const DateFormat = "2006-01-02 15:04:05"
