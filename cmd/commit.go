package cmd

import (
	"fmt"
	"got/common"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

// type fileStatus struct {
// 	file   string
// 	status git.StatusCode
// }

type commitOptions struct {
	except  string
	all     bool
	message string
	push    bool
}

var commitOpts = &commitOptions{}

func init() {
	globalName, globalMail = getGitConfig()
	commitCmd.Flags().StringVarP(&commitOpts.except, "except", "e", "", "except files")
	commitCmd.Flags().BoolVarP(&commitOpts.all, "all", "a", false, "commit all files")
	commitCmd.Flags().StringVarP(&commitOpts.message, "message", "m", "", "commit message")
	commitCmd.Flags().BoolVarP(&commitOpts.push, "push", "p", false, "push to remote repository")
}

/*
// Print fileStatus print
func (fs fileStatus) Print(i int) {
	switch fs.status {
	case git.Untracked:
		// print serial number and file name and format alignment
		serial := fmt.Sprintf("%d. %s", i, fs.file)
		tm.Printf("%s", tm.Color(serial, tm.WHITE))
		tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.WHITE))
		tm.Printf("%s\n", tm.Color("Untracked", tm.WHITE))
		tm.Flush()
	case git.Added:
		serial := fmt.Sprintf("%d. %s", i, fs.file)
		tm.Printf("%s", tm.Color(serial, tm.GREEN))
		tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.GREEN))
		tm.Printf("%s\n", tm.Color("Added", tm.GREEN))
		tm.Flush()
	case git.Modified:
		serial := fmt.Sprintf("%d. %s", i, fs.file)
		tm.Printf("%s", tm.Color(serial, tm.YELLOW))
		tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.YELLOW))
		tm.Printf("%s\n", tm.Color("Modified", tm.YELLOW))
		tm.Flush()
	case git.Deleted:
		serial := fmt.Sprintf("%d. %s", i, fs.file)
		tm.Printf("%s", tm.Color(serial, tm.RED))
		tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.RED))
		tm.Printf("%s\n", tm.Color("Deleted", tm.RED))
		tm.Flush()
	default:
		serial := fmt.Sprintf("%d. %s", i, fs.file)
		tm.Printf("%s", tm.Color(serial, tm.WHITE))
		tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.WHITE))
		tm.Printf("%s\n", tm.Color("Untracked", tm.WHITE))
		tm.Flush()
	}
}

func (fs fileStatus) String() string {
	switch fs.status {
	case git.Untracked:
		return tm.Color(fs.file, tm.WHITE)
	case git.Added:
		return tm.Color(fs.file, tm.GREEN)
	case git.Modified:
		return tm.Color(fs.file, tm.YELLOW)
	case git.Deleted:
		return tm.Color(fs.file, tm.RED)
	default:
		return tm.Color(fs.file, tm.WHITE)
	}
}


func printStatus(fileStatusList []fileStatus) {
	for i, fs := range fileStatusList {
		fs.Print(i + 1)
	}
	tm.Flush()
}
*/

func initFiles(fs git.Status) {
	if len(fs) == 0 {
		_, err := tm.Println(tm.Color("There is no file to printStatus.", tm.RED))
		cobra.CheckErr(err)
		tm.Flush()
		os.Exit(1)
	}
	var i = 1
	for file, s := range fs {
		if s.Staging == git.Deleted || s.Worktree == git.Deleted {
			files[strconv.Itoa(i)] = newGotFileStatus(file, s)
			i++
		} else if s.Staging == git.Added || s.Worktree == git.Added {
			files[strconv.Itoa(i)] = newGotFileStatus(file, s)
			i++
		} else if s.Staging == git.Modified || s.Worktree == git.Modified {
			files[strconv.Itoa(i)] = newGotFileStatus(file, s)
			i++
		} else {
			files[strconv.Itoa(i)] = newGotFileStatus(file, s)
			i++
		}
	}
}

func printStatus() {
	for i := 1; i <= len(files); i++ {
		if file, ok := files[strconv.Itoa(i)]; ok {
			if file.staging == git.Deleted || file.worktree == git.Deleted {
				serial := fmt.Sprintf("%d. %s", i, file.file)
				tm.Printf("%s", tm.Color(serial, tm.RED))
				tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.RED))
				tm.Printf("%s\n", tm.Color("Deleted", tm.RED))
				tm.Flush()
			} else if file.staging == git.Added || file.worktree == git.Added {
				serial := fmt.Sprintf("%d. %s", i, file.file)
				tm.Printf("%s", tm.Color(serial, tm.GREEN))
				tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.GREEN))
				tm.Printf("%s\n", tm.Color("Added", tm.GREEN))
				tm.Flush()
			} else if file.staging == git.Modified || file.worktree == git.Modified {
				serial := fmt.Sprintf("%d. %s", i, file.file)
				tm.Printf("%s", tm.Color(serial, tm.YELLOW))
				tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.YELLOW))
				tm.Printf("%s\n", tm.Color("Modified", tm.YELLOW))
				tm.Flush()
			} else if file.staging == git.Untracked || file.worktree == git.Untracked {
				serial := fmt.Sprintf("%d. %s", i, file.file)
				tm.Printf("%s", tm.Color(serial, tm.WHITE))
				tm.Printf("%s", tm.Color(strings.Repeat("-", 50-len(serial)), tm.WHITE))
				tm.Printf("%s\n", tm.Color("Untracked", tm.WHITE))
				tm.Flush()
			}
		}
	}
}

