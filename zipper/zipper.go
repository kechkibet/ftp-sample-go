package zipper

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"path/filepath"
	"strings"
)

func ListTopLevelFilesWithBase64(zipFilePath string) (map[string]string, error) {
	extensions := []string{".txt", ".PAR", ".par", ".JPG", ".jpg", ".PNG", ".png", ".PEM", ".pem", ".CRT", ".crt", ".KEY", ".key", ".BIN", ".bin"}
	// Open the ZIP file
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	fileMap := make(map[string]string)
	extMap := make(map[string]bool)

	// Create a map for quick extension checking
	for _, ext := range extensions {
		extMap[strings.ToLower(ext)] = true
	}

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Normalize file path to avoid issues with different path separators
		normalizedPath := filepath.ToSlash(f.Name)
		// Check if the file is a top-level file and not empty
		if !f.FileInfo().IsDir() && f.FileInfo().Size() > 0 && strings.Count(normalizedPath, "/") <= 1 {
			ext := strings.ToLower(filepath.Ext(normalizedPath))
			if _, ok := extMap[ext]; ok {
				// Open the file
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}

				// Read the file's contents
				contents, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					return nil, err
				}

				// Encode the contents in base64
				encoded := base64.StdEncoding.EncodeToString(contents)
				fileMap[filepath.Base(normalizedPath)] = encoded
			}
		}
	}

	return fileMap, nil
}
