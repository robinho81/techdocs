package conversion

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/robinho81/techdocs/crawl"
)

type Hint struct {
	Key  string
	Text string
}

func FindAndSaveAllHintFilesToDb(isRecursive bool, folderPath string, connectionString string, dbName string, versionTag string) {
	files := getListOfHintFiles(isRecursive, folderPath)

	db, err := connect(connectionString, dbName)
	if err != nil {
		fmt.Println("Error connecting to db: " + err.Error())
		return
	}

	removeAllItemsInCollection(db, versionTag, "hints")

	for _, file := range files {
		hints, hintsErr := parseTextFileForHints(file.Path)
		if hintsErr != nil {
			fmt.Println("Error getting hints from file " + file.Path + ": " + hintsErr.Error())
		} else {
			saveHintsToDb(db, hints, versionTag)
		}
	}
}

func parseTextFileForHints(filePath string) ([]Hint, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var hints []Hint
	for scanner.Scan() {
		hint, hintErr := parseHint(scanner.Text())
		if hintErr != nil {
			fmt.Println("Error parsing hint in file " + filePath + ": " + hintErr.Error())
		} else {
			hints = append(hints, hint)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hints, nil
}

func parseHint(text string) (hint Hint, err error) {
	if text == "" {
		return hint, errors.New("Input line was empty")
	}

	if !strings.Contains(text, "=") {
		return hint, errors.New("Input line did not contain a valid seperator char")
	}
	components := strings.Split(text, "=")

	if len(components) < 2 {
		return hint, errors.New("Invalid hint format")
	}

	hint = Hint{Key: components[0], Text: components[1]}
	return hint, nil
}

func getListOfHintFiles(isRecursive bool, folderPath string) []crawl.Helpfile {

	fileExtensions := []string{".txt"}

	if isRecursive {
		recErr, files := crawl.FindAllFilesRecursively(folderPath, fileExtensions)

		if recErr != nil {
			fmt.Println("Error finding files recursively in folder " + folderPath + ": " + recErr.Error())
			return nil
		}

		return files
	} else {
		_, files := crawl.FindFilesInFolder(folderPath, fileExtensions)
		return files
	}
}
