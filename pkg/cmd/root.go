package cmd

import (
	"github.com/kellegous/got/pkg"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	Dir      string
	Platform pkg.Platform
}

func Root() *cobra.Command {
	flags := rootFlags{
		Platform: pkg.DefaultPlatform(),
	}

	cmd := &cobra.Command{
		Use:   "got",
		Short: "Manages versions of Go in a directory",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	cmd.Flags().StringVar(
		&flags.Dir,
		"gotdir",
		"~/.go",
		"the directory where go versions will be kept")
	cmd.Flags().Var(
		&flags.Platform,
		"platform",
		"the platform to use",
	)

	cmd.AddCommand(cmdNeed())
	cmd.AddCommand(cmdUse())
	return cmd
}
