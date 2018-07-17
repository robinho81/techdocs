package crawl

import (
	"fmt"
	"os"
	"path/filepath"
)

type Helpfile struct {
	Name      string
	Path      string
	Extension string
}

func FindFilesInFolder(directory string, fileExtensions []string) (error, []Helpfile) {

	helpFilesFound := []Helpfile{}

	f, err := os.Open(directory)
	if err != nil {
		fmt.Println("Error opening directory: " + err.Error())
		return err, helpFilesFound
	}

	filesInDirectory, readDirErr := f.Readdir(-1) // This indicates that all FileInfos should be returned

	if readDirErr != nil {
		fmt.Println("Error reading from directory: " + readDirErr.Error())
		return readDirErr, helpFilesFound
	}

	for _, file := range filesInDirectory {
		if hasSpecifiedFileExtensions(file.Name(), fileExtensions) {
			md := Helpfile{Name: file.Name(), Path: file.Name(), Extension: getExtension(file.Name())}
			helpFilesFound = append(helpFilesFound, md)
		}
	}
	return nil, helpFilesFound
}

func FindAllFilesRecursively(rootFolder string, fileExtensions []string) (error, []Helpfile) {

	helpFilesFound := []Helpfile{}

	// calls the specified (anonymous) method to walk the folder structure
	// we do this so that we can use the parameters to this function within the method body
	err := filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}
		if hasSpecifiedFileExtensions(path, fileExtensions) {
			fileName := filepath.Base(path)
			helpFile := Helpfile{Name: fileName, Path: path, Extension: getExtension(path)}
			helpFilesFound = append(helpFilesFound, helpFile)
		}
		return err
	})
	if err != nil {
		fmt.Printf("\n An error occurred")
	}
	return nil, helpFilesFound
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

func getExtension(path string) string {
	return filepath.Ext(path)
}
