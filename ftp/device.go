package ftp

import (
	"encoding/base64"
	"fmt"
)

func isValidSerialNumber(serial string) bool {
	if len(serial) != 24 {
		return false
	}
	//203037333301059417812603
	for _, ch := range serial {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

func getFileMap() map[string]string {
	return map[string]string{
		"file1.txt": "SGVsbG8sIFdvcmxkIQ==",         // "Hello, World!" in base64
		"file2.txt": "VGhpcyBpcyBhIHNhbXBsZSB0ZXh0", // "This is a sample text" in base64
	}
}

func listFilesBySerialNumber(serialNumber string) []string {
	fileMap := getFileMap()

	// Extract file names from the map
	var files []string
	for fileName := range fileMap {
		files = append(files, fileName)
	}

	// You can add your logic here for filtering or processing based on serialNumber
	// For now, it returns all file names

	return files
}

func getFileBytes(fileName string) ([]byte, error) {
	fileMap := getFileMap()

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
