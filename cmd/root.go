package cmd

import (
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "got",
	Short: "got is a simple git tool",
	Long: `got is a simple git tool,
which can help you to commit and push your code easily.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	fileStatusList []fileStatus
	workTree       *git.Worktree
	globalMail     string
	globalName     string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.got.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(getCmd)
}
