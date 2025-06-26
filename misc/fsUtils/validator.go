package fsUtils

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidatePath checks if a path exists and is accessible
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	// Check if the path exists
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error accessing path %s: %w", path, err)
	}
	
	return nil
}

// ValidateDirectory checks if a directory exists and is accessible
func ValidateDirectory(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	
	// Check if the directory exists and is a directory
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("error accessing directory %s: %w", dirPath, err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dirPath)
	}
	
	return nil
}

// ValidateWritableDirectory checks if a directory is writable
func ValidateWritableDirectory(dirPath string) error {
	if err := ValidateDirectory(dirPath); err != nil {
		return err
	}
	
	// Create a temporary file to test write permissions
	tempFile := filepath.Join(dirPath, ".write-test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %s", dirPath)
	}
	
	// Clean up the test file
	file.Close()
	os.Remove(tempFile)
	
	return nil
}
import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidatePath checks if a path exists and is accessible
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	// Check if the path exists
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("error accessing path %s: %w", path, err)
	}
	
	return nil
}

// ValidateDirectory checks if a directory exists and is accessible
func ValidateDirectory(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	
	// Check if the directory exists and is a directory
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dirPath)
		}
		return fmt.Errorf("error accessing directory %s: %w", dirPath, err)
	}
	
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dirPath)
	}
	
	return nil
}

// ValidateWritableDirectory checks if a directory is writable
func ValidateWritableDirectory(dirPath string) error {
	if err := ValidateDirectory(dirPath); err != nil {
		return err
	}
	
	// Create a temporary file to test write permissions
	tempFile := filepath.Join(dirPath, ".write-test")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("directory is not writable: %s", dirPath)
	}
	
	// Clean up the test file
	file.Close()
	os.Remove(tempFile)
	
	return nil
}
