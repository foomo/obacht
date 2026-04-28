package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commitHash     = "none"
	buildTimestamp = "unknown"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		return renderVersion(cmd.OutOrStdout(), getVersionInfo(), format)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersionInfo() versionInfo {
	info := versionInfo{
		Version: version,
		Commit:  commitHash,
		Date:    buildTimestamp,
	}

	if version != "dev" {
		return info
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if info.Commit == "none" {
				info.Commit = s.Value
			}
		case "vcs.time":
			if info.Date == "unknown" {
				info.Date = s.Value
			}
		}
	}

	return info
}

func renderVersion(w io.Writer, info versionInfo, format string) error {
	if format == "json" {
		return json.NewEncoder(w).Encode(info)
	}

	_, err := fmt.Fprintf(w, "Version: %s\nCommit:  %s\nDate:    %s\n", info.Version, info.Commit, info.Date)

	return err
}
