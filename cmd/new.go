package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/emi-ran/Tempi/internal/config"
	"github.com/emi-ran/Tempi/internal/registry"
	"github.com/spf13/cobra"
)

var deadtimeFlag string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new temporary folder",
	Long:  `Create a new temporary folder with a specified or default deadtime (expiration time).`,
	RunE:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVar(&deadtimeFlag, "deadtime", "4h", "Folder expiration time (e.g., 2h, 30m, 1h30m)")
}

func runNew(cmd *cobra.Command, args []string) error {
	// Parse deadtime
	deadtime, err := config.ParseDuration(deadtimeFlag)
	if err != nil {
		return fmt.Errorf("invalid deadtime format: %w", err)
	}

	// Generate folder name with timestamp
	timestamp := time.Now().Format("20060102_150405")
	folderName := fmt.Sprintf("tempi_%s", timestamp)
	
	// Get temp directory
	tempDir := config.GetTempDir()
	folderPath := filepath.Join(tempDir, folderName)

	// Create the directory
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}

	// Create metadata
	now := time.Now()
	metadata := registry.FolderMetadata{
		Path:      folderPath,
		CreatedAt: now,
		Deadtime:  deadtimeFlag,
		ExpiresAt: now.Add(deadtime),
	}

	// Load registry and add folder
	reg := registry.New(config.GetRegistryPath())
	if err := reg.Load(); err != nil {
		// Cleanup folder if registry fails
		os.RemoveAll(folderPath)
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if err := reg.Add(metadata); err != nil {
		// Cleanup folder if adding to registry fails
		os.RemoveAll(folderPath)
		return fmt.Errorf("failed to add folder to registry: %w", err)
	}

	fmt.Printf("Created folder %s\n", folderPath)
	fmt.Printf("This folder will expire in %s.\n\n", config.FormatDuration(deadtime))
	fmt.Printf("cd %s\n", folderPath)

	return nil
}

