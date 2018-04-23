package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()

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
		markdownFiles := readDirectory(markdownFileOrFolder, *recursiveFlag)

		ch := make(chan string, len(markdownFiles))

		for _, mdFile := range markdownFiles {
			outputFilePath := generateOutputFilePath(mdFile, *outputFolderName)
			go convertMarkdownFileToHtmlInParallel(mdFile.Path, outputFilePath, ch)
		}

		for range markdownFiles {
			filename := <-ch
			fmt.Println("Generated file: " + filename)
		}

	} else {
		outputFileName := convertMarkdownFileToHtml(markdownFileOrFolder, *outputFolderName)
		fmt.Println("Generated file: " + outputFileName)
	}
	duration := time.Since(start).Seconds()
	fmt.Printf("Generated files in %f (s)", duration)
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
		if _, existsErr := os.Stat(outputFolder); os.IsNotExist(existsErr) {
			fmt.Println("No output folder specified, creating folder \"output\"")
			createErr := os.Mkdir(outputFolder, os.ModePerm)
			if createErr != nil {
				return createErr
			}
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

func readDirectory(folderPath string, isRecursive bool) []Markdownfile {

	// TODO configure file extensions
	fileExtensions := []string{".md"}

	// TODO handle errors properly here
	if isRecursive {
		_, markdownFiles := findAllFilesRecursively(folderPath, fileExtensions)
		return markdownFiles
	} else {
		_, markdownFiles := findFilesInFolder(folderPath, fileExtensions)
		return markdownFiles
	}
}
