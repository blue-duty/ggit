package cmd

import (
	"errors"
	"got/common"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone a repository into a new directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("too many arguments")
		}
		if len(args) == 0 {
			return errors.New("missing repository url")
		}
		if len(args) == 1 {
			common.Clone(args[0], "")
		} else {
			common.Clone(args[0], args[1])
		}
		return nil
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "fetch from and integrate with another repository or a local branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("too many arguments")
		}
		common.Pull()
		return nil
	},
}
