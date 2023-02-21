package cmd

import "github.com/spf13/cobra"

func cmdUse() *cobra.Command {
	return &cobra.Command{
		Use:   "use [flags] version",
		Short: "Set current symlink to the given version",
		Args:  cobra.ExactArgs(1),
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}
