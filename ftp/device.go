package ftp

import (
	"encoding/base64"
	"fmt"
	"tms/zipper"
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
	//return map[string]string{
	//	"hanan.txt":   "SGkKCk15IG5hbWUgaXMgSGFuYW4=",
	//	"meshack.txt": "TmlhamUsCk5haXR3YSBtZXNoYWNrCgpXZSByb2NrISE=",
	//}
	fileMap, err := zipper.ListTopLevelFilesWithBase64("Files.zip")
	if err != nil {
		fmt.Println("There was an error reading the zip file...")
		return map[string]string{}
	}
	return fileMap
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
