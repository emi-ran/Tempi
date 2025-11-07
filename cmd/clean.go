package cmd

import (
	"fmt"
	"os"

	"github.com/emi-ran/Tempi/internal/cleaner"
	"github.com/emi-ran/Tempi/internal/config"
	"github.com/emi-ran/Tempi/internal/registry"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean expired temporary folders",
	Long:  `Remove all temporary folders that have exceeded their expiration time or are no longer active.`,
	RunE:  runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func runClean(cmd *cobra.Command, args []string) error {
	return cleanExpiredFolders()
}

// cleanExpiredFolders is the shared cleanup logic that can be called from other commands
func cleanExpiredFolders() error {
	// Load registry
	reg := registry.New(config.GetRegistryPath())
	if err := reg.Load(); err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Get all folders and check which ones to clean
	folders := reg.GetAll()
	cleanedCount := 0
	errors := []string{}

	for _, folder := range folders {
		shouldClean := false
		reason := ""

		// Check if folder exists
		if _, err := os.Stat(folder.Path); os.IsNotExist(err) {
			// Folder doesn't exist, remove from registry
			if err := reg.Remove(folder.Path); err != nil {
				errors = append(errors, fmt.Sprintf("failed to remove non-existent folder from registry: %s", folder.Path))
			} else {
				cleanedCount++
			}
			continue
		}

		// Check if expired based on time
		if folder.IsExpired() {
			shouldClean = true
			reason = "time expired"
		}

		// Check if inactive based on last modification time
		if !shouldClean {
			deadtime, err := config.ParseDuration(folder.Deadtime)
			if err == nil {
				isActive, err := cleaner.IsActiveBasedOnModTime(folder.Path, deadtime)
				if err == nil && !isActive {
					shouldClean = true
					reason = "no recent activity"
				}
			}
		}

		if shouldClean {
			// Try to delete the folder
			if err := cleaner.DeleteFolder(folder.Path); err != nil {
				errors = append(errors, fmt.Sprintf("failed to delete %s: %v", folder.Path, err))
			} else {
				fmt.Printf("Cleaned %s (%s)\n", folder.Path, reason)
				cleanedCount++
			}

			// Remove from registry regardless of deletion success
			if err := reg.Remove(folder.Path); err != nil {
				errors = append(errors, fmt.Sprintf("failed to remove from registry: %s", folder.Path))
			}
		}
	}

	if cleanedCount == 0 {
		fmt.Println("No expired folders to clean.")
	} else {
		fmt.Printf("\nCleaned %d expired folder(s).\n", cleanedCount)
	}

	if len(errors) > 0 {
		fmt.Println("\nErrors encountered:")
		for _, errMsg := range errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	return nil
}

// AutoClean runs cleanup automatically at startup (called from root command)
func AutoClean() {
	// Silently clean expired folders
	reg := registry.New(config.GetRegistryPath())
	if err := reg.Load(); err != nil {
		return // Silently fail on auto-clean
	}

	folders := reg.GetAll()
	for _, folder := range folders {
		shouldClean := false

		// Check if folder doesn't exist
		if _, err := os.Stat(folder.Path); os.IsNotExist(err) {
			reg.Remove(folder.Path)
			continue
		}

		// Check if expired
		if folder.IsExpired() {
			shouldClean = true
		}

		// Check if inactive
		if !shouldClean {
			deadtime, err := config.ParseDuration(folder.Deadtime)
			if err == nil {
				isActive, err := cleaner.IsActiveBasedOnModTime(folder.Path, deadtime)
				if err == nil && !isActive {
					shouldClean = true
				}
			}
		}

		if shouldClean {
			cleaner.DeleteFolder(folder.Path)
			reg.Remove(folder.Path)
		}
	}
}
