package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

			b, err := json.MarshalIndent(bi.Raw, "", "  ")
			if err != nil {
				log.Panic(err)
			}
			fmt.Printf("%s\n", b)
		},
	}
}

type buildInfo struct {
	SHA  string
	Time time.Time

	Raw *debug.BuildInfo
}

func (i *buildInfo) Name() string {
	return buildname.FromVersion(i.SHA)
}

func readBuildInfo() (*buildInfo, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, errors.New("build info unavailable")
	}

	i := buildInfo{
		Raw: bi,
	}

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
