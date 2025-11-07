package cmd

import (
	"fmt"

	"github.com/emi-ran/Tempi/internal/cleaner"
	"github.com/emi-ran/Tempi/internal/config"
	"github.com/emi-ran/Tempi/internal/registry"
	"github.com/spf13/cobra"
)

// deletenowCmd represents the deletenow command
var deletenowCmd = &cobra.Command{
	Use:   "deletenow",
	Short: "Delete all temporary folders immediately",
	Long:  `Delete all tracked temporary folders immediately, regardless of their expiration time.`,
	RunE:  runDeleteNow,
}

func init() {
	rootCmd.AddCommand(deletenowCmd)
}

func runDeleteNow(cmd *cobra.Command, args []string) error {
	// Load registry
	reg := registry.New(config.GetRegistryPath())
	if err := reg.Load(); err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	folders := reg.GetAll()

	if len(folders) == 0 {
		fmt.Println("No temporary folders to delete.")
		return nil
	}

	deletedCount := 0
	errors := []string{}

	for _, folder := range folders {
		// Try to delete the folder
		if err := cleaner.DeleteFolder(folder.Path); err != nil {
			errors = append(errors, fmt.Sprintf("failed to delete %s: %v", folder.Path, err))
		} else {
			fmt.Printf("Deleted %s\n", folder.Path)
			deletedCount++
		}
	}

	// Clear the registry
	if err := reg.Clear(); err != nil {
		return fmt.Errorf("failed to clear registry: %w", err)
	}

	fmt.Printf("\nDeleted %d folder(s).\n", deletedCount)

	if len(errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, errMsg := range errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	return nil
}

