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

func init() {
	globalName, globalMail = getGitConfig()
	gitStatus("")
}

type fileStatus struct {
	file   string
	status git.StatusCode
}

type commitOptions struct {
	except  string
	all     bool
	message string
}

var commitOpts = &commitOptions{}

func init() {
	commitCmd.Flags().StringVarP(&commitOpts.except, "except", "e", "", "except files")
	commitCmd.Flags().BoolVarP(&commitOpts.all, "all", "a", false, "commit all files")
	commitCmd.Flags().StringVarP(&commitOpts.message, "message", "m", "", "commit message")
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
	// dir非必须参数
	Use:   "commit [<dir>] [flags]",
	Short: "Commit files to local repository",
	Long: `This command is used to commit files to the local repository by interactive mode.
If you appoint the directory, it will use the directory as the root directory of the repository.
If you don't appoint the directory, it will use the current directory as the root directory of the repository.
It can list all the files which can be committed, and you can input the file's serial number to commit, they are separated by ','.
And their status is distinguished by color. Green means added, Yellow means modified, Red means deleted, White means untracked.
If you want to commit all the files, you can use the -a flag.
If you want to commit all the files except some files, you can use the -e flag, and they are separated by ','.
If you user the -a flag or -e flag, you can also use the -m flag to appoint the commit message, if it not exist, you should input the message in the next step.
If you use the -a flag and -e flag at the same time, it will be invalid.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 判断是一个目录且存在
			if ff, err := os.Stat(args[0]); err != nil || !ff.IsDir() {
				_, err := tm.Println(tm.Color("The directory arg is not a directory or not exist.", tm.RED))
				cobra.CheckErr(err)
				tm.Flush()
				return
			}
			gitStatus(args[0])
		}
		if len(fileStatusList) == 0 {
			_, err := tm.Println(tm.Color("There is no file to commit.", tm.RED))
			cobra.CheckErr(err)
			tm.Flush()
			return
		}

		if commitOpts.except != "" && commitOpts.all {
			_, err := tm.Println(tm.Color("The except and all flags can't be used at the same time.", tm.RED))
			cobra.CheckErr(err)
			tm.Flush()
			return
		}

		if commitOpts.all {
			for _, fs := range fileStatusList {
				_, err := workTree.Add(fs.file)
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
			return
		}
		if commitOpts.except != "" {
			//exceptFiles := strings.Split(commitOpts.except, ",")
			eMap := make(map[string]bool)
			exceptFiles := strings.Split(commitOpts.except, ",")
			for _, e := range exceptFiles {
				eMap[e] = true
			}
			for _, fs := range fileStatusList {
				// get the file name
				fn := fs.file[strings.LastIndex(fs.file, "/")+1:]
				if _, ok := eMap[fn]; !ok {
					_, err := workTree.Add(fs.file)
					cobra.CheckErr(err)
				}
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
			return
		}
		_, err := tm.Println(`The following is the file which can be committed ('Green' means added, 'Yellow' means modified, 'Red' means deleted, 'White' means untracked):`)
		cobra.CheckErr(err)
		commit(fileStatusList)
		tm.Flush()
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
		if strings.Contains(fileIndex, ";") {
			var fsl = fileStatusList
			for _, index := range strings.Split(fileIndex, ";") {
				if len(strings.TrimSpace(index)) == 0 {
					_, err := tm.Println(tm.Color("You don't input the serial number of the file you want to commit.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}
				i, err := strconv.Atoi(strings.TrimSpace(index))
				cobra.CheckErr(err)
				// remove
				fsl = append(fsl[:i-1], fsl[i:]...)
			}
			for _, fs := range fsl {
				_, err = tm.Println(fs.file)
				cobra.CheckErr(err)
				tm.Flush()
				_, err := workTree.Add(fs.file)
				cobra.CheckErr(err)
			}

		} else {
			for _, index := range strings.Split(fileIndex, ",") {
				if len(strings.TrimSpace(index)) == 0 {
					_, err := tm.Println(tm.Color("You don't input the serial number of the file you want to commit.", tm.RED))
					cobra.CheckErr(err)
					tm.Flush()
					return
				}
				i, err := strconv.Atoi(strings.TrimSpace(index))
				cobra.CheckErr(err)
				_, err = tm.Println(fileStatusList[i-1].file)
				cobra.CheckErr(err)
				tm.Flush()

				// add file to staging area
				_, err = workTree.Add(fileStatusList[i-1].file)
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
	},

	Args: cobra.MaximumNArgs(1),
}

// git the status of the current repository
func gitStatus(dir string) {
	var err error
	if len(dir) == 0 {
		dir, err = os.Getwd()
		cobra.CheckErr(err)
	}

	r, err := git.PlainOpen(dir)
	cobra.CheckErr(err)

	// getCmd the worktree
	workTree, err = r.Worktree()
	cobra.CheckErr(err)

	status, err := workTree.Status()
	cobra.CheckErr(err)

	fileStatusList = make([]fileStatus, 0)
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
		fileStatusList = append(fileStatusList, fs)
	}
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
