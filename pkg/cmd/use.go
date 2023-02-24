package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/kellegous/got/pkg"
	"github.com/spf13/cobra"
)

func cmdUse(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "use [flags] version",
		Short: "Set current symlink to the given version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gotDir := flags.needGotDir(cmd)
			ctx := context.Background()

			var versions pkg.Versions
			norm, err := versions.Normalize(ctx, args[0])
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			name, err := pkg.Download(
				ctx,
				gotDir,
				&flags.Platform,
				norm,
				&pkg.DownloadOptions{},
			)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			link := filepath.Join(gotDir, "current")
			os.Remove(link)

			if err := os.Symlink(name, link); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}
}
