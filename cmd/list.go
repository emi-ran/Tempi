package cmd

import (
	"fmt"

	"github.com/emi-ran/Tempi/internal/config"
	"github.com/emi-ran/Tempi/internal/registry"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active temporary folders",
	Long:  `Display all tracked temporary folders with their paths and remaining time until expiration.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Load registry
	reg := registry.New(config.GetRegistryPath())
	if err := reg.Load(); err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	folders := reg.GetAll()

	if len(folders) == 0 {
		fmt.Println("No active temporary folders.")
		return nil
	}

	fmt.Printf("Active temporary folders (%d):\n\n", len(folders))

	for i, folder := range folders {
		remaining := folder.TimeRemaining()
		remainingStr := config.FormatDuration(remaining)
		
		status := ""
		if remaining < 0 {
			status = " (expired)"
		}

		fmt.Printf("[%d] %s (expires in %s)%s\n", i+1, folder.Path, remainingStr, status)
	}

	return nil
}

