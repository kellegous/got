package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/kellegous/got/pkg"
)

func cmdNeed(flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "need [flags] version",
		Aliases: []string{"get", "download", "dl"},
		Short:   "Downloads the specified version of Go",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gotDir, err := flags.ensureGotDir()
			if err != nil {
				cmd.PrintErrf("unable to determine got dir: %s\n", err)
				os.Exit(1)
			}

			ctx := context.Background()
			seen := map[string]bool{}
			var versions pkg.Versions
			for _, version := range args {
				norm, err := versions.Normalize(ctx, version)
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}

				if seen[norm] {
					continue
				}
				seen[norm] = true

				if _, err := pkg.Download(
					ctx,
					gotDir,
					&flags.Platform,
					norm,
					&pkg.DownloadOptions{},
				); err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
			}
		},
	}
}
