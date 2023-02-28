package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [<url>] [flags]",
	Short: "Get a file from remote repository",
	Long: `Get a file from remote repository.
You use the url of the file in the remote repository.
With no arguments, it will exit with an error.
With '--output' flag, it will output the file to the specified path.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			getOpts.url = args[0]
			return get(getOpts)
		} else {
			if _, err := tm.Println(tm.Color("You must specify a url.", tm.RED)); err != nil {
				return err
			}
			tm.Flush()
			return errors.New("you must specify a url")
		}
	},
}

type getOptions struct {
	output string
	url    string
	proxy  string
}

var getOpts = &getOptions{}

func init() {
	getCmd.Flags().StringVarP(&getOpts.output, "output", "o", "", "output file path")
	getCmd.Flags().StringVarP(&getOpts.proxy, "proxy", "p", "", "proxy url")
}

func get(getOpts *getOptions) error {
	// 判断是url还是path
	// 如果是url，直接下载
	// 如果是path，转为url，下载
	if isNotURL(getOpts.url) {
		return errors.New("not a url")
	}
	// download file from url
	// 替换url中的github.com为raw.githubusercontent.com，删掉blob
	getOpts.url = strings.Replace(getOpts.url, "github.com", "raw.githubusercontent.com", 1)
	getOpts.url = strings.Replace(getOpts.url, "blob/", "", 1)
	return getFile(getOpts)
}

func isNotURL(url string) bool {
	return !strings.HasPrefix(url, "https://")
}

func getFile(getOpts *getOptions) error {
	// download file from url
	// 如果有--output，下载到指定路径
	// 如果没有--output，下载到当前路径
	if len(getOpts.output) > 0 {
		// download to output path
		// 取出目录，判断路径是否存在
		dir := getOpts.output[:strings.LastIndex(getOpts.output, "/")]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// 目录不存在，创建目录
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
		// 判断文件是否存在
		if _, err := os.Stat(getOpts.output); err == nil {
			// 文件存在，提示会被覆盖
			_, err := tm.Println(tm.Color("File already exists, it will be overwritten.", tm.RED))
			if err != nil {
				return err
			}
			tm.Flush()
			// 删除文件
			err = os.Remove(getOpts.output)
			if err != nil {
				return err
			}
		}
	} else {
		// download to current path
		// 取出文件名
		fileName := getOpts.url[strings.LastIndex(getOpts.url, "/")+1:]
		// 获取当前路径
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		// 判断文件是否存在
		if _, err := os.Stat(path.Join(dir, fileName)); err == nil {
			// 文件存在，提示会被覆盖
			_, err := tm.Println(tm.Color("File already exists, it will be overwritten.", tm.RED))
			if err != nil {
				return err
			}
			tm.Flush()
			// 删除文件
			err = os.Remove(path.Join(dir, fileName))
			if err != nil {
				return err
			}
		}

		getOpts.output = path.Join(dir, fileName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if len(getOpts.proxy) > 0 {
		cmd = exec.CommandContext(ctx, "curl", "-x", getOpts.proxy, "-o", getOpts.output, getOpts.url)
	} else {
		cmd = exec.CommandContext(context.TODO(), "curl", "-o", getOpts.output, getOpts.url)
	}

	tm.Println(tm.Color("Downloading...", tm.GREEN))
	tm.Flush()

	if err := cmd.Run(); err != nil {
		return err
	}
	if ctx.Err() == context.DeadlineExceeded {
		tm.Println(tm.Color("Download timeout.You can try again with a proxy.", tm.RED))
		tm.Flush()
	}

	tm.Println(tm.Color("Download success.The file is saved in "+getOpts.output, tm.GREEN))
	tm.Flush()

	return nil
}
