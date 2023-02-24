package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/kellegous/buildname"
	"github.com/spf13/cobra"
)

func cmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Get the version of this program",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			bi, err := readBuildInfo()
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			fmt.Printf("Version: %s\n", bi.SHA)
			fmt.Printf("Name:    %s\n", bi.Name())
			fmt.Printf("Time:    %s\n", bi.Time.Format(time.RFC3339))
		},
	}
}

type buildInfo struct {
	SHA  string
	Time time.Time
}

func (i *buildInfo) Name() string {
	return buildname.FromVersion(i.SHA)
}

func readBuildInfo() (*buildInfo, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("build info unavailable")
	}

	var i buildInfo
	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			i.SHA = setting.Value
		case "vcs.time":
			t, err := time.Parse(time.RFC3339, setting.Value)
			if err != nil {
				return nil, err
			}
			i.Time = t
		}
	}

	return &i, nil
}
