package cli

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var (
	format   string
	verbose  bool
	rulesDir string
)

var (
	version    = "dev"
	commitHash = "none"
)

var rootCmd = &cobra.Command{
	Use:   "obacht",
	Short: "Security configuration scanner for developer environments",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch format {
		case "pretty", "json":
			// valid
		default:
			return fmt.Errorf("invalid format %q: must be pretty or json", format)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "pretty", "output format (pretty, json)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&rulesDir, "rules-dir", "", "use rules from this directory instead of embedded rules")
}

// Execute runs the root command and exits with the appropriate code.
func Execute() {
	err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version),
		fang.WithCommit(commitHash),
		fang.WithNotifySignal(os.Interrupt, syscall.SIGTERM),
	)
	if err != nil {
		os.Exit(Error)
	}
}
