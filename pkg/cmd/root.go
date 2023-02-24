package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kellegous/got/pkg"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	Dir      string
	Platform pkg.Platform
}

func (f *rootFlags) gotDir() (string, error) {
	dir := f.Dir
	if !strings.HasPrefix(dir, "~/") {
		return dir, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, dir[2:]), nil
}

func (f *rootFlags) ensureGotDir() (string, error) {
	dir, err := f.gotDir()
	if err != nil {
		return "", err
	}
	if s, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	} else if !s.IsDir() {
		return "", fmt.Errorf("%s: not a directory", dir)
	}

	return dir, nil
}

func (f *rootFlags) needGotDir(cmd *cobra.Command) string {
	dir, err := f.ensureGotDir()
	if err != nil {
		cmd.PrintErrf("unable to access gotdir: %s", err)
		os.Exit(1)
	}
	return dir
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

	cmd.AddCommand(cmdNeed(&flags))
	cmd.AddCommand(cmdUse(&flags))
	cmd.AddCommand(cmdHas(&flags))
	cmd.AddCommand(cmdVersion())
	return cmd
}
