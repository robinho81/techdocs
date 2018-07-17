package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	recursiveFlag := flag.Bool("r", false, "Recursively search the specified directory for markdown files")
	outputFolderName := flag.String("o", "", "Specify the output folder")
	persistData := flag.Bool("p", false, "Persist the data to mongo db")
	versionTag := flag.String("t", "1.0.0", "Specify the version tag")
	dbName := flag.String("d", "techdocs", "Specify the database / catalog name")
	// Try the default Mongo connection string if none is specified
	connectionString := flag.String("c", "mongodb://localhost:27017", "Specify connection string to Mongo")

	flag.Parse()

	if *persistData {
		fmt.Println("Version Tag: " + *versionTag)
		fmt.Println("Db name: " + *dbName)
		fmt.Println("ConnectionString: " + *connectionString)
	}

	if *outputFolderName != "" {
		outputFolderError := checkOutputFolder(*outputFolderName)
		if outputFolderError != nil {
			fmt.Println("Error checking the output folder: " + outputFolderError.Error())
			return
		}
	}

	inputFolder := getFolderOrFileFromArgs()
	isDirectory, inputFolderErr := checkInputFolder(inputFolder)
	if inputFolderErr != nil {
		fmt.Println("Error checking the specified folder: " + inputFolderErr.Error())
		return
	}

	// if an output folder is specified then save the files to disk, otherwise save to Db
	if *outputFolderName != "" {
		convertFilesAndWriteToDisk(isDirectory, *recursiveFlag, inputFolder, *outputFolderName)
	} else if *persistData {
		convertFilesAndSaveToDb(isDirectory, *recursiveFlag, inputFolder, *versionTag, *dbName, *connectionString)
	}
}

func convertFilesAndSaveToDb(isDirectory bool, recursiveFlag bool, markdownFileOrFolder string, versionTag string, dbName string, connectionString string) {
	start := time.Now()

	db, dbErr := connect(connectionString, dbName)
	if dbErr != nil {
		fmt.Println("Error connecting to db: " + dbErr.Error())
		return
	}

	removeAllFilesForVersion(db, versionTag)

	if isDirectory {
		markdownFiles := readDirectory(markdownFileOrFolder, recursiveFlag)
		for _, mdFile := range markdownFiles {
			html := convertMarkdownFileToHtml(mdFile.Path)
			saveHtmlFileToDb(db, mdFile.Name, html, versionTag)
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

func convertFilesAndWriteToDisk(isDirectory bool, recursiveFlag bool, markdownFileOrFolder string, outputFolderName string) {
	start := time.Now()
	if isDirectory {
		markdownFiles := readDirectory(markdownFileOrFolder, recursiveFlag)

		ch := make(chan string, len(markdownFiles))

		for _, mdFile := range markdownFiles {
			outputFilePath := generateOutputFilePath(mdFile, outputFolderName)
			go convertMarkdownFilesAndSaveInParallel(mdFile.Path, outputFilePath, ch)
		}

		for range markdownFiles {
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
