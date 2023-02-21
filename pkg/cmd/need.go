package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/kellegous/got/pkg"
)

func cmdNeed() *cobra.Command {
	return &cobra.Command{
		Use:     "need [flags] version",
		Aliases: []string{"get", "download", "dl"},
		Short:   "Downloads the specified version of Go",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			latest, err := pkg.GetLatestVersion(ctx)
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("latest = %s\n", latest)
		},
	}
}
