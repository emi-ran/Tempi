package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tempi",
	Short: "Tempi - Temporary folder manager",
	Long: `Tempi is a CLI tool for managing temporary folders on Windows.
It creates temporary folders with automatic expiration and cleanup.

When run without any subcommands, it creates a new temporary folder with default settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default behavior: create a new folder with default settings
		// This is equivalent to running "tempi new"
		return runNew(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Run auto-cleanup before each command (except for list and help)
		cmdName := cmd.Name()
		if cmdName != "list" && cmdName != "help" && cmdName != "completion" {
			AutoClean()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tempi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&deadtimeFlag, "deadtime", "4h", "Folder expiration time (e.g., 2h, 30m, 1h30m)")
}

