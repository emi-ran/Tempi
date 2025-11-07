package cleaner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetLastModifiedTime returns the most recent modification time in a directory tree
func GetLastModifiedTime(dirPath string) (time.Time, error) {
	var latestTime time.Time

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip errors for individual files (e.g., permission denied)
			return nil
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return nil // Skip files we can't stat
			}

			if info.ModTime().After(latestTime) {
				latestTime = info.ModTime()
			}
		}

		return nil
	})

	if err != nil {
		return time.Time{}, err
	}

	// If no files found, use directory's own modification time
	if latestTime.IsZero() {
		info, err := os.Stat(dirPath)
		if err != nil {
			return time.Time{}, err
		}
		latestTime = info.ModTime()
	}

	return latestTime, nil
}

// IsActiveBasedOnModTime checks if a folder is still active based on modification time
// Returns true if any file was modified within the deadtime period
func IsActiveBasedOnModTime(dirPath string, deadtime time.Duration) (bool, error) {
	lastModified, err := GetLastModifiedTime(dirPath)
	if err != nil {
		return false, err
	}

	timeSinceModification := time.Since(lastModified)
	return timeSinceModification < deadtime, nil
}

// GetProcessesUsingPath returns a list of process IDs using files in the given path
func GetProcessesUsingPath(dirPath string) ([]uint32, error) {
	// This is a simplified implementation
	// A full implementation would enumerate all processes and check their open file handles
	// For now, we'll return an empty list and rely on Windows delete-on-close semantics
	return []uint32{}, nil
}

// TerminateProcess attempts to terminate a process by PID
func TerminateProcess(pid uint32) error {
	handle, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return fmt.Errorf("failed to open process %d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	err = windows.TerminateProcess(handle, 1)
	if err != nil {
		return fmt.Errorf("failed to terminate process %d: %w", pid, err)
	}

	return nil
}

// DeleteFolder removes a folder and all its contents
// It attempts to handle processes using the folder
func DeleteFolder(dirPath string) error {
	// First, try to find and terminate processes using the folder
	pids, err := GetProcessesUsingPath(dirPath)
	if err != nil {
		// Log but continue - we'll try to delete anyway
	} else {
		for _, pid := range pids {
			_ = TerminateProcess(pid) // Ignore errors
		}
		// Give processes time to cleanup
		if len(pids) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Try to remove the folder
	err = os.RemoveAll(dirPath)
	if err != nil {
		// If removal failed, try to mark files for deletion on reboot
		if markErr := markForDeletionOnReboot(dirPath); markErr != nil {
			return fmt.Errorf("failed to delete folder and mark for deletion: %w", err)
		}
		return fmt.Errorf("folder marked for deletion on reboot: %s", dirPath)
	}

	return nil
}

// markForDeletionOnReboot uses MoveFileEx with MOVEFILE_DELAY_UNTIL_REBOOT
func markForDeletionOnReboot(path string) error {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procMoveFileEx := modkernel32.NewProc("MoveFileExW")

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	const MOVEFILE_DELAY_UNTIL_REBOOT = 0x4

	ret, _, err := procMoveFileEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0, // NULL destination means delete
		uintptr(MOVEFILE_DELAY_UNTIL_REBOOT),
	)

	if ret == 0 {
		return fmt.Errorf("MoveFileEx failed: %w", err)
	}

	return nil
}

