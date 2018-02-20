package main

import (
	"fmt"
	"io/ioutil"
	"os"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify a markdown file.")
		return
	}

	markdownFile := os.Args[1]

	bytes, err := ioutil.ReadFile(markdownFile)

	if err != nil {
		fmt.Println("Error reading markdown file " + err.Error())
		return
	}

	fileName := "output.html"

	output := blackfriday.Run(bytes)

	writeFileErr := ioutil.WriteFile(fileName, output, 0644)

	if writeFileErr != nil {
		fmt.Println("Error writing the file " + writeFileErr.Error())
		return
	}

	fmt.Println("Wrote file to " + fileName)
}