var commitCmd = &cobra.Command{
	// dir非必须参数
	Use:   "commit [flags]",
	Short: "Commit files to local repository",
	Long: `This command is used to commit files to the local repository by interactive mode.
It can list all the files which can be committed at the current directory, and you can input the file's serial number to commit, they are separated by ','.
And their status is distinguished by color. Green means added, Yellow means modified, Red means deleted, White means untracked.
If you want to commit all the files, you can use the -a flag.
If you want to commit all the files except some files, you can use the -e flag, and they are separated by ','.
If you want to push after commit, you can use the -p flag.
If you user the -a flag or -e flag, you can also use the -m flag to appoint the commit message, if it not exist, you should input the message in the next step.
If you use the -a flag and -e flag at the same time, it will be invalid.`,
	Run: func(cmd *cobra.Command, args []string) {
		if commitOpts.except != "" && commitOpts.all {
			_, err := tm.Println(tm.Color("The except and all flags can't be used at the same time.", tm.RED))
			cobra.CheckErr(err)
			tm.Flush()
			return
		}

		newWorkTree()

		if commitOpts.all {
			err := workTree.AddWithOptions(&git.AddOptions{
				All: true,
			})
			cobra.CheckErr(err)
			if commitOpts.message == "" {
				var message string
				prompt := &survey.Input{
					Message: "Please input your commit message:",
				}
				err := survey.AskOne(prompt, &message)
				cobra.CheckErr(err)
				_, err = workTree.Commit(message, &git.CommitOptions{
					Author: &object.Signature{
						Name:  globalName,
						Email: globalMail,
						When:  time.Now(),
					},
				})
				cobra.CheckErr(err)
			} else {
				_, err := workTree.Commit(commitOpts.message, &git.CommitOptions{
					Author: &object.Signature{
						Name:  globalName,
						Email: globalMail,
						When:  time.Now(),
					},
				})
				cobra.CheckErr(err)
			}
			_, err = tm.Println(tm.Color("Commit successfully.", tm.GREEN))
			cobra.CheckErr(err)
			tm.Flush()
			if commitOpts.push {
				err := common.Push()
				if err != nil {
					_, err := tm.Println(tm.Color("Push failed.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}
				_, err = tm.Println(tm.Color("Push successfully.", tm.GREEN))
				cobra.CheckErr(err)
				tm.Flush()
			}
			return
		}

		fileStatus, err := workTree.Status()
		cobra.CheckErr(err)
		initFiles(fileStatus)

		if commitOpts.except != "" {
			iMap := files
			exceptFiles := strings.Split(commitOpts.except, ",")
			for _, ef := range exceptFiles {
				delete(iMap, ef)
			}

			for _, fs := range iMap {
				// get the file name
				fn := fs.file
				_, err := workTree.Add(fn)
				cobra.CheckErr(err)
			}

			if commitOpts.message == "" {
				var message string
				prompt := &survey.Input{
					Message: "Please input your commit message:",
				}
				err := survey.AskOne(prompt, &message)
				cobra.CheckErr(err)
				_, err = workTree.Commit(message, &git.CommitOptions{
					Author: &object.Signature{
						Name:  globalName,
						Email: globalMail,
						When:  time.Now(),
					},
				})
				cobra.CheckErr(err)
			} else {
				_, err := workTree.Commit(commitOpts.message, &git.CommitOptions{
					Author: &object.Signature{
						Name:  globalName,
						Email: globalMail,
						When:  time.Now(),
					},
				})
				cobra.CheckErr(err)
			}
			_, err := tm.Println(tm.Color("Commit successfully.", tm.GREEN))
			cobra.CheckErr(err)
			tm.Flush()
			if commitOpts.push {
				err := common.Push()
				if err != nil {
					_, err := tm.Println(tm.Color("Push failed.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}
				_, err = tm.Println(tm.Color("Push successfully.", tm.GREEN))
				cobra.CheckErr(err)
				tm.Flush()
			}
			return
		}

		tm.Println(`The following is the file which can be committed: `)
		cobra.CheckErr(err)
		tm.Flush()
		printStatus()
		var fileIndex string
		prompt := &survey.Input{
			Message: "Please input the serial number of the file (You can use ',' to separate you want to commit or use ';' to separate you don't want to commit):",
		}
		err = survey.AskOne(prompt, &fileIndex)
		cobra.CheckErr(err)
		// split by ,
		if len(fileIndex) == 0 {
			_, err := tm.Println(tm.Color("You don't input the serial number of the file you want to commit.", tm.RED))
			cobra.CheckErr(err)
			tm.Flush()
			return
		}
		if fileIndex == "q" || fileIndex == "Q" {
			return
		}
		if strings.Contains(fileIndex, ";") {
			var fsl = files
			for _, index := range strings.Split(fileIndex, ";") {
				if len(strings.TrimSpace(index)) == 0 {
					_, err := tm.Println(tm.Color("You don't input the serial number of the file you want to commit.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}
				// remove
				delete(fsl, index)
			}
			for _, fs := range fsl {
				_, err := workTree.Add(fs.file)
				cobra.CheckErr(err)
				fmt.Printf("Add %s to staging area successfully.\n", fs.file)
			}
		} else {
			for _, index := range strings.Split(fileIndex, ",") {
				if len(strings.TrimSpace(index)) == 0 {
					_, err := tm.Println(tm.Color("You don't input the serial number of the file you want to commit.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}

				// add file to staging area
				_, err = workTree.Add(files[index].file)
				cobra.CheckErr(err)
				fmt.Printf("Add %s to staging area successfully.\n", files[index].file)
			}
		}

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
		_, err = tm.Println(tm.Color("Commit successfully.", tm.GREEN))
		cobra.CheckErr(err)

		if commitOpts.push {
			err := common.Push()
			if err != nil {
				_, err := tm.Println(tm.Color("Push failed.", tm.RED))
				cobra.CheckErr(err)
				tm.Flush()
				return
			}
			_, err = tm.Println(tm.Color("Push successfully.", tm.GREEN))
			cobra.CheckErr(err)
			tm.Flush()
		}
	},
}
