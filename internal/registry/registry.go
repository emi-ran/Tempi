package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Registry manages temporary folder metadata
type Registry struct {
	registryPath string
	folders      []FolderMetadata
}

// New creates a new Registry instance with the given registry path
func New(registryPath string) *Registry {
	return &Registry{
		registryPath: registryPath,
		folders:      make([]FolderMetadata, 0),
	}
}

// Load reads the registry from disk
func (r *Registry) Load() error {
	// Ensure the directory exists
	dir := filepath.Dir(r.registryPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	// If file doesn't exist, start with empty registry
	if _, err := os.Stat(r.registryPath); os.IsNotExist(err) {
		r.folders = make([]FolderMetadata, 0)
		return nil
	}

	data, err := os.ReadFile(r.registryPath)
	if err != nil {
		return fmt.Errorf("failed to read registry file: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		r.folders = make([]FolderMetadata, 0)
		return nil
	}

	if err := json.Unmarshal(data, &r.folders); err != nil {
		return fmt.Errorf("failed to parse registry file: %w", err)
	}

	return nil
}

// Save writes the registry to disk
func (r *Registry) Save() error {
	data, err := json.MarshalIndent(r.folders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(r.registryPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	if err := os.WriteFile(r.registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}

	return nil
}

// Add adds a new folder to the registry
func (r *Registry) Add(folder FolderMetadata) error {
	r.folders = append(r.folders, folder)
	return r.Save()
}

// Remove removes a folder from the registry by path
func (r *Registry) Remove(path string) error {
	for i, folder := range r.folders {
		if folder.Path == path {
			r.folders = append(r.folders[:i], r.folders[i+1:]...)
			return r.Save()
		}
	}
	return nil // Not found, no error
}

// GetAll returns all folders in the registry
func (r *Registry) GetAll() []FolderMetadata {
	return r.folders
}

// GetExpired returns all expired folders
func (r *Registry) GetExpired() []FolderMetadata {
	expired := make([]FolderMetadata, 0)
	for _, folder := range r.folders {
		if folder.IsExpired() {
			expired = append(expired, folder)
		}
	}
	return expired
}

// Clear removes all entries from the registry
func (r *Registry) Clear() error {
	r.folders = make([]FolderMetadata, 0)
	return r.Save()
}

