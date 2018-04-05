package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func convertMarkdownFileToHtml(specifiedFilePath string, outputFileName string) {
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

func generateOutputFilePath(markdownFile Markdownfile, outputFolder string) string {
	var extension = filepath.Ext(markdownFile.Name)
	var name = markdownFile.Name[0 : len(markdownFile.Name)-len(extension)]
	var newFilename = name + ".html"
	return filepath.Join(outputFolder, newFilename)
}
