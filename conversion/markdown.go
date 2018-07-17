package conversion

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/robinho81/techdocs/crawl"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func ConvertMarkdownFilesAndSaveToDb(isDirectory bool, recursiveFlag bool, markdownFileOrFolder string, versionTag string, dbName string, connectionString string) {
	start := time.Now()

	db, dbErr := connect(connectionString, dbName)
	if dbErr != nil {
		fmt.Println("Error connecting to db: " + dbErr.Error())
		return
	}

	removeAllItemsInCollection(db, versionTag, "pages")

	if isDirectory {
		files := readDirectory(markdownFileOrFolder, recursiveFlag)
		for _, file := range files {
			html := convertMarkdownFileToHtml(file.Path)
			saveHtmlFileToDb(db, file.Name, html, versionTag)
		}
	} else {
		// convert a single file
		html := convertMarkdownFileToHtml(markdownFileOrFolder)
		fileName := filepath.Base(markdownFileOrFolder)
		saveHtmlFileToDb(db, fileName, html, versionTag)
	}

	duration := time.Since(start).Seconds()
	fmt.Printf("Generated files in %f (s) \n", duration)
}

func ConvertMarkdownFilesAndWriteToDisk(isDirectory bool, recursiveFlag bool, markdownFileOrFolder string, outputFolderName string) {
	start := time.Now()
	if isDirectory {
		files := readDirectory(markdownFileOrFolder, recursiveFlag)

		ch := make(chan string, len(files))

		for _, file := range files {
			outputFilePath := generateOutputFilePath(file, outputFolderName)
			go convertMarkdownFilesAndSaveInParallel(file.Path, outputFilePath, ch)
		}

		for range files {
			filename := <-ch
			fmt.Println("Generated file: " + filename)
		}

	} else {
		outputFileName := convertMarkdownFileAndSave(markdownFileOrFolder, outputFolderName)
		fmt.Println("Generated file: " + outputFileName)
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("Generated files in %f (s)", duration)
}

func convertMarkdownFilesAndSaveInParallel(specifiedFilePath string, outputFileName string, ch chan string) {
	bytes, err := ioutil.ReadFile(specifiedFilePath)

	if err != nil {
		fmt.Println("Error reading markdown file " + err.Error())
		return
	}

	output := blackfriday.Run(bytes)
	writeFileErr := ioutil.WriteFile(outputFileName, output, 0644)

	if writeFileErr != nil {
		fmt.Println("Error writing the file " + writeFileErr.Error())
		return
	}
	ch <- outputFileName
}

func convertMarkdownFileToHtml(specifiedFilePath string) string {
	inputBytes, err := ioutil.ReadFile(specifiedFilePath)

	if err != nil {
		fmt.Println("Error reading markdown file " + err.Error())
		return ""
	}

	outputBytes := blackfriday.Run(inputBytes)
	html := string(outputBytes[:])
	return html
}

func convertMarkdownFileAndSave(specifiedFilePath string, outputFileName string) string {
	bytes, err := ioutil.ReadFile(specifiedFilePath)

	if err != nil {
		fmt.Println("Error reading markdown file " + err.Error())
		return ""
	}

	output := blackfriday.Run(bytes)
	writeFileErr := ioutil.WriteFile(outputFileName, output, 0644)

	if writeFileErr != nil {
		fmt.Println("Error writing the file " + writeFileErr.Error())
		return ""
	}

	fmt.Println("Wrote file to " + outputFileName)
	return outputFileName
}

func generateOutputFilePath(file crawl.Helpfile, outputFolder string) string {
	var extension = filepath.Ext(file.Name)
	var name = file.Name[0 : len(file.Name)-len(extension)]
	var newFilename = name + ".html"
	return filepath.Join(outputFolder, newFilename)
}

func readDirectory(folderPath string, isRecursive bool) []crawl.Helpfile {

	// TODO configure file extensions
	fileExtensions := []string{".md"}

	// TODO handle errors properly here
	if isRecursive {
		_, markdownFiles := crawl.FindAllFilesRecursively(folderPath, fileExtensions)
		return markdownFiles
	} else {
		_, markdownFiles := crawl.FindFilesInFolder(folderPath, fileExtensions)
		return markdownFiles
	}
}
