package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/thalesfsp/sypl"
	"github.com/thalesfsp/sypl/level"
)

var cliLogger = sypl.NewDefault("looper-cli", level.Error)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "looper",
	Short: "Looper CLI",
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
