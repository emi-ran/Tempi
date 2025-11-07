package registry

import "time"

// FolderMetadata represents metadata for a temporary folder
type FolderMetadata struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	Deadtime  string    `json:"deadtime"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired checks if the folder has expired based on ExpiresAt time
func (f *FolderMetadata) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// TimeRemaining returns the duration until expiration
func (f *FolderMetadata) TimeRemaining() time.Duration {
	return time.Until(f.ExpiresAt)
}

