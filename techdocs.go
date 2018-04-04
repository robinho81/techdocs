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
	outputFolderName := flag.String("o", "output", "Specify the output folder")
	flag.Parse()

	outputFolderError := checkOutputFolder(*outputFolderName)
	if outputFolderError != nil {
		fmt.Println("Error checking the output folder: " + outputFolderError.Error())
		return
	}

	markdownFileOrFolder := getFolderOrFileFromArgs()
	isDirectory, inputFolderErr := checkInputFolder(markdownFileOrFolder)
	if inputFolderErr != nil {
		fmt.Println("Error checking the specified folder: " + inputFolderErr.Error())
		return
	}

	if isDirectory {
		readDirectory(markdownFileOrFolder, *outputFolderName, *recursiveFlag)
	} else {
		readMarkdownFile(markdownFileOrFolder, *outputFolderName)
	}
}

func checkInputFolder(markdownFileOrFolder string) (bool, error) {
	stat, errStat := os.Stat(markdownFileOrFolder)

	if errStat != nil {
		fmt.Println("Error opening specified file or folder \"" + markdownFileOrFolder + "\": " + errStat.Error())
		return false, errStat
	}

	return stat.IsDir(), nil
}

func checkOutputFolder(outputFolder string) error {

	if outputFolder == "output" {
		fmt.Println("No output folder specified, creating folder \"output\"")
		createErr := os.Mkdir(outputFolder, os.ModePerm)
		if createErr != nil {
			return createErr
		}
	}

	_, errStat := os.Stat(outputFolder)
	if errStat != nil {
		return errStat
	}

	// TODO check if this is a folder

	return nil
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

	var markdownFiles []Markdownfile

	if isRecursive {
		markdownFiles = findAllFilesRecursively(folderPath, fileExtensions)
		fmt.Printf("Found markdown files: %d", len(markdownFiles))
	} else {
		_, markdownFiles = findFilesInFolder(folderPath, fileExtensions)
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
