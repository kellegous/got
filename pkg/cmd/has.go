package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kellegous/got/pkg"
	"github.com/spf13/cobra"
)

func cmdHas(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "has [flags] [version]",
		Aliases: []string{"versions"},
		Short:   "List existing versions or check for a version",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gotDir := flags.needGotDir(cmd)

			ctx := context.Background()
			if len(args) == 1 {
				var versions pkg.Versions
				norm, err := versions.Normalize(ctx, args[0])
				if err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}

				if _, err := os.Stat(filepath.Join(gotDir, norm)); err != nil {
					fmt.Printf("%s 👍\n", norm)
				} else {
					fmt.Printf("%s 👎\n", norm)
				}
				return
			}

			cver, cplat, err := getCurrent(gotDir)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			files, err := ioutil.ReadDir(gotDir)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			for _, file := range files {
				if !file.IsDir() {
					continue
				}

				version, platform, err := parseDirname(file.Name())
				if err != nil {
					continue
				}

				if version == cver &&
					platform.OS == cplat.OS &&
					platform.Arch == cplat.Arch {
					fmt.Printf("%s\t%s\t✅\n", version, platform)
				} else {
					fmt.Printf("%s\t%s\n", version, platform)
				}
			}
		},
	}

	return cmd
}

func getCurrent(gotDir string) (string, *pkg.Platform, error) {
	name, err := os.Readlink(filepath.Join(gotDir, "current"))
	if err != nil {
		return "", nil, err
	}

	return parseDirname(name)
}

func parseDirname(name string) (string, *pkg.Platform, error) {
	elems := strings.SplitN(name, "-", 3)
	if len(elems) != 3 {
		return "", nil, fmt.Errorf("not a go directory: %s", name)
	}
	return elems[0], &pkg.Platform{OS: elems[1], Arch: elems[2]}, nil
}
