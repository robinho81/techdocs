package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func findFiles(rootFolder string, fileExtensions []string) []string {

	filesFound := []string{}

	// calls the specified (anonymous) method to walk the folder structure
	// we do this so that we can use the parameters to this function within the method body
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		if hasSpecifiedFileExtensions(path, fileExtensions) {
			filesFound = append(filesFound, path)
		}
		return err
	})
	if err != nil {
		fmt.Printf("\n An error occurred")
	}
	return filesFound
}

func hasSpecifiedFileExtensions(path string, fileExtensions []string) bool {
	for _, extensionToLookFor := range fileExtensions {
		currentFileExtension := filepath.Ext(path)
		if currentFileExtension == extensionToLookFor {
			return true
		}
	}
	return false
}
