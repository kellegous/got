package cmd

import (
	"context"
	"fmt"
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
			ctx := context.Background()
			normalized, err := normalizeVersions(ctx, args)
			if err != nil {
				cmd.PrintErrln(err.Error())
				os.Exit(1)
			}

			gotDir, err := flags.gotDir()
			if err != nil {
				cmd.PrintErrf("unable to determine got dir: %s", err)
				os.Exit(1)
			}

			for _, version := range normalized {
				fmt.Printf("gotdir = %s, version = %s, platform = %s\n", gotDir, version, &flags.Platform)
			}
		},
	}
}

func normalizeVersions(
	ctx context.Context,
	versions []string,
) ([]string, error) {
	s := pkg.NewOrderedSet[string]()
	var latest string
	var err error
	for _, version := range versions {
		if version == "latest" {
			if latest == "" {
				latest, err = pkg.GetLatestVersion(ctx)
				if err != nil {
					return nil, err
				}
			}
			s.Add(latest)
			continue
		}

		n, err := pkg.NormalizeVersion(version)
		if err != nil {
			return nil, err
		}

		s.Add(n)
	}
	return s.Values(), nil
}
