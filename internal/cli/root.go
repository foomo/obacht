package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	format   string
	verbose  bool
	rulesDir string
)

var rootCmd = &cobra.Command{
	Use:   "bouncer",
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
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "pretty", "output format (pretty, json)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&rulesDir, "rules-dir", "", "path to rules directory")
}

// Execute runs the root command and exits with the appropriate code.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(Error)
	}
}
