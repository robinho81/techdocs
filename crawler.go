package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Markdownfile struct {
	Name string
	Path string
}

func findFilesInFolder(directory string, fileExtensions []string) (error, []Markdownfile) {

	markdownFilesFound := []Markdownfile{}

	f, err := os.Open(directory)
	if err != nil {
		fmt.Println("Error opening directory: " + err.Error())
		return err, markdownFilesFound
	}

	filesInDirectory, readDirErr := f.Readdir(-1) // This indicates that all FileInfos should be returned

	if readDirErr != nil {
		fmt.Println("Error reading from directory: " + readDirErr.Error())
		return readDirErr, markdownFilesFound
	}

	for _, file := range filesInDirectory {
		if hasSpecifiedFileExtensions(file.Name(), fileExtensions) {
			md := Markdownfile{Name: file.Name(), Path: file.Name()}
			markdownFilesFound = append(markdownFilesFound, md)
		}
	}
	return nil, markdownFilesFound
}

func findAllFilesRecursively(rootFolder string, fileExtensions []string) (error, []Markdownfile) {

	markdownFilesFound := []Markdownfile{}

	// calls the specified (anonymous) method to walk the folder structure
	// we do this so that we can use the parameters to this function within the method body
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		if hasSpecifiedFileExtensions(path, fileExtensions) {
			fileName := filepath.Base(path)
			markdownFile := Markdownfile{Name: fileName, Path: path}
			markdownFilesFound = append(markdownFilesFound, markdownFile)
		}
		return err
	})
	if err != nil {
		fmt.Printf("\n An error occurred")
	}
	return nil, markdownFilesFound
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
