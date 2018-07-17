package main

import (
	"flag"
	"fmt"
	"os"

	conv "github.com/robinho81/techdocs/conversion"
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
		conv.ConvertMarkdownFilesAndWriteToDisk(isDirectory, *recursiveFlag, inputFolder, *outputFolderName)

	} else if *persistData {
		conv.ConvertMarkdownFilesAndSaveToDb(isDirectory, *recursiveFlag, inputFolder, *versionTag, *dbName, *connectionString)
		conv.FindAndSaveAllHintFilesToDb(*recursiveFlag, inputFolder, *connectionString, *dbName, *versionTag)
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
