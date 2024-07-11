package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `
<!DOCTYPE html>

<html lang="en">

<head>
  <meta charset="UTF-8" />
  <title>Markdown Preview</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <meta name="description" content="" />
  <link rel="icon" href="favicon.png">
</head>
<body>`

	footer = `
</body>
</html>`
)

func main() {
	fileFlag := flag.String("f", "", "Provide Markdown flag")
	flag.Parse()

	if *fileFlag == "" {
		flag.Usage()
	} else {
		err := Readfile(*fileFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func Readfile(file string) error {
	extension := path.Ext(file)
	if extension != ".md" {
		return fmt.Errorf("Provide a markdown file.")
	}

	_, errOpen := os.Open(file)
	if errOpen != nil {
		return errOpen
	}
	data, errRead := os.ReadFile(file)
	if errRead != nil {
		return errRead
	}

	parsedData := ParseFile(data)
	htmlFilename := fmt.Sprintf("%s.html", file)

	return SaveHTML(parsedData, htmlFilename)
}

func ParseFile(markdown []byte) []byte {
	output := blackfriday.Run(markdown)
	santizedHTML := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(santizedHTML)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func SaveHTML(data []byte, filename string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}
