package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func main() {
	recursiveFlag := flag.Bool("r", false, "Recursively search the specified directory for markdown files")
	outputFileName := flag.String("o", "output.html", "Specify the output file or folder")
	flag.Parse()

	markdownFileOrFolder := getFolderOrFileFromArgs()

	stat, errStat := os.Stat(markdownFileOrFolder)

	if errStat != nil {
		fmt.Println("Error opening specified file or folder \"" + markdownFileOrFolder + "\": " + errStat.Error())
		return
	}

	isDirectory := stat.IsDir()

	if isDirectory {
		readDirectory(markdownFileOrFolder, *outputFileName, *recursiveFlag)
	} else {
		readMarkdownFile(markdownFileOrFolder, *outputFileName)
	}
}

func getFolderOrFileFromArgs() string {
	numberOfArguments := len(os.Args)
	if numberOfArguments == 1 {
		fmt.Println("Please specify a markdown file or folder.")
		return ""
	}
	return os.Args[numberOfArguments-1]
}

func readDirectory(folderPath string, outputFileName string, isRecursive bool) {

	// TODO configure file extensions
	fileExtensions := []string{".md"}

	if isRecursive {

		markdownFiles := findAllFilesRecursively(folderPath, fileExtensions)
		fmt.Printf("Found markdown files: %d", len(markdownFiles))
	} else {
		_, markdownFiles := findFilesInFolder(folderPath, fileExtensions)
		fmt.Printf("Found markdown files: %d", len(markdownFiles))
	}
}

func readMarkdownFile(specifiedFilePath string, outputFileName string) {
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

	fmt.Println("Wrote file to " + outputFileName)
}
