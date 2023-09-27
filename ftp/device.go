package ftp

import (
	"encoding/base64"
	"fmt"
)

func isValidSerialNumber(serial string) bool {
	if len(serial) != 10 {
		return false
	}

	for _, ch := range serial {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

func listFilesBySerialNumber(serialNumber string) []string {
	// Replace this block with your logic for listing files by serial number
	files := []string{"file1.txt", "file2.txt"} // Dummy file list
	return files
}

func getFileBytes(fileName string) ([]byte, error) {
	fileMap := map[string]string{
		"sample1.txt": "SGVsbG8sIFdvcmxkIQ==",         // "Hello, World!" in base64
		"sample2.txt": "VGhpcyBpcyBhIHNhbXBsZSB0ZXh0", // "This is a sample text" in base64
	}

	base64Content, exists := fileMap[fileName]
	if !exists {
		return nil, fmt.Errorf("file %s does not exist", fileName)
	}

	content, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 for file %s: %w", fileName, err)
	}

	return content, nil
}
