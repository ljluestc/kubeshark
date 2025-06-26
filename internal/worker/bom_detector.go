package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DetectBOM checks if a file has a BOM (Byte Order Mark)
func DetectBOM(filename string) (bool, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// UTF-8 BOM is EF BB BF
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return true, nil
	}

	return false, nil
}

// ScanForBOM scans for BOM in all Go files in a directory
func ScanForBOM(dir string) ([]string, error) {
	var filesWithBOM []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			hasBOM, err := DetectBOM(path)
			if err != nil {
				return err
			}

			if hasBOM {
				filesWithBOM = append(filesWithBOM, path)
				fmt.Printf("Found BOM in file: %s\n", path)
			}
		}

		return nil
	})

	return filesWithBOM, err
}

// RemoveBOM removes BOM from a file
func RemoveBOM(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Check if file has BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		// Remove BOM
		err = ioutil.WriteFile(filename, data[3:], 0644)
		if err != nil {
			return err
		}
		fmt.Printf("Removed BOM from file: %s\n", filename)
	}

	return nil
}
